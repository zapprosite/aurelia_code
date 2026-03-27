package tools

import (
	"context"
	"os/exec"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
)

// GrepSearchHandler searches for text patterns in files using grep.
func GrepSearchHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	query, err := requireStringArg(args, "query")
	if err != nil {
		return "", err
	}
	path, err := requireStringArg(args, "path")
	if err != nil {
		path = "."
	}

	workdir, _ := agent.WorkdirFromContext(ctx)
	
	// Use ripgrep (rg) if available, fallback to grep -r
	cmdName := "rg"
	if _, err := exec.LookPath("rg"); err != nil {
		cmdName = "grep"
	}

	var cmdArgs []string
	if cmdName == "rg" {
		cmdArgs = []string{"--column", "--line-number", "--no-heading", "--color", "never", "-e", query, path}
	} else {
		cmdArgs = []string{"-rnE", query, path}
	}

	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	if workdir != "" {
		cmd.Dir = workdir
	}

	out, _ := cmd.CombinedOutput()
	result := string(out)
	
	if result == "" {
		return "Nenhum resultado encontrado.", nil
	}

	// Truncate if too long (standard agent limit)
	if len(result) > 30000 {
		result = result[:30000] + "\n... [RESULTADO TRUNCADO]"
	}

	return result, nil
}

// FindFilesHandler searches for files by name/pattern using find.
func FindFilesHandler(ctx context.Context, args map[string]interface{}) (string, error) {
	pattern, err := requireStringArg(args, "pattern")
	if err != nil {
		return "", err
	}
	path := optionalStringArg(args, "path")
	if path == "" {
		path = "."
	}

	workdir, _ := agent.WorkdirFromContext(ctx)

	// find <path> -name <pattern>
	cmd := exec.CommandContext(ctx, "find", path, "-name", pattern)
	if workdir != "" {
		cmd.Dir = workdir
	}

	out, _ := cmd.CombinedOutput()
	result := strings.TrimSpace(string(out))

	if result == "" {
		return "Nenhum arquivo encontrado.", nil
	}

	return result, nil
}
