package data

import (
	"strings"

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
			version = "2023-06-01"
		}
		headers["anthropic-version"] = version
		headers["user-agent"] = "claude-cli/" + version
		headers["x-go-stock-client-version"] = version
	case AIClientProfileCodex:
		if version == "" {
			version = "0.0.0"
		}
		headers["user-agent"] = "codex/" + version
		headers["x-stainless-package-version"] = version
		headers["x-go-stock-client-version"] = version
	}

	return headers
}

func ApplyAIClientProfileHeaders(req *resty.Request, profile, version string) *resty.Request {
	for k, v := range AIClientProfileHeaders(profile, version) {
		req = req.SetHeader(k, v)
	}
	return req
}
