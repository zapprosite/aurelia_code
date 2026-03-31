package markdownbrain

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	DefaultCollection = "aurelia_markdown_brain"
	repoSourceSystem  = "repo_markdown"
	vaultSourceSystem = "obsidian"
)

type Metadata struct {
	Title string
	Tags  []string
}

type Document struct {
	SourceSystem string
	SourceKind   string
	RootPath     string
	RelPath      string
	Title        string
	Content      string
	Tags         []string
	Frontmatter  map[string]string
	ModifiedAt   time.Time
}

func ReadRepository(repoRoot string) ([]Document, error) {
	repoRoot = filepath.Clean(strings.TrimSpace(repoRoot))
	if repoRoot == "" {
		return nil, nil
	}

	var docs []Document
	err := filepath.WalkDir(repoRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if shouldSkipRepoDir(repoRoot, path, d) {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}
		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		doc, err := parseMarkdownDocument(path, relPath, repoRoot, repoSourceSystem, "repo_document", info.ModTime())
		if err != nil {
			return nil
		}
		docs = append(docs, doc)
		return nil
	})
	return docs, err
}

func ReadObsidianVault(vaultPath string) ([]Document, error) {
	vaultPath = filepath.Clean(strings.TrimSpace(vaultPath))
	if vaultPath == "" {
		return nil, nil
	}

	// Como o internal/obsidian foi removido, tratamos o vault como um repositório markdown padrão.
	// No futuro, podemos adicionar lógica específica para links [[Obsidian]] aqui.
	var docs []Document
	err := filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if shouldSkipRepoDir(vaultPath, path, d) {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}
		relPath, err := filepath.Rel(vaultPath, path)
		if err != nil {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		doc, err := parseMarkdownDocument(path, relPath, vaultPath, vaultSourceSystem, "vault_note", info.ModTime())
		if err != nil {
			return nil
		}
		docs = append(docs, doc)
		return nil
	})
	return docs, err
}

func parseMarkdownDocument(path, relPath, rootPath, sourceSystem, sourceKind string, modTime time.Time) (Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return Document{}, err
	}
	defer func() { _ = f.Close() }()

	doc := Document{
		SourceSystem: sourceSystem,
		SourceKind:   sourceKind,
		RootPath:     rootPath,
		RelPath:      filepath.ToSlash(relPath),
		Frontmatter:  make(map[string]string),
		ModifiedAt:   modTime.UTC(),
	}

	scanner := bufio.NewScanner(f)
	var lines []string
	inFrontmatter := false
	frontmatterDone := false
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if lineNum == 1 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			continue
		}
		if inFrontmatter && !frontmatterDone {
			if trimmed := strings.TrimSpace(line); trimmed == "---" || trimmed == "..." {
				frontmatterDone = true
				continue
			}
			if idx := strings.Index(line, ":"); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				value := normalizeFrontmatterValue(line[idx+1:])
				doc.Frontmatter[key] = value
				switch key {
				case "title":
					doc.Title = value
				case "tags":
					doc.Tags = parseTags(value)
				}
			}
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return Document{}, err
	}

	doc.Content = strings.TrimSpace(strings.Join(lines, "\n"))
	if doc.Title == "" {
		doc.Title = firstMarkdownHeading(doc.Content)
	}
	if doc.Title == "" {
		base := filepath.Base(relPath)
		doc.Title = strings.TrimSuffix(base, filepath.Ext(base))
	}
	return doc, nil
}

func shouldSkipRepoDir(repoRoot, path string, entry fs.DirEntry) bool {
	name := entry.Name()
	switch name {
	case ".git", "node_modules", ".aurelia":
		return true
	}

	relPath, err := filepath.Rel(repoRoot, path)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	switch relPath {
	case "homelab-bibliotheca/skills/open-claw/skills":
		return true
	default:
		return false
	}
}

func normalizeFrontmatterValue(raw string) string {
	value := strings.TrimSpace(raw)
	value = strings.Trim(value, `"'`)
	return value
}

func firstMarkdownHeading(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimLeft(line, "#"))
		if line != "" {
			return line
		}
	}
	return ""
}

func parseTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	raw = strings.Trim(raw, "[]")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		tag := strings.TrimSpace(strings.Trim(part, `"'`))
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		tags = append(tags, tag)
	}
	return tags
}
