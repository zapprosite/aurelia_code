package markdownbrain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadRepository_IncludesStrategicMarkdownAndSkipsVendoredNoise(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	mustWriteFile(t, filepath.Join(repoRoot, "AGENTS.md"), "# AGENTS\n\nCore governance")
	mustWriteFile(t, filepath.Join(repoRoot, ".agent", "skills", "alpha", "SKILL.md"), "---\ntitle: Skill Alpha\ntags: [skill, core]\n---\n\n# Uso\n\nConteudo")
	mustWriteFile(t, filepath.Join(repoRoot, "node_modules", "pkg", "README.md"), "# ignore")
	mustWriteFile(t, filepath.Join(repoRoot, "homelab-bibliotheca", "skills", "open-claw", "skills", "x", "README.md"), "# vendored")

	docs, err := ReadRepository(repoRoot)
	if err != nil {
		t.Fatalf("ReadRepository() error = %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("expected 2 markdown docs, got %d", len(docs))
	}
}

func TestReadObsidianVault_MapsVaultNotesToDocuments(t *testing.T) {
	t.Parallel()

	vault := t.TempDir()
	mustWriteFile(t, filepath.Join(vault, "projects", "aurelia.md"), "---\ntitle: Aurelia\ntags: [brain, code]\n---\n\n# Context\n\nRuntime")

	docs, err := ReadObsidianVault(vault)
	if err != nil {
		t.Fatalf("ReadObsidianVault() error = %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}
	if docs[0].SourceSystem != vaultSourceSystem || docs[0].SourceKind != "vault_note" {
		t.Fatalf("unexpected doc metadata: %#v", docs[0])
	}
}

func TestChunkDocument_CreatesStableChunks(t *testing.T) {
	t.Parallel()

	doc := Document{
		SourceSystem: repoSourceSystem,
		SourceKind:   "repo_document",
		RelPath:      "docs/brain.md",
		Title:        "Brain",
		Content:      "# One\n\n" + strings.Repeat("context ", 120) + "\n\n# Two\n\nTail",
	}
	chunks := ChunkDocument(doc)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for _, chunk := range chunks {
		if chunk.ID == "" || chunk.Section == "" || chunk.Checksum == "" {
			t.Fatalf("invalid chunk: %#v", chunk)
		}
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
