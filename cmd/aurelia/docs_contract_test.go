package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var markdownLinkPattern = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)

func TestGovernanceEntrypointLinksResolve(t *testing.T) {
	t.Parallel()

	repoRoot := repoRootFromCaller(t)
	files := []string{
		filepath.Join(repoRoot, "AGENTS.md"),
		filepath.Join(repoRoot, ".agent", "rules", "README.md"),
		filepath.Join(repoRoot, "docs", "governance", "REPOSITORY_CONTRACT.md"),
		filepath.Join(repoRoot, "docs", "adr", "README.md"),
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		matches := markdownLinkPattern.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			target := normalizeMarkdownTarget(match[1])
			if target == "" {
				continue
			}

			resolved := target
			if !filepath.IsAbs(target) {
				resolved = filepath.Join(filepath.Dir(file), filepath.FromSlash(target))
			}

			if _, err := os.Stat(resolved); err != nil {
				t.Fatalf("broken markdown link in %s -> %s (%v)", file, match[1], err)
			}
		}
	}
}

func repoRootFromCaller(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func normalizeMarkdownTarget(target string) string {
	target = strings.TrimSpace(target)
	if target == "" || strings.HasPrefix(target, "#") {
		return ""
	}
	for _, prefix := range []string{"http://", "https://", "mailto:", "tel:"} {
		if strings.HasPrefix(target, prefix) {
			return ""
		}
	}
	if idx := strings.Index(target, "#"); idx >= 0 {
		target = target[:idx]
	}
	return filepath.Clean(filepath.FromSlash(target))
}
