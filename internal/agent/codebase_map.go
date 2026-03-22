package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/observability"
)

type CodebaseMapPayload struct {
	Architecture struct {
		Layers []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"layers"`
	} `json:"architecture"`
	KeyFiles []struct {
		Path string `json:"path"`
		Role string `json:"role"`
	} `json:"keyFiles"`
	Structure struct {
		TopDirectories []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"topDirectories"`
	} `json:"structure"`
}

var (
	cachedMapSummary string
	lastLoadTime     time.Time
	mapMu            sync.Mutex
)

// GetCodebaseMapSummary reads the .context/docs/codebase-map.json file and formats it into a dense markdown summary.
// It caches the output for 5 minutes to avoid excessive I/O during tight LLM loops.
func GetCodebaseMapSummary() string {
	mapMu.Lock()
	defer mapMu.Unlock()

	// Cache TTL of 5 minutes
	if time.Since(lastLoadTime) < 5*time.Minute && cachedMapSummary != "" {
		return cachedMapSummary
	}

	logger := observability.Logger("agent.codebasemap")

	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	mapPath := filepath.Join(cwd, ".context", "docs", "codebase-map.json")
	data, err := os.ReadFile(mapPath)
	if err != nil {
		logger.Debug("codebase-map.json not found or unreadable", "err", err)
		return ""
	}

	var payload CodebaseMapPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		logger.Warn("failed to parse codebase-map.json", "err", err)
		return ""
	}

	var sb strings.Builder

	sb.WriteString("Arquitetura (Layers):\n")
	for _, layer := range payload.Architecture.Layers {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", layer.Name, layer.Description))
	}

	sb.WriteString("\nDiretorios Top-Level:\n")
	for _, dir := range payload.Structure.TopDirectories {
		sb.WriteString(fmt.Sprintf("- %s/: %s\n", dir.Name, dir.Description))
	}

	sb.WriteString("\nArquivos Chaves (Key Files):\n")
	for _, file := range payload.KeyFiles {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", file.Path, file.Role))
	}

	cachedMapSummary = sb.String()
	lastLoadTime = time.Now()

	return cachedMapSummary
}
