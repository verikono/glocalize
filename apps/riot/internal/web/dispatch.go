package web

import (
	"fmt"
	"strings"
	"time"

	T "github.com/verikono/glocalize/types/go"
)

// dispatchBuffer bounds the queued dispatch messages. A burst of file changes
// must never block the watcher goroutine, so the queue is bounded and a full
// queue drops the message: a pending refresh already covers the change.
const dispatchBuffer = 64

// Dispatch queues message for the web UI. It never blocks, so the caller (the
// fsnotify watcher) is never held up by a slow or absent browser.
func (rw *RiotWebUi) Dispatch(message string) {
	select {
	case rw.dispatch <- message:
	default:
	}
}

// OnDispatch handles a message delivered to the web UI. A refresh request is
// broadcast to every connected websocket client with a timestamp; any other
// message is broadcast verbatim, so the same path carries non-HMR traffic.
func (rw *RiotWebUi) OnDispatch(message string) {
	if strings.HasPrefix(message, T.HmrRefreshMessage) {
		rw.hub.broadcast(fmt.Sprintf("%s::%d", T.HmrRefreshMessage, time.Now().UnixMilli()))
		return
	}
	rw.hub.broadcast(message)
}

// run drains the dispatch channel into OnDispatch. It is the single owner of
// the channel and exits when the channel is closed.
func (rw *RiotWebUi) run() {
	for message := range rw.dispatch {
		rw.OnDispatch(message)
	}
}
