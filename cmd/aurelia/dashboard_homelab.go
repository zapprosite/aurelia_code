package main

import (
	"encoding/json"
	"net/http"

	"github.com/kocar/aurelia/internal/homelab"
)

func buildHomelabSnapshotHandler() http.HandlerFunc {
	collector := homelab.NewCollector()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_ = json.NewEncoder(w).Encode(collector.Collect(r.Context()))
	}
}
