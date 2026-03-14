package util

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sync"
)

var (
	// Multiple listeners for web mode
	listeners   = make(map[chan WebEvent]bool)
	listenersMu sync.Mutex
	IsWebMode   = false
	IsReady     = false // Track if initialization is done
)

type WebEvent struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// RegisterListener adds a new channel to receive events
func RegisterListener() chan WebEvent {
	ch := make(chan WebEvent, 100)
	listenersMu.Lock()
	listeners[ch] = true
	listenersMu.Unlock()
	return ch
}

// UnregisterListener removes a channel
func UnregisterListener(ch chan WebEvent) {
	listenersMu.Lock()
	delete(listeners, ch)
	listenersMu.Unlock()
	close(ch)
}

// Emit redirects events to either Wails or all registered web event channels
func Emit(ctx context.Context, name string, data interface{}) {
	if name == "loadingMsg" && data == "done" {
		IsReady = true
	}

	if IsWebMode {
		event := WebEvent{Name: name, Data: data}
		listenersMu.Lock()
		defer listenersMu.Unlock()
		for ch := range listeners {
			select {
			case ch <- event:
			default:
				// Channel full, skip to avoid blocking
			}
		}
	} else {
		runtime.EventsEmit(ctx, name, data)
	}
}
