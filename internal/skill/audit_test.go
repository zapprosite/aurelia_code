package skill

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAuditCatalog_DetectsDuplicateNamesAndBrokenLinks(t *testing.T) {
	t.Parallel()

	dir1 := t.TempDir()
	dir2 := t.TempDir()

	writeSkillMD(t, dir1, "alpha", "Shared Skill", "alpha description")
	writeSkillMD(t, dir2, "beta", "Shared Skill", "beta description")

	brokenDir := filepath.Join(dir2, "broken")
	writeSkillMD(t, dir2, "broken", "Broken Skill", "broken description")
	skillPath := filepath.Join(brokenDir, "SKILL.md")
	content := "---\nname: Broken Skill\ndescription: broken description\n---\n\nVeja [algo](./missing.md).\n"
	if err := os.WriteFile(skillPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write broken skill: %v", err)
	}

	report, err := AuditCatalog(dir1, dir2)
	if err != nil {
		t.Fatalf("AuditCatalog() error = %v", err)
	}
	if report.WarningCount() == 0 {
		t.Fatalf("expected duplicate warning, got %#v", report.Issues)
	}
	if report.ErrorCount() == 0 {
		t.Fatalf("expected broken-link error, got %#v", report.Issues)
	}
}

func TestAuditCatalog_RepositoryCanonicalCatalogIsClean(t *testing.T) {
	t.Parallel()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))

	report, err := AuditCatalog(filepath.Join(repoRoot, ".agent", "skills"))
	if err != nil {
		t.Fatalf("AuditCatalog() error = %v", err)
	}
	if report.ErrorCount() != 0 {
		t.Fatalf("expected clean canonical catalog, got %#v", report.Issues)
	}
}
