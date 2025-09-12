package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunFirstSync(t *testing.T) {
	tests := []struct {
		name        string
		agentNames  []string
		expectError bool
	}{
		{
			name:        "valid single agent",
			agentNames:  []string{"copilot"},
			expectError: false,
		},
		{
			name:        "invalid agent",
			agentNames:  []string{"invalid"},
			expectError: true,
		},
		{
			name:        "multiple agents not allowed",
			agentNames:  []string{"copilot", "claude"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing (auto-clean)
			tempDir := t.TempDir()

			// Run first sync with dry-run and predefined parameters to avoid interactive prompts
			err := RunFirstSyncWithParams(tempDir, tt.agentNames, "test-project", "Test project description", true)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAgentNames(t *testing.T) {
	tests := []struct {
		name        string
		agentNames  []string
		expectError bool
		expectCount int
	}{
		{
			name:        "valid single agent",
			agentNames:  []string{"copilot"},
			expectError: false,
			expectCount: 1,
		},
		{
			name:        "multiple agents not allowed",
			agentNames:  []string{"copilot", "claude", "qdev"},
			expectError: true,
			expectCount: 0,
		},
		{
			name:        "invalid agent",
			agentNames:  []string{"invalid"},
			expectError: true,
			expectCount: 0,
		},
		{
			name:        "empty list",
			agentNames:  []string{},
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agents, err := validateAgentNames(tt.agentNames)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(agents) != tt.expectCount {
				t.Errorf("Expected %d agents, got %d", tt.expectCount, len(agents))
			}
		})
	}
}

func TestCreateAgentsFile(t *testing.T) {
	// Create a temporary directory for testing (auto-clean)
	tempDir := t.TempDir()
	var err error

	params := &InitParams{
		ProjectName:        "test-project",
		ProjectDescription: "A test project for anyagent",
		ProjectDir:         tempDir,
		DynamicParameters: map[string]string{
			"PROJECT_NAME":        "test-project",
			"PROJECT_DESCRIPTION": "A test project for anyagent",
			"PRIMARY_LANGUAGE":    "Go",
			"TEAM_NAME":           "Development Team",
		},
	}

	// Test dry run
	err = createAgentsFile(params, true)
	if err != nil {
		t.Errorf("Dry run failed: %v", err)
	}

	// Verify file doesn't exist after dry run
	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); !os.IsNotExist(err) {
		t.Errorf("AGENTS.md should not exist after dry run")
	}

	// Test actual creation
	err = createAgentsFile(params, false)
	if err != nil {
		t.Errorf("File creation failed: %v", err)
	}

	// Verify file exists and contains expected content
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		t.Errorf("AGENTS.md was not created")
	}

	content, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Errorf("Failed to read AGENTS.md: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "test-project") {
		t.Errorf("AGENTS.md does not contain project name")
	}
	if !strings.Contains(contentStr, "A test project for anyagent") {
		t.Errorf("AGENTS.md does not contain project description")
	}
}

func TestCreateAgentSymlinks(t *testing.T) {
	// Create a temporary directory for testing (auto-clean)
	tempDir := t.TempDir()
	var err error

	// Create AGENTS.md file first
	agentsPath := filepath.Join(tempDir, "AGENTS.md")
	err = os.WriteFile(agentsPath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create AGENTS.md: %v", err)
	}

	params := &InitParams{
		ProjectName:        "test-project",
		ProjectDescription: "A test project",
		ProjectDir:         tempDir,
		DynamicParameters: map[string]string{
			"PROJECT_NAME":        "test-project",
			"PROJECT_DESCRIPTION": "A test project",
			"PRIMARY_LANGUAGE":    "Go",
			"TEAM_NAME":           "Development Team",
		},
		SelectedAgents: []AIAgent{
			{
				Name:         "copilot",
				DisplayName:  "GitHub Copilot",
				ConfigPath:   "", // シンボリックリンク不要
				NeedsSymlink: false,
			},
		},
	}

	// Test dry run
	err = createAgentSymlinks(params, true)
	if err != nil {
		t.Errorf("Dry run failed: %v", err)
	}

	// Test actual creation (no-op for Copilot now)
	err = createAgentSymlinks(params, false)
	if err != nil {
		t.Errorf("createAgentSymlinks should not fail when symlink not needed: %v", err)
	}

	// Ensure that old symlink is not created (backward compatibility: if existing project still has it, that's fine, but new生成しない)
	symlinkPath := filepath.Join(tempDir, ".github", "copilot-instructions.md")
	if _, err := os.Lstat(symlinkPath); err == nil {
		t.Errorf("copilot-instructions.md symlink should no longer be created")
	}
}
