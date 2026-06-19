package data

import "testing"

func TestAIClientProfileHeadersDefaultIsEmpty(t *testing.T) {
	headers := AIClientProfileHeaders("", "")
	if len(headers) != 0 {
		t.Fatalf("expected no headers for empty profile, got %#v", headers)
	}
}

func TestAIClientProfileHeadersClaude(t *testing.T) {
	headers := AIClientProfileHeaders("claude", "2023-06-01")

	if got := headers["anthropic-version"]; got != "2023-06-01" {
		t.Fatalf("anthropic-version = %q, want %q", got, "2023-06-01")
	}
	if got := headers["user-agent"]; got != "claude-cli/2023-06-01" {
		t.Fatalf("user-agent = %q, want %q", got, "claude-cli/2023-06-01")
	}
	if got := headers["x-go-stock-client-profile"]; got != "claude" {
		t.Fatalf("x-go-stock-client-profile = %q, want %q", got, "claude")
	}
	if got := headers["x-go-stock-client-version"]; got != "2023-06-01" {
		t.Fatalf("x-go-stock-client-version = %q, want %q", got, "2023-06-01")
	}
}

func TestAIClientProfileHeadersCodex(t *testing.T) {
	headers := AIClientProfileHeaders("codex", "0.0.0")

	if got := headers["user-agent"]; got != "codex/0.0.0" {
		t.Fatalf("user-agent = %q, want %q", got, "codex/0.0.0")
	}
	if got := headers["x-stainless-package-version"]; got != "0.0.0" {
		t.Fatalf("x-stainless-package-version = %q, want %q", got, "0.0.0")
	}
	if got := headers["x-go-stock-client-profile"]; got != "codex" {
		t.Fatalf("x-go-stock-client-profile = %q, want %q", got, "codex")
	}
	if got := headers["x-go-stock-client-version"]; got != "0.0.0" {
		t.Fatalf("x-go-stock-client-version = %q, want %q", got, "0.0.0")
	}
}
