package data

import (
	"fmt"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"go-stock/backend/logger"
)

// 雪球行情页：匿名访问即可，雪球会通过 Set-Cookie 直接下发反爬 token。
const xueqiuQuotePage = "https://xueqiu.com/hq"

// 雪球反爬关键 cookie：命中其一即认为拿到有效匿名 token，
// 足以调用 hot_stock / hot_event 等公开行情接口。
var xueqiuKeyCookieNames = map[string]bool{
	"xq_a_token":  true,
	"xqat":        true,
	"xq_r_token":  true,
	"xq_id_token": true,
	"cookiesu":    true,
	"u":           true,
}

// fetchXueqiuCookiesViaHTTP 纯 HTTP 获取雪球反爬 cookie，无需本地浏览器。
//
// 桌面 Wails 环境能起本地 Chrome（chromedp），但 web server 部署的机器上往往没有浏览器，
// 导致 CheckBrowser 失败、热榜空数据。实测访问行情页 https://xueqiu.com/hq 时，
// 雪球会通过 Set-Cookie 直接下发 xq_a_token / xqat / xq_r_token / xq_id_token / cookiesu / u
// 等匿名 token，带上即可正常请求公开行情接口，完全不依赖浏览器。
func fetchXueqiuCookiesViaHTTP(pageURL string) (string, error) {
	target := strings.TrimSpace(pageURL)
	if target == "" {
		target = xueqiuQuotePage
	}

	// 复用 shared transport（IPv4 强制 + 代理设置），但挂独立 cookie jar，
	// 避免污染 SharedHTTPClient 的全局状态。
	client := CreateHTTPClientWithTimeout(30 * time.Second)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", fmt.Errorf("创建 cookie jar 失败: %w", err)
	}
	client.GetClient().Jar = jar

	_, err = client.R().
		SetHeader("User-Agent", getRandomUA()).
		SetHeader("Referer", "https://xueqiu.com/").
		Get(target)
	if err != nil {
		return "", fmt.Errorf("访问雪球行情页失败: %w", err)
	}

	u, err := url.Parse("https://xueqiu.com")
	if err != nil {
		return "", err
	}
	cookies := jar.Cookies(u)
	if len(cookies) == 0 {
		return "", fmt.Errorf("雪球未下发任何 cookie")
	}

	var b strings.Builder
	hasKey := false
	for i, c := range cookies {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(c.Name)
		b.WriteByte('=')
		b.WriteString(c.Value)
		if xueqiuKeyCookieNames[c.Name] {
			hasKey = true
		}
	}
	if !hasKey {
		return "", fmt.Errorf("雪球 cookie 缺少关键反爬 token（拿到 %d 个 cookie）", len(cookies))
	}

	logger.SugaredLogger.Infof("雪球 HTTP 获取 cookie 成功（%d 个）", len(cookies))
	return b.String(), nil
}

// GetXueqiuCookieHeader 雪球反爬 cookie 的统一入口，供 hot_stock / hot_event 等接口复用。
//
// 策略：先查缓存 → 纯 HTTP 获取（无需浏览器，适配 web server 部署）→ 失败再降级 chromedp
// （桌面 Wails 环境兜底）。任一路径成功都写入共享缓存（TTL = XueqiuCookieCacheTTL）。
// 全部失败时返回空串 + error，调用方可选择不带 cookie 继续请求。
func GetXueqiuCookieHeader(pageURL string) (string, error) {
	target := strings.TrimSpace(pageURL)
	if target == "" {
		target = xueqiuQuotePage
	}

	now := time.Now()
	cacheKey := "xueqiu-cookie||" + getURLCacheKey(target)

	xueqiuCookieCache.mu.Lock()
	if item, ok := xueqiuCookieCache.items[cacheKey]; ok && now.Before(item.expiry) {
		header := item.header
		xueqiuCookieCache.mu.Unlock()
		return header, nil
	}
	xueqiuCookieCache.mu.Unlock()

	// 1) 首选纯 HTTP，不依赖本地浏览器。
	header, httpErr := fetchXueqiuCookiesViaHTTP(target)
	if httpErr == nil && header != "" {
		storeXueqiuCookie(cacheKey, header, now)
		return header, nil
	}
	logger.SugaredLogger.Warnf("雪球 HTTP 获取 cookie 失败，尝试降级 chromedp: %v", httpErr)

	// 2) 降级 chromedp（桌面环境兜底）。注意 chromedp 路径自带缓存，但其 cacheKey
	//    含 browserPath，与此处不同；这里再写一份到统一 key，便于后续命中。
	h, cdpErr := fetchXueqiuCookiesViaChromedp("", 30*time.Second, target)
	if cdpErr == nil && h != "" {
		storeXueqiuCookie(cacheKey, h, time.Now())
		return h, nil
	}

	return "", fmt.Errorf("雪球 cookie 获取失败（HTTP: %v；chromedp: %v）", httpErr, cdpErr)
}

func storeXueqiuCookie(cacheKey, header string, now time.Time) {
	xueqiuCookieCache.mu.Lock()
	xueqiuCookieCache.items[cacheKey] = &cookieCacheItem{
		header: header,
		expiry: now.Add(XueqiuCookieCacheTTL),
	}
	xueqiuCookieCache.mu.Unlock()
}
