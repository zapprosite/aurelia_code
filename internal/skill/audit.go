package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var skillMarkdownLinkPattern = regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)

type AuditIssue struct {
	Severity  string
	Code      string
	SkillName string
	Path      string
	Message   string
}

type AuditReport struct {
	ScannedSkills int
	Issues        []AuditIssue
}

func (r AuditReport) ErrorCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			count++
		}
	}
	return count
}

func (r AuditReport) WarningCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Severity == "warning" {
			count++
		}
	}
	return count
}

// AuditCatalog inspects raw skill directories before merge/override so silent structural drift is visible.
func AuditCatalog(baseDirs ...string) (AuditReport, error) {
	report := AuditReport{}
	seenNames := make(map[string]string)

	for _, baseDir := range baseDirs {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return report, fmt.Errorf("audit skills directory %q: %w", baseDir, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			skillDir := filepath.Join(baseDir, entry.Name())
			skillPath := filepath.Join(skillDir, "SKILL.md")
			data, err := os.ReadFile(skillPath)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				report.Issues = append(report.Issues, AuditIssue{
					Severity: "error",
					Code:     "read_failed",
					Path:     skillPath,
					Message:  err.Error(),
				})
				continue
			}

			meta, content, err := parseFrontmatter(data)
			if err != nil {
				report.Issues = append(report.Issues, AuditIssue{
					Severity: "error",
					Code:     "invalid_frontmatter",
					Path:     skillPath,
					Message:  err.Error(),
				})
				continue
			}

			report.ScannedSkills++
			if strings.TrimSpace(meta.Description) == "" {
				report.Issues = append(report.Issues, AuditIssue{
					Severity:  "warning",
					Code:      "missing_description",
					SkillName: meta.Name,
					Path:      skillPath,
					Message:   "description is empty",
				})
			}

			if prev, ok := seenNames[meta.Name]; ok {
				report.Issues = append(report.Issues, AuditIssue{
					Severity:  "warning",
					Code:      "duplicate_name",
					SkillName: meta.Name,
					Path:      skillPath,
					Message:   fmt.Sprintf("duplicates %s", prev),
				})
			} else {
				seenNames[meta.Name] = skillPath
			}

			for _, issue := range auditSkillMarkdownLinks(skillPath, content) {
				issue.SkillName = meta.Name
				report.Issues = append(report.Issues, issue)
			}
		}
	}

	return report, nil
}

func auditSkillMarkdownLinks(sourcePath, content string) []AuditIssue {
	var issues []AuditIssue
	matches := skillMarkdownLinkPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		target := normalizeSkillMarkdownTarget(match[1])
		if target == "" {
			continue
		}

		resolved := target
		if !filepath.IsAbs(target) {
			resolved = filepath.Join(filepath.Dir(sourcePath), filepath.FromSlash(target))
		}

		if _, err := os.Stat(resolved); err != nil {
			issues = append(issues, AuditIssue{
				Severity: "error",
				Code:     "broken_link",
				Path:     sourcePath,
				Message:  fmt.Sprintf("%s -> %s", match[1], err),
			})
		}
	}
	return issues
}

func normalizeSkillMarkdownTarget(target string) string {
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
