package modules

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-stock/server"
)

func callRPC(t *testing.T, method string, args []any) map[string]any {
	t.Helper()
	body, err := json.Marshal(map[string]any{"args": args})
	if err != nil {
		t.Fatalf("marshal rpc body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/rpc/"+method, bytes.NewReader(body))
	rec := httptest.NewRecorder()
	server.RPCDispatcher(&server.Core{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s status = %d, body = %s", method, rec.Code, rec.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode %s response: %v", method, err)
	}
	return payload
}

func TestGetSponsorInfoReportsFreeEntitlement(t *testing.T) {
	payload := callRPC(t, "GetSponsorInfo", nil)

	if got := int(payload["vipLevel"].(float64)); got != 2 {
		t.Fatalf("vipLevel = %d, want 2", got)
	}
	if active, _ := payload["active"].(bool); !active {
		t.Fatal("active = false, want true")
	}
	if got, _ := payload["sponsorCode"].(string); got != "" {
		t.Fatalf("sponsorCode = %q, want empty", got)
	}
}

func TestCheckSponsorCodeNoLongerValidatesCodes(t *testing.T) {
	for _, args := range [][]any{nil, {""}, {"not-hex"}} {
		payload := callRPC(t, "CheckSponsorCode", args)
		if got := int(payload["code"].(float64)); got != 1 {
			t.Fatalf("code = %d, want 1 for args %#v", got, args)
		}
	}
}

func TestCheckDeviceBindingDoesNotBlockLocalFreeUse(t *testing.T) {
	payload := callRPC(t, "CheckDeviceBinding", []any{"token", "https://example.invalid"})

	if bound, _ := payload["bound"].(bool); !bound {
		t.Fatal("bound = false, want true")
	}
	if deviceCount := int(payload["deviceCount"].(float64)); deviceCount != 0 {
		t.Fatalf("deviceCount = %d, want 0", deviceCount)
	}
}
