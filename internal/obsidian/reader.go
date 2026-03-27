// Package obsidian provides read-only access to an Obsidian vault.
// The canonical indexing pipeline lives in internal/markdownbrain; this package
// only parses vault markdown into stable note structs.
package obsidian

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// VaultNote represents a single Obsidian markdown file.
type VaultNote struct {
	// RelPath is the path relative to vault root, e.g. "projects/aurelia.md"
	RelPath     string
	Title       string
	Content     string // body text after frontmatter
	Tags        []string
	Frontmatter map[string]string // raw key=value pairs from YAML front matter
	ModifiedAt  time.Time
}

// ReadVault walks vaultPath and returns all .md files as VaultNotes.
// It skips hidden directories (starting with ".") except the vault root itself.
func ReadVault(vaultPath string) ([]VaultNote, error) {
	vaultPath = filepath.Clean(vaultPath)
	var notes []VaultNote

	err := filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip hidden sub-dirs (e.g. .obsidian, .trash) but not the root itself
			if path != vaultPath && strings.HasPrefix(d.Name(), ".") {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		relPath, err := filepath.Rel(vaultPath, path)
		if err != nil {
			return nil // skip unresolvable
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		note, err := parseNoteFile(path, relPath, info.ModTime())
		if err != nil {
			return nil // skip unparseable files, don't abort
		}
		notes = append(notes, note)
		return nil
	})
	return notes, err
}

func parseNoteFile(path, relPath string, modTime time.Time) (VaultNote, error) {
	f, err := os.Open(path)
	if err != nil {
		return VaultNote{}, err
	}
	defer f.Close()

	note := VaultNote{
		RelPath:     relPath,
		ModifiedAt:  modTime,
		Frontmatter: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	var lines []string
	inFrontmatter := false
	frontmatterDone := false
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if lineNum == 1 && line == "---" {
			inFrontmatter = true
			continue
		}
		if inFrontmatter && !frontmatterDone {
			if line == "---" || line == "..." {
				frontmatterDone = true
				continue
			}
			// Parse "key: value"
			if idx := strings.Index(line, ":"); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])
				note.Frontmatter[key] = val
				switch key {
				case "title":
					note.Title = val
				case "tags":
					note.Tags = parseTags(val)
				}
			}
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return note, err
	}

	note.Content = strings.TrimSpace(strings.Join(lines, "\n"))

	// Fall back to filename as title if frontmatter doesn't have one
	if note.Title == "" {
		base := filepath.Base(relPath)
		note.Title = strings.TrimSuffix(base, filepath.Ext(base))
	}

	return note, nil
}

// parseTags handles both YAML inline list "[tag1, tag2]" and plain "tag1, tag2"
func parseTags(raw string) []string {
	raw = strings.Trim(raw, "[]")
	parts := strings.Split(raw, ",")
	var tags []string
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}
