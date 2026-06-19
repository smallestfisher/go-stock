package data

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAIClientProfileHeadersDefaultIsEmpty(t *testing.T) {
	headers := AIClientProfileHeaders("", "")
	if len(headers) != 0 {
		t.Fatalf("expected no headers for empty profile, got %#v", headers)
	}
}

func TestAIClientProfileHeadersClaude(t *testing.T) {
	headers := AIClientProfileHeaders("claude", "2.1.178")

	if got := headers["accept"]; got != "application/json" {
		t.Fatalf("accept = %q, want %q", got, "application/json")
	}
	if got := headers["x-app"]; got != "cli" {
		t.Fatalf("x-app = %q, want %q", got, "cli")
	}
	if got := headers["anthropic-version"]; got != "2023-06-01" {
		t.Fatalf("anthropic-version = %q, want %q", got, "2023-06-01")
	}
	if got := headers["anthropic-beta"]; !strings.Contains(got, "claude-code-20250219") || !strings.Contains(got, "interleaved-thinking-2025-05-14") {
		t.Fatalf("anthropic-beta = %q", got)
	}
	if got := headers["user-agent"]; got != "claude-cli/2.1.178 (external, cli)" {
		t.Fatalf("user-agent = %q, want %q", got, "claude-cli/2.1.178 (external, cli)")
	}
	if got := headers["x-stainless-os"]; got != "Linux" {
		t.Fatalf("x-stainless-os = %q", got)
	}
	if got := headers["x-stainless-arch"]; got != "x64" {
		t.Fatalf("x-stainless-arch = %q", got)
	}
	if got := headers["x-stainless-lang"]; got != "js" {
		t.Fatalf("x-stainless-lang = %q", got)
	}
	if got := headers["x-stainless-runtime"]; got != "node" {
		t.Fatalf("x-stainless-runtime = %q", got)
	}
	if got := headers["x-stainless-timeout"]; got != "600" {
		t.Fatalf("x-stainless-timeout = %q", got)
	}
	if got := headers["x-stainless-retry-count"]; got != "0" {
		t.Fatalf("x-stainless-retry-count = %q", got)
	}
	if got := headers["x-stainless-package-version"]; got != "0.94.0" {
		t.Fatalf("x-stainless-package-version = %q", got)
	}
	if got := headers["x-stainless-runtime-version"]; got != "v24.3.0" {
		t.Fatalf("x-stainless-runtime-version = %q", got)
	}
	if got := headers["anthropic-dangerous-direct-browser-access"]; got != "true" {
		t.Fatalf("anthropic-dangerous-direct-browser-access = %q", got)
	}
	if got := headers["x-claude-code-session-id"]; got == "" {
		t.Fatalf("x-claude-code-session-id should be set")
	}
}

func TestAIClientProfileHeadersCodex(t *testing.T) {
	headers := AIClientProfileHeaders("codex", "0.139.0")

	if got := headers["accept"]; got != "application/json" {
		t.Fatalf("accept = %q, want %q", got, "application/json")
	}
	if got := headers["originator"]; got != "codex-tui" {
		t.Fatalf("originator = %q, want %q", got, "codex-tui")
	}
	if got := headers["user-agent"]; got != "codex-tui/0.139.0 (Ubuntu 24.4.0; x86_64) WindowsTerminal (codex-tui; 0.139.0)" {
		t.Fatalf("user-agent = %q", got)
	}
	if headers["session-id"] == "" || headers["thread-id"] == "" {
		t.Fatalf("session/thread headers should be set: %#v", headers)
	}
	if got := headers["session-id"]; got != headers["thread-id"] {
		t.Fatalf("session-id = %q, thread-id = %q, want same id", got, headers["thread-id"])
	}
	if got := headers["x-client-request-id"]; got != headers["session-id"] {
		t.Fatalf("x-client-request-id = %q, want session id %q", got, headers["session-id"])
	}
	if got := headers["x-codex-window-id"]; got != headers["session-id"]+":0" {
		t.Fatalf("x-codex-window-id = %q, want %q", got, headers["session-id"]+":0")
	}
	if got := headers["x-codex-beta-features"]; got != "terminal_resize_reflow" {
		t.Fatalf("x-codex-beta-features = %q", got)
	}
	if strings.Contains(headers["x-codex-turn-metadata"], headers["user-agent"]) {
		t.Fatalf("turn metadata should not include user agent")
	}
	var metadata map[string]any
	if err := json.Unmarshal([]byte(headers["x-codex-turn-metadata"]), &metadata); err != nil {
		t.Fatalf("metadata should be json: %v", err)
	}
	if got := metadata["session_id"]; got != headers["session-id"] {
		t.Fatalf("metadata session_id = %v, want %q", got, headers["session-id"])
	}
	if got := metadata["request_kind"]; got != "turn" {
		t.Fatalf("metadata request_kind = %v, want turn", got)
	}
}
