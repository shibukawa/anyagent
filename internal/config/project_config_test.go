package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegenerateAgentsFilePrefersProjectTemplates(t *testing.T) {
	dir := t.TempDir()
	// Prepare .anyagent template
	anyagentDir := filepath.Join(dir, ".anyagent")
	extraDir := filepath.Join(anyagentDir, "extra_rules")
	if err := os.MkdirAll(extraDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Simple template with placeholders
	tmpl := "# Project: {{PROJECT_NAME}}\n\n{{EXTRA_RULES}}\n"
	if err := os.WriteFile(filepath.Join(anyagentDir, "AGENTS.md.tmpl"), []byte(tmpl), 0644); err != nil {
		t.Fatalf("write tmpl: %v", err)
	}
	// Local extra rule for go
	goRule := "# Go Local Rule\ncontent"
	if err := os.WriteFile(filepath.Join(extraDir, "go.md"), []byte(goRule), 0644); err != nil {
		t.Fatalf("write rule: %v", err)
	}

	cfg := &ProjectConfig{
		ProjectName:    "X",
		InstalledRules: []string{"go"},
		Parameters:     map[string]string{"PROJECT_NAME": "X"},
	}
	if err := cfg.RegenerateAgentsFileAt(dir); err != nil {
		t.Fatalf("regenerate: %v", err)
	}
	// Verify AGENTS.md contains local rule header
	b, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	s := string(b)
	if want := "# Go Local Rule"; !strings.Contains(s, want) {
		t.Fatalf("AGENTS.md missing local rule: %s", want)
	}
}
