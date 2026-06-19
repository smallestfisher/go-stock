package data

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	AIClientProfileDefault = ""
	AIClientProfileClaude  = "claude"
	AIClientProfileCodex   = "codex"
)

func AIClientProfileHeaders(profile, version string) map[string]string {
	profile = strings.ToLower(strings.TrimSpace(profile))
	version = strings.TrimSpace(version)
	if profile == "" || profile == "default" || profile == "none" {
		return map[string]string{}
	}

	headers := map[string]string{
		"x-go-stock-client-profile": profile,
	}
	if version != "" {
		headers["x-go-stock-client-version"] = version
	}

	switch profile {
	case AIClientProfileClaude:
		if version == "" {
			version = "2.1.178"
		}
		headers = map[string]string{
			"accept":                                    "application/json",
			"x-app":                                     "cli",
			"user-agent":                                fmt.Sprintf("claude-cli/%s (external, cli)", version),
			"anthropic-beta":                            "claude-code-20250219,interleaved-thinking-2025-05-14,redact-thinking-2026-02-12,context-management-2025-06-27,prompt-caching-scope-2026-01-05,mid-conversation-system-2026-04-07,effort-2025-11-24",
			"anthropic-version":                         "2023-06-01",
			"x-stainless-os":                            "Linux",
			"x-stainless-arch":                          "x64",
			"x-stainless-lang":                          "js",
			"x-stainless-runtime":                       "node",
			"x-stainless-timeout":                       "600",
			"x-stainless-retry-count":                   "0",
			"x-stainless-package-version":               "0.94.0",
			"x-stainless-runtime-version":               "v24.3.0",
			"x-claude-code-session-id":                  newClientProfileID(),
			"anthropic-dangerous-direct-browser-access": "true",
		}
	case AIClientProfileCodex:
		if version == "" {
			version = "0.139.0"
		}
		sessionID := newClientProfileID()
		turnID := newClientProfileID()
		windowID := sessionID + ":0"
		headers = map[string]string{
			"accept":                "application/json",
			"originator":            "codex-tui",
			"session-id":            sessionID,
			"thread-id":             sessionID,
			"user-agent":            fmt.Sprintf("codex-tui/%s (Ubuntu 24.4.0; x86_64) WindowsTerminal (codex-tui; %s)", version, version),
			"x-client-request-id":   sessionID,
			"x-codex-beta-features": "terminal_resize_reflow",
			"x-codex-window-id":     windowID,
		}
		metadata := map[string]any{
			"session_id":              sessionID,
			"thread_id":               sessionID,
			"thread_source":           "user",
			"turn_id":                 turnID,
			"sandbox":                 "seccomp",
			"turn_started_at_unix_ms": time.Now().UnixMilli(),
			"request_kind":            "turn",
			"window_id":               windowID,
		}
		if raw, err := json.Marshal(metadata); err == nil {
			headers["x-codex-turn-metadata"] = string(raw)
		}
	}

	return headers
}

func newClientProfileID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		now := time.Now().UnixNano()
		return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", uint32(now>>32), uint16(now>>16), uint16(now), uint16(now>>48), now&0xffffffffffff)
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	)
}

func ApplyAIClientProfileHeaders(req *resty.Request, profile, version string) *resty.Request {
	for k, v := range AIClientProfileHeaders(profile, version) {
		req = req.SetHeader(k, v)
	}
	return req
}
