package skill

import (
	"strings"
	"testing"
)

func TestChunkSkill_SplitsByHeadingAndParagraphBudget(t *testing.T) {
	t.Parallel()

	body := strings.Repeat("linha longa para chunking seguro. ", 80)
	skill := Skill{
		Metadata: Metadata{
			Description: "skill de teste",
		},
		Content: "---\nname: test-skill\ndescription: skill de teste\n---\n\n# Intro\n\nResumo curto.\n\n# Deep Dive\n\n" + body,
	}

	chunks := ChunkSkill("test-skill", skill)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for _, chunk := range chunks {
		if chunk.Section == "" || chunk.ID == "" || chunk.Checksum == "" {
			t.Fatalf("invalid chunk metadata: %#v", chunk)
		}
		if !strings.Contains(chunk.Text, "Skill: test-skill") {
			t.Fatalf("chunk missing skill prefix: %q", chunk.Text)
		}
	}
}

func TestChunkSkill_FallsBackToDescription(t *testing.T) {
	t.Parallel()

	chunks := ChunkSkill("fallback-skill", Skill{
		Metadata: Metadata{Description: "descricao fallback"},
	})
	if len(chunks) != 1 {
		t.Fatalf("expected one fallback chunk, got %d", len(chunks))
	}
	if !strings.Contains(chunks[0].Text, "descricao fallback") {
		t.Fatalf("chunk missing fallback description: %q", chunks[0].Text)
	}
}
