package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"

	"github.com/verikono/glocalize/apps/riot/internal/web"
)

// TestDispatchReachesWebsocketClient proves the dispatch path end to end: a message handed to the web UI is written to a connected websocket client.
func TestDispatchReachesWebsocketClient(t *testing.T) {
	ui, err := web.NewWebService()
	if err != nil {
		t.Fatalf("NewWebService: %v", err)
	}

	mux := http.NewServeMux()
	if err := ui.BindMux(mux); err != nil {
		t.Fatalf("BindMux: %v", err)
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatalf("dial %s: %v", url, err)
	}
	defer conn.CloseNow()

	// Dispatch is queued and the client is registered asynchronously, so retry until the message lands rather than assuming ordering.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		ui.Dispatch("riot::refresh")

		readCtx, readCancel := context.WithTimeout(ctx, 500*time.Millisecond)
		_, data, err := conn.Read(readCtx)
		readCancel()
		if err != nil {
			continue
		}
		if !strings.HasPrefix(string(data), "riot::refresh") {
			t.Fatalf("message = %q, want prefix %q", data, "riot::refresh")
		}
		return
	}
	t.Fatal("no message received on the websocket")
}

// TestDispatchIsNonBlocking proves a dispatch never blocks the caller, even
// with no clients connected and the queue full.
func TestDispatchIsNonBlocking(t *testing.T) {
	ui, err := web.NewWebService()
	if err != nil {
		t.Fatalf("NewWebService: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 1000; i++ {
			ui.Dispatch("riot::refresh")
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Dispatch blocked")
	}
}
