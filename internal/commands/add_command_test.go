package commands

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/shibukawa/anyagent/internal/config"
)

func TestRunAddCommand(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create AGENTS.md to simulate initialized project
	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("# Test Project"), 0644); err != nil {
		t.Fatalf("Failed to create AGENTS.md: %v", err)
	}

	tests := []struct {
		name        string
		command     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid create-readme command",
			command:     "create-readme",
			expectError: false,
		},
		{
			name:        "Valid editorconfig command",
			command:     "editorconfig",
			expectError: false,
		},
		{
			name:        "Invalid command",
			command:     "nonexistent-command",
			expectError: true,
			errorMsg:    "", // Just check for error, don't check exact message
		},
		{
			name:        "Empty command",
			command:     "",
			expectError: true, // RunAddCommand should return error for empty command
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunAddCommand(tt.command, tempDir, true) // Use dry run

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestRunAddCommandClaude(t *testing.T) {
    tempDir := t.TempDir()
    // AGENTS.md
    if err := os.WriteFile(filepath.Join(tempDir, "AGENTS.md"), []byte("# Test"), 0644); err != nil {
        t.Fatalf("failed to create AGENTS.md: %v", err)
    }
    // Enable claude
    cfg := &config.ProjectConfig{EnabledAgents: []string{"claude"}}
    if err := config.SaveProjectConfig(tempDir, cfg); err != nil {
        t.Fatalf("failed to save config: %v", err)
    }
    // Dry run should succeed
    if err := RunAddCommand("create-readme", tempDir, true); err != nil {
        t.Fatalf("dry run failed: %v", err)
    }
    // Actual
    if err := RunAddCommand("create-readme", tempDir, false); err != nil {
        t.Fatalf("RunAddCommand failed: %v", err)
    }
    // Expect .claude/commands/create-readme.md and YAML frontmatter with allowed-tools/description
    path := filepath.Join(tempDir, ".claude", "commands", "create-readme.md")
    b, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("Claude command not readable: %v", err)
    }
    s := string(b)
    if !strings.HasPrefix(s, "---") || !strings.Contains(s, "allowed-tools:") || !strings.Contains(s, "description:") {
        t.Fatalf("Claude command missing required frontmatter:\n%s", s)
    }
}

func TestRunAddCommandNotInitialized(t *testing.T) {
	// Create temporary directory without AGENTS.md
	tempDir := t.TempDir()

	err := RunAddCommand("create-readme", tempDir, true)
	if err == nil {
		t.Error("Expected error for uninitialized project")
	}

	expectedMsg := "project is not initialized with anyagent"
	if err != nil && !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateCommand(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "Valid command - create-readme",
			command:     "create-readme",
			expectError: false,
		},
		{
			name:        "Valid command - editorconfig",
			command:     "editorconfig",
			expectError: false,
		},
		{
			name:        "Invalid command",
			command:     "nonexistent",
			expectError: true,
		},
		{
			name:        "Empty command",
			command:     "",
			expectError: true,
		},
		{
			name:        "Command with invalid characters",
			command:     "test/command",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommand(tt.command)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRemoveYAMLFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Content with YAML frontmatter",
			input: `---
mode: 'agent'
description: 'Test description'
---

# Test Content

This is the main content.`,
			expected: `
# Test Content

This is the main content.`,
		},
		{
			name: "Content without YAML frontmatter",
			input: `# Test Content

This is content without frontmatter.`,
			expected: `# Test Content

This is content without frontmatter.`,
		},
		{
			name:     "Empty content",
			input:    "",
			expected: "",
		},
		{
			name: "Only YAML frontmatter",
			input: `---
mode: 'agent'
---`,
			expected: "",
		},
		{
			name: "Malformed YAML frontmatter (no end)",
			input: `---
mode: 'agent'
description: 'Test'

# Content`,
			expected: `---
mode: 'agent'
description: 'Test'

# Content`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeYAMLFrontmatter(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
