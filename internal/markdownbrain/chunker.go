package markdownbrain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

const maxChunkChars = 1500

type Chunk struct {
	ID       string
	Index    int
	Count    int
	Section  string
	Text     string
	Checksum string
}

type section struct {
	title string
	body  string
}

func ChunkDocument(doc Document) []Chunk {
	content := strings.TrimSpace(doc.Content)
	if content == "" {
		return nil
	}

	sections := splitSections(content)
	rawChunks := make([]Chunk, 0, len(sections))
	for _, sec := range sections {
		body := strings.TrimSpace(sec.body)
		if body == "" {
			continue
		}

		prefix := buildChunkPrefix(doc, sec.title)
		limit := maxChunkChars - len(prefix)
		if limit < 320 {
			limit = 320
		}

		for _, part := range splitText(body, limit) {
			text := prefix + part
			checksum := checksumText(text)
			rawChunks = append(rawChunks, Chunk{
				Section:  sec.title,
				Text:     text,
				Checksum: checksum,
			})
		}
	}

	for idx := range rawChunks {
		rawChunks[idx].Index = idx
		rawChunks[idx].Count = len(rawChunks)
		rawChunks[idx].ID = checksumText(fmt.Sprintf("%s|%s|%d|%s", buildSourceID(doc), rawChunks[idx].Section, idx, rawChunks[idx].Checksum))
	}

	return rawChunks
}

func buildChunkPrefix(doc Document, sectionTitle string) string {
	var b strings.Builder
	b.WriteString("Title: ")
	b.WriteString(doc.Title)
	b.WriteString("\nSource: ")
	b.WriteString(doc.SourceKind)
	b.WriteString("\nPath: ")
	b.WriteString(chunkPath(doc))
	b.WriteString("\nSection: ")
	b.WriteString(sectionTitle)
	if len(doc.Tags) > 0 {
		b.WriteString("\nTags: ")
		b.WriteString(strings.Join(doc.Tags, ", "))
	}
	b.WriteString("\n\n")
	return b.String()
}

func chunkPath(doc Document) string {
	if doc.RelPath == "" {
		return doc.Title
	}
	return filepath.ToSlash(doc.RelPath)
}

func splitSections(content string) []section {
	lines := strings.Split(content, "\n")
	sections := make([]section, 0, 4)
	current := section{title: "Overview"}

	flush := func() {
		current.body = strings.TrimSpace(current.body)
		if current.body == "" {
			return
		}
		sections = append(sections, current)
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if title, ok := parseHeading(trimmed); ok {
			flush()
			current = section{title: title}
			continue
		}
		if current.body == "" {
			current.body = line
		} else {
			current.body += "\n" + line
		}
	}
	flush()
	return sections
}

func parseHeading(line string) (string, bool) {
	if !strings.HasPrefix(line, "#") {
		return "", false
	}
	title := strings.TrimSpace(strings.TrimLeft(line, "#"))
	if title == "" {
		return "", false
	}
	return title, true
}

func splitText(text string, limit int) []string {
	if len(text) <= limit {
		return []string{text}
	}

	var chunks []string
	var current strings.Builder
	paragraphs := strings.Split(text, "\n\n")

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		if current.Len() > 0 && current.Len()+2+len(paragraph) > limit {
			chunks = append(chunks, current.String())
			current.Reset()
		}
		if len(paragraph) > limit {
			if current.Len() > 0 {
				chunks = append(chunks, current.String())
				current.Reset()
			}
			chunks = append(chunks, splitLongParagraph(paragraph, limit)...)
			continue
		}
		if current.Len() > 0 {
			current.WriteString("\n\n")
		}
		current.WriteString(paragraph)
	}
	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}
	return chunks
}

func splitLongParagraph(paragraph string, limit int) []string {
	var chunks []string
	remaining := strings.TrimSpace(paragraph)
	for len(remaining) > limit {
		cut := strings.LastIndex(remaining[:limit], " ")
		if cut < 1 {
			cut = limit
		}
		chunks = append(chunks, strings.TrimSpace(remaining[:cut]))
		remaining = strings.TrimSpace(remaining[cut:])
	}
	if remaining != "" {
		chunks = append(chunks, remaining)
	}
	return chunks
}

func checksumText(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}
