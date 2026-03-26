package dashboard

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

//go:embed dist
var content embed.FS

// Event representa uma mensagem enviada pro dashboard em realtime
type Event struct {
	Type      string      `json:"type"`
	Agent     string      `json:"agent"`
	Action    string      `json:"action"`
	Payload   interface{} `json:"payload,omitempty"`
	Timestamp string      `json:"timestamp"`
	BotID     string      `json:"bot_id,omitempty"` // S-32: multi-bot source identifier
}

const (
	ringBufferSize = 500
	replayCount    = 50
)

var (
	subscribers = make(map[chan Event]bool)
	subMu       sync.Mutex

	// ring buffer: stores the last ringBufferSize events for replay on reconnect
	ringBuf  [ringBufferSize]Event
	ringHead int // index of oldest event (write position)
	ringLen  int // number of valid entries (0..ringBufferSize)

	customRoutes   = make(map[string]http.HandlerFunc)
	customRoutesMu sync.Mutex
)

// RegisterRoute allows registering external handlers (like /api/squad) safely.
func RegisterRoute(path string, handler http.HandlerFunc) {
	customRoutesMu.Lock()
	defer customRoutesMu.Unlock()
	customRoutes[path] = handler
}

// Publish envia um evento para todos os clientes conectados ao dashboard
// e armazena no ring buffer para replay em novos connects.
func Publish(e Event) {
	subMu.Lock()
	defer subMu.Unlock()

	// Store in ring buffer
	ringBuf[ringHead] = e
	ringHead = (ringHead + 1) % ringBufferSize
	if ringLen < ringBufferSize {
		ringLen++
	}

	for ch := range subscribers {
		select {
		case ch <- e:
		default:
			// subscriber channel full — drop rather than block
		}
	}
}

// recentEvents returns up to n events from the ring buffer (oldest first).
// Must be called with subMu held.
func recentEvents(n int) []Event {
	if ringLen == 0 || n <= 0 {
		return nil
	}
	if n > ringLen {
		n = ringLen
	}
	out := make([]Event, n)
	// start position: walk back n steps from ringHead
	start := (ringHead - ringLen + ringBufferSize) % ringBufferSize
	// advance start to skip the oldest we don't want
	skip := ringLen - n
	start = (start + skip) % ringBufferSize
	for i := 0; i < n; i++ {
		out[i] = ringBuf[(start+i)%ringBufferSize]
	}
	return out
}

// StartServer inicia o servidor Web do Dashboard ULTRATRINK na porta configurada.
// Se port <= 0, usa o padrão 3334.
func StartServer(logger *slog.Logger, port int) error {
	if port <= 0 {
		port = 3334
	}
	subFS, err := fs.Sub(content, "dist")
	if err != nil {
		logger.Error("erro ao carregar arquivos estáticos do dashboard", slog.Any("err", err))
		return err
	}

	mux := http.NewServeMux()

	// Endpoint de Real-time (SSE) com replay dos últimos replayCount eventos
	mux.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ch := make(chan Event, 10)

		subMu.Lock()
		replay := recentEvents(replayCount)
		subscribers[ch] = true
		subMu.Unlock()

		defer func() {
			subMu.Lock()
			delete(subscribers, ch)
			subMu.Unlock()
			close(ch)
		}()

		logger.Debug("novo cliente SSE conectado ao dashboard", slog.Int("replay", len(replay)))

		flusher, hasFlusher := w.(http.Flusher)

		// Send replay events first
		for _, event := range replay {
			data, _ := json.Marshal(event)
			fmt.Fprintf(w, "data: %s\n\n", data)
		}
		if hasFlusher && len(replay) > 0 {
			flusher.Flush()
		}

		for {
			select {
			case event := <-ch:
				data, _ := json.Marshal(event)
				fmt.Fprintf(w, "data: %s\n\n", data)
				if hasFlusher {
					flusher.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	})

	mux.Handle("/", http.FileServer(http.FS(subFS)))

	customRoutesMu.Lock()
	for path, handler := range customRoutes {
		mux.HandleFunc(path, handler)
	}
	customRoutesMu.Unlock()

	addr := ":" + strconv.Itoa(port)
	go func() {
		logger.Info("ULTRATRINK Dashboard Online", slog.String("url", "http://localhost:"+strconv.Itoa(port)))
		if err := http.ListenAndServe(addr, mux); err != nil {
			logger.Error("servidor do dashboard parou", slog.Any("err", err))
		}
	}()

	return nil
}
