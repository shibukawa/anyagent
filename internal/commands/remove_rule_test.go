package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/anyagent/internal/config"
)

func TestRunRemoveRule(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create AGENTS.md to simulate initialized project
	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("# Test Project"), 0644); err != nil {
		t.Fatalf("Failed to create AGENTS.md: %v", err)
	}

	// Create instructions directory and rule file
	instructionsDir := filepath.Join(tempDir, ".github", "instructions")
	if err := os.MkdirAll(instructionsDir, 0755); err != nil {
		t.Fatalf("Failed to create instructions directory: %v", err)
	}

	ruleFilePath := filepath.Join(instructionsDir, "go.instructions.md")
	if err := os.WriteFile(ruleFilePath, []byte("# Go Rules"), 0644); err != nil {
		t.Fatalf("Failed to create rule file: %v", err)
	}

	tests := []struct {
		name        string
		language    string
		expectError bool
	}{
		{
			name:        "Valid rule removal",
			language:    "go",
			expectError: false,
		},
		{
			name:        "Non-existent rule",
			language:    "python",
			expectError: true,
		},
		{
			name:        "Invalid language",
			language:    "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunRemoveRule(tt.language, tempDir, true) // Use dry run

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRunRemoveRuleQDev(t *testing.T) {
	tempDir := t.TempDir()
	// AGENTS.md present
	if err := os.WriteFile(filepath.Join(tempDir, "AGENTS.md"), []byte("# AGENTS"), 0644); err != nil {
		t.Fatalf("failed to write AGENTS.md: %v", err)
	}
	// Enable qdev
	if err := config.SaveProjectConfig(tempDir, &config.ProjectConfig{EnabledAgents: []string{"qdev"}}); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}
	// Create Q Dev rule file
	rulesDir := filepath.Join(tempDir, ".amazonq", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("failed to create rules dir: %v", err)
	}
	rulePath := filepath.Join(rulesDir, "go.md")
	if err := os.WriteFile(rulePath, []byte("# Go Rules"), 0644); err != nil {
		t.Fatalf("failed to create rule file: %v", err)
	}
	// Dry run should succeed
	if err := RunRemoveRule("go", tempDir, true); err != nil {
		t.Fatalf("dry run failed: %v", err)
	}
}

func TestRunListRules(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create AGENTS.md to simulate initialized project
	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("# Test Project"), 0644); err != nil {
		t.Fatalf("Failed to create AGENTS.md: %v", err)
	}

	// Create instructions directory and some rule files
	instructionsDir := filepath.Join(tempDir, ".github", "instructions")
	if err := os.MkdirAll(instructionsDir, 0755); err != nil {
		t.Fatalf("Failed to create instructions directory: %v", err)
	}

	// Create go.instructions.md
	goRuleFilePath := filepath.Join(instructionsDir, "go.instructions.md")
	if err := os.WriteFile(goRuleFilePath, []byte("# Go Rules"), 0644); err != nil {
		t.Fatalf("Failed to create go rule file: %v", err)
	}

	err := RunListRules(tempDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRunListRulesNotInitialized(t *testing.T) {
	// Create temporary directory without AGENTS.md
	tempDir := t.TempDir()

	err := RunListRules(tempDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
