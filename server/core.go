package server

import (
	"context"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/robfig/cron/v3"
	"go-stock/backend/data"
	"go-stock/server/events"
)

// Core 持有进程级的运行时状态：缓存、定时任务、AI 工具集、AI 流式取消句柄、告警状态等。
// HTTP handler 通过它访问共享状态。
type Core struct {
	Cache   *freecache.Cache
	Cron    *cron.Cron
	AiTools []data.Tool
	Events  *events.Hub

	cronEntrys   map[string]cron.EntryID
	cronEntrysMu sync.Mutex

	// AI 流式任务的取消句柄
	summaryMu      sync.Mutex
	summaryCancels map[string]context.CancelFunc
	agentMu        sync.Mutex
	agentCancels   map[string]context.CancelFunc

	// 股价告警去重状态
	stockAlertMu       sync.Mutex
	stockAlertLastSent map[string]time.Time
	priceAtAlertReset  map[string]float64
}

// NewCore 构造并启动 Core（cron 在此启动）。
func NewCore() *Core {
	cache := freecache.NewCache(512 * 1024)
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger)))
	c.Start()
	tools := data.Tools(nil)
	return &Core{
		Cache:              cache,
		Cron:               c,
		AiTools:            tools,
		Events:             events.NewHub(),
		cronEntrys:         make(map[string]cron.EntryID),
		summaryCancels:     make(map[string]context.CancelFunc),
		agentCancels:       make(map[string]context.CancelFunc),
		stockAlertLastSent: make(map[string]time.Time),
		priceAtAlertReset:  make(map[string]float64),
	}
}

// --- 定时任务条目管理（对应 App 的 setCronEntry/getCronEntry/removeCronEntry）---

func (c *Core) SetCronEntry(key string, id cron.EntryID) {
	c.cronEntrysMu.Lock()
	c.cronEntrys[key] = id
	c.cronEntrysMu.Unlock()
}

func (c *Core) GetCronEntry(key string) (cron.EntryID, bool) {
	c.cronEntrysMu.Lock()
	defer c.cronEntrysMu.Unlock()
	id, ok := c.cronEntrys[key]
	return id, ok
}

func (c *Core) RemoveCronEntry(key string) {
	c.cronEntrysMu.Lock()
	delete(c.cronEntrys, key)
	c.cronEntrysMu.Unlock()
}

// --- AI 流式任务取消句柄（Phase 2 的 SSE 端点使用）---

func streamScope(scope string) string {
	if scope == "" {
		return "__default__"
	}
	return scope
}

func (c *Core) SetSummaryCancel(scope string, cf context.CancelFunc) {
	scope = streamScope(scope)
	c.summaryMu.Lock()
	if c.summaryCancels == nil {
		c.summaryCancels = make(map[string]context.CancelFunc)
	}
	old := c.summaryCancels[scope]
	if cf == nil {
		delete(c.summaryCancels, scope)
	} else {
		c.summaryCancels[scope] = cf
	}
	c.summaryMu.Unlock()
	if old != nil && cf != nil {
		old()
	}
}

func (c *Core) CancelSummary(scope string) {
	scope = streamScope(scope)
	c.summaryMu.Lock()
	cf := c.summaryCancels[scope]
	delete(c.summaryCancels, scope)
	c.summaryMu.Unlock()
	if cf != nil {
		cf()
	}
}

func (c *Core) SetAgentCancel(scope string, cf context.CancelFunc) {
	scope = streamScope(scope)
	c.agentMu.Lock()
	if c.agentCancels == nil {
		c.agentCancels = make(map[string]context.CancelFunc)
	}
	old := c.agentCancels[scope]
	if cf == nil {
		delete(c.agentCancels, scope)
	} else {
		c.agentCancels[scope] = cf
	}
	c.agentMu.Unlock()
	if old != nil && cf != nil {
		old()
	}
}

func (c *Core) CancelAgent(scope string) {
	scope = streamScope(scope)
	c.agentMu.Lock()
	cf := c.agentCancels[scope]
	delete(c.agentCancels, scope)
	c.agentMu.Unlock()
	if cf != nil {
		cf()
	}
}
