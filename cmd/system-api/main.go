package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kocar/aurelia/internal/purity/alog"
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Uptime    string            `json:"uptime"`
}

var startTime time.Time

func init() {
	startTime = time.Now()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "Soberano",
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		Services: map[string]string{
			"ollama":   checkService(getEnv("OLLAMA_URL", "http://localhost:11434")),
			"qdrant":   checkService(getEnv("QDRANT_URL", "http://localhost:6333") + "/readyz"),
			"postgres": checkTCP(getEnv("POSTGRES_HOST", "localhost"), getEnv("POSTGRES_PORT", "5432")),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func checkService(url string) string {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "🔴 DOWN"
	}
	return "🟢 UP"
}

func envAuditHandler(w http.ResponseWriter, r *http.Request) {
	// Simulação de auditoria de paridade .env vs .env.example
	// No futuro: Ler os arquivos fisicamente e comparar chaves
	audit := map[string]interface{}{
		"status":      "Conforme",
		"drift_count": 0,
		"missing":     []string{},
		"timestamp":   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audit)
}

func memorySyncHandler(w http.ResponseWriter, r *http.Request) {
	// Integrando com as variáveis de ambiente extraídas do .env
	qdrantURL := getEnv("QDRANT_URL", "http://localhost:6333")
	
	// Simulação de sincronização L2 <-> L3
	syncInfo := map[string]interface{}{
		"status":           "Sincronização Ativa",
		"qdrant_target":    "aurelia_markdown_brain",
		"postgres_source":  "app_journal.entries",
		"last_sync":        time.Now(),
		"items_processed":  42, // Mock por enquanto
		"drift_detected":   false,
	}

	alog.Info("omni-memory sync triggered", alog.With("url", qdrantURL))
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(syncInfo)
}

func checkTCP(host, port string) string {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return "🔴 DOWN"
	}
	conn.Close()
	return "🟢 UP"
}

func main() {
	alog.Configure(alog.Options{})
	port := "8081"
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/env-audit", envAuditHandler)
	http.HandleFunc("/memory-sync", memorySyncHandler)

	alog.Info("Aurelia System API starting", alog.With("port", port), alog.With("mode", "Go-Native SOTA 2026"))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		alog.Error("System API failed to start", alog.With("err", err))
		os.Exit(1)
	}
}
