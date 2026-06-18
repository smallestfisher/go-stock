package events

import (
	"testing"
	"time"
)

func TestEmitToDeliversOnlyToMatchingClient(t *testing.T) {
	h := NewHub()
	clientA := h.subscribe("client-a")
	clientB := h.subscribe("client-b")
	unscoped := h.subscribe("")
	defer h.unsubscribe(clientA)
	defer h.unsubscribe(clientB)
	defer h.unsubscribe(unscoped)

	h.EmitTo("client-a", "newChatStream", "chunk-a")

	assertEvent(t, clientA.ch, "newChatStream", "chunk-a")
	assertNoEvent(t, clientB.ch)
	assertNoEvent(t, unscoped.ch)
}

func TestEmitBroadcastsToAllClients(t *testing.T) {
	h := NewHub()
	clientA := h.subscribe("client-a")
	clientB := h.subscribe("client-b")
	defer h.unsubscribe(clientA)
	defer h.unsubscribe(clientB)

	h.Emit("warnMsg", "market-open")

	assertEvent(t, clientA.ch, "warnMsg", "market-open")
	assertEvent(t, clientB.ch, "warnMsg", "market-open")
}

func assertEvent(t *testing.T, ch <-chan eventPayload, wantName string, wantData any) {
	t.Helper()
	select {
	case got := <-ch:
		if got.name != wantName || got.data != wantData {
			t.Fatalf("event = (%q, %#v), want (%q, %#v)", got.name, got.data, wantName, wantData)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timed out waiting for event %q", wantName)
	}
}

func assertNoEvent(t *testing.T, ch <-chan eventPayload) {
	t.Helper()
	select {
	case got := <-ch:
		t.Fatalf("unexpected event: (%q, %#v)", got.name, got.data)
	case <-time.After(50 * time.Millisecond):
	}
}
