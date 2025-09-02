package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRemoveCommand(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create AGENTS.md to simulate initialized project
	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("# Test Project"), 0644); err != nil {
		t.Fatalf("Failed to create AGENTS.md: %v", err)
	}

	// Create prompts directory and command file
	promptsDir := filepath.Join(tempDir, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	commandFilePath := filepath.Join(promptsDir, "create-readme.prompt.md")
	if err := os.WriteFile(commandFilePath, []byte("# Create README"), 0644); err != nil {
		t.Fatalf("Failed to create command file: %v", err)
	}

	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "Valid command removal",
			command:     "create-readme",
			expectError: false,
		},
		{
			name:        "Non-existent command",
			command:     "non-existent",
			expectError: true,
		},
		{
			name:        "Empty command",
			command:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunRemoveCommand(tt.command, tempDir, true) // Use dry run

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRunListCommands(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create AGENTS.md to simulate initialized project
	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("# Test Project"), 0644); err != nil {
		t.Fatalf("Failed to create AGENTS.md: %v", err)
	}

	// Create prompts directory and some command files
	promptsDir := filepath.Join(tempDir, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	// Create create-readme.prompt.md
	commandFilePath := filepath.Join(promptsDir, "create-readme.prompt.md")
	if err := os.WriteFile(commandFilePath, []byte("# Create README"), 0644); err != nil {
		t.Fatalf("Failed to create command file: %v", err)
	}

	err := RunListCommands(tempDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestRunListCommandsNotInitialized(t *testing.T) {
	// Create temporary directory without AGENTS.md
	tempDir := t.TempDir()

	err := RunListCommands(tempDir)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
