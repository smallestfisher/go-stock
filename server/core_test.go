package server

import "testing"

func TestSummaryCancelIsScopedByClient(t *testing.T) {
	core := &Core{}
	cancelledA := false
	cancelledB := false

	core.SetSummaryCancel("client-a", func() { cancelledA = true })
	core.SetSummaryCancel("client-b", func() { cancelledB = true })

	core.CancelSummary("client-b")

	if cancelledA {
		t.Fatal("client-a summary stream was cancelled by client-b")
	}
	if !cancelledB {
		t.Fatal("client-b summary stream was not cancelled")
	}
}

func TestAgentCancelIsScopedByClient(t *testing.T) {
	core := &Core{}
	cancelledA := false
	cancelledB := false

	core.SetAgentCancel("client-a", func() { cancelledA = true })
	core.SetAgentCancel("client-b", func() { cancelledB = true })

	core.CancelAgent("client-a")

	if !cancelledA {
		t.Fatal("client-a agent stream was not cancelled")
	}
	if cancelledB {
		t.Fatal("client-b agent stream was cancelled by client-a")
	}
}
