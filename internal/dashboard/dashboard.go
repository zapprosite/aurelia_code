package dashboard

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
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
}

var (
	subscribers = make(map[chan Event]bool)
	subMu       sync.Mutex

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
func Publish(e Event) {
	subMu.Lock()
	defer subMu.Unlock()
	for ch := range subscribers {
		ch <- e
	}
}

// StartServer inicia o servidor Web do Dashboard ULTRATRINK na porta 3333.
func StartServer(logger *slog.Logger) error {
	subFS, err := fs.Sub(content, "dist")
	if err != nil {
		logger.Error("erro ao carregar arquivos estáticos do dashboard", slog.Any("err", err))
		return err
	}

	mux := http.NewServeMux()
	
	// Endpoint de Real-time (SSE)
	mux.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		ch := make(chan Event, 10)
		subMu.Lock()
		subscribers[ch] = true
		subMu.Unlock()

		defer func() {
			subMu.Lock()
			delete(subscribers, ch)
			subMu.Unlock()
			close(ch)
		}()

		logger.Debug("novo cliente SSE conectado ao dashboard")

		for {
			select {
			case event := <-ch:
				data, _ := json.Marshal(event)
				fmt.Fprintf(w, "data: %s\n\n", data)
				if flusher, ok := w.(http.Flusher); ok {
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

	go func() {
		logger.Info("ULTRATRINK Dashboard Online", slog.String("url", "http://localhost:3334"))
		if err := http.ListenAndServe(":3334", mux); err != nil {
			logger.Error("servidor do dashboard parou", slog.Any("err", err))
		}
	}()
	
	return nil
}
