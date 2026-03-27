package skill

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const maxSkillChunkChars = 1400

type SkillChunk struct {
	ID       string
	Index    int
	Count    int
	Section  string
	Text     string
	Checksum string
}

type skillSection struct {
	title string
	body  string
}

// ChunkSkill splits a skill markdown document into semantic sections for vector indexing.
func ChunkSkill(name string, skill Skill) []SkillChunk {
	body := strings.TrimSpace(stripSkillFrontmatter(skill.Content))
	if body == "" {
		body = strings.TrimSpace(skill.Metadata.Description)
	}
	if body == "" {
		body = name
	}

	sections := splitSkillSections(body)
	var rawChunks []SkillChunk
	for _, section := range sections {
		sectionBody := strings.TrimSpace(section.body)
		if sectionBody == "" {
			continue
		}

		prefix := fmt.Sprintf("Skill: %s\nDescription: %s\nSection: %s\n\n", name, skill.Metadata.Description, section.title)
		limit := maxSkillChunkChars - len(prefix)
		if limit < 240 {
			limit = 240
		}

		for _, piece := range splitSkillText(sectionBody, limit) {
			text := prefix + piece
			checksum := checksumText(text)
			rawChunks = append(rawChunks, SkillChunk{
				Section:  section.title,
				Text:     text,
				Checksum: checksum,
			})
		}
	}

	if len(rawChunks) == 0 {
		text := fmt.Sprintf("Skill: %s\nDescription: %s\nSection: Overview\n\n%s", name, skill.Metadata.Description, body)
		checksum := checksumText(text)
		rawChunks = append(rawChunks, SkillChunk{
			Section:  "Overview",
			Text:     text,
			Checksum: checksum,
		})
	}

	for idx := range rawChunks {
		rawChunks[idx].Index = idx
		rawChunks[idx].Count = len(rawChunks)
		rawChunks[idx].ID = checksumText(fmt.Sprintf("%s|%s|%d|%s", name, rawChunks[idx].Section, idx, rawChunks[idx].Checksum))
	}

	return rawChunks
}

func stripSkillFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---") {
		return content
	}

	lines := strings.Split(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return content
	}

	for idx := 1; idx < len(lines); idx++ {
		if strings.TrimSpace(lines[idx]) == "---" {
			return strings.Join(lines[idx+1:], "\n")
		}
	}
	return content
}

func splitSkillSections(body string) []skillSection {
	lines := strings.Split(body, "\n")
	var sections []skillSection
	current := skillSection{title: "Overview"}

	flush := func() {
		current.body = strings.TrimSpace(current.body)
		if current.body == "" {
			return
		}
		sections = append(sections, current)
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if title, ok := parseSkillHeading(trimmed); ok {
			flush()
			current = skillSection{title: title}
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

func parseSkillHeading(line string) (string, bool) {
	if !strings.HasPrefix(line, "#") {
		return "", false
	}
	title := strings.TrimLeft(line, "#")
	title = strings.TrimSpace(title)
	if title == "" {
		return "", false
	}
	return title, true
}

func splitSkillText(text string, limit int) []string {
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
			chunks = append(chunks, splitLongSkillParagraph(paragraph, limit)...)
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
	if len(chunks) == 0 {
		return []string{text}
	}
	return chunks
}

func splitLongSkillParagraph(paragraph string, limit int) []string {
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
