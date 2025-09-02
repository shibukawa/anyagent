package config

import (
	"strings"
	"testing"
)

func TestGetCommandTemplate(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectError bool
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
			command:     "nonexistent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := GetCommandTemplate(tt.command)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if content == "" {
					t.Error("Expected content but got empty string")
				}
				// Check if content starts with YAML frontmatter
				if !strings.HasPrefix(content, "---") {
					t.Error("Expected content to start with YAML frontmatter")
				}
			}
		})
	}
}

func TestGetAvailableCommands(t *testing.T) {
	commands, err := GetAvailableCommands()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(commands) == 0 {
		t.Error("Expected at least one command")
	}

	// Check for expected commands
	expectedCommands := []string{"create-readme", "editorconfig"}
	for _, expected := range expectedCommands {
		found := false
		for _, command := range commands {
			if command == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command '%s' not found in available commands: %v", expected, commands)
		}
	}
}
