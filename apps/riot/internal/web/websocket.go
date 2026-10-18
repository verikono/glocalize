package web

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	// wsPath is the websocket endpoint the browser client connects to.
	wsPath = "/ws"

	// wsWriteTimeout bounds a single write to a client, so one stalled connection cannot hold the broadcaster.
	wsWriteTimeout = 5 * time.Second

	// wsSendBuffer is how many messages may queue for one client before it is considered behind and skipped.
	wsSendBuffer = 8
)

// wsClient is one connected browser. send is buffered so the broadcaster never blocks on a slow client.
type wsClient struct {
	conn *websocket.Conn
	send chan string
}

// hub fans a message out to every connected websocket client.
type hub struct {
	mu      sync.RWMutex
	clients map[*wsClient]struct{}
}

func newHub() *hub {
	return &hub{clients: make(map[*wsClient]struct{})}
}

func (h *hub) add(client *wsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = struct{}{}
}

func (h *hub) remove(client *wsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client)
}

// broadcast delivers message to every client. A client whose buffer is full is skipped: it is already behind and the next message will catch it up.
func (h *hub) broadcast(message string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
		}
	}
}

// websocketHandler upgrades the request to a websocket and serves the client until it disconnects.
func (w *RiotWebUi) websocketHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(rw, r, nil)
		if err != nil {
			return
		}

		client := &wsClient{conn: conn, send: make(chan string, wsSendBuffer)}
		w.hub.add(client)
		defer func() {
			w.hub.remove(client)
			_ = conn.CloseNow()
		}()

		w.serveClient(r.Context(), client)
	})
}

// serveClient writes queued messages to the client until it disconnects. A reader runs alongside because the write side alone cannot tell a quiet client from a closed one.
func (w *RiotWebUi) serveClient(ctx context.Context, client *wsClient) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer cancel()
		for {
			if _, _, err := client.conn.Read(ctx); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case message := <-client.send:
			writeCtx, writeCancel := context.WithTimeout(ctx, wsWriteTimeout)
			err := client.conn.Write(writeCtx, websocket.MessageText, []byte(message))
			writeCancel()
			if err != nil {
				return
			}
		}
	}
}
