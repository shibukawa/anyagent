package commands

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/shibukawa/anyagent/internal/config"
)

func TestRunAddRule(t *testing.T) {
	tests := []struct {
		name         string
		language     string
		createAGENTS bool
		expectError  bool
		expectedFile string
	}{
		{
			name:         "add go rules to initialized project",
			language:     "go",
			createAGENTS: true,
			expectError:  false,
			expectedFile: "go.instructions.md",
		},
		{
			name:         "add typescript rules",
			language:     "typescript",
			createAGENTS: true,
			expectError:  false,
			expectedFile: "typescript.instructions.md",
		},
		{
			name:         "add python rules with alias",
			language:     "py",
			createAGENTS: true,
			expectError:  false,
			expectedFile: "python.instructions.md",
		},
		{
			name:         "add typescript rules with js alias",
			language:     "js",
			createAGENTS: true,
			expectError:  false,
			expectedFile: "typescript.instructions.md",
		},
		{
			name:         "project not initialized",
			language:     "go",
			createAGENTS: false,
			expectError:  true,
		},
		{
			name:         "unsupported language",
			language:     "unsupported",
			createAGENTS: true,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test_add_rule")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create AGENTS.md if required
			if tt.createAGENTS {
				agentsPath := filepath.Join(tempDir, "AGENTS.md")
				err = os.WriteFile(agentsPath, []byte("# AI Agents Configuration\ntest content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create AGENTS.md: %v", err)
				}
			}

			// Test dry run first
			err = RunAddRule(tt.language, tempDir, true)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for dry run, got nil")
				}
				return // Skip the rest of the test for error cases
			}
			if err != nil {
				t.Errorf("Dry run failed: %v", err)
				return
			}

			// Test actual execution
			err = RunAddRule(tt.language, tempDir, false)
			if err != nil {
				t.Errorf("Actual run failed: %v", err)
				return
			}

			// Verify the file was created
			if tt.expectedFile != "" {
				expectedPath := filepath.Join(tempDir, ".github", "instructions", tt.expectedFile)
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("Expected file was not created: %s", expectedPath)
				}

				// Verify file content is not empty
				content, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Errorf("Failed to read created file: %v", err)
				}
				if len(content) == 0 {
					t.Errorf("Created file is empty")
				}
			}
		})
	}
}

func TestRunAddRuleCodexUpdatesAgentsOnly(t *testing.T) {
    // Create temporary directory
    tempDir := t.TempDir()

    // Create AGENTS.md to simulate initialized project
    agentsPath := filepath.Join(tempDir, "AGENTS.md")
    if err := os.WriteFile(agentsPath, []byte("# AI Agents Configuration\n"), 0644); err != nil {
        t.Fatalf("Failed to create AGENTS.md: %v", err)
    }

    // Save project config with Codex enabled
    cfg := &config.ProjectConfig{EnabledAgents: []string{"codex"}}
    if err := config.SaveProjectConfig(tempDir, cfg); err != nil {
        t.Fatalf("Failed to save project config: %v", err)
    }

    // Run add rule for Go
    if err := RunAddRule("go", tempDir, false); err != nil {
        t.Fatalf("RunAddRule failed: %v", err)
    }

    // Verify no external rule file was created
    rulePath := filepath.Join(tempDir, ".github", "instructions", "go.instructions.md")
    if _, err := os.Stat(rulePath); !os.IsNotExist(err) {
        t.Errorf("External rule file should not be created for Codex: %s", rulePath)
    }

    // Verify AGENTS.md contains the Go extra rule content
    content, err := os.ReadFile(agentsPath)
    if err != nil {
        t.Fatalf("Failed to read AGENTS.md: %v", err)
    }
    if !strings.Contains(string(content), "# Go Language Specific Rules") {
        t.Errorf("AGENTS.md was not updated with Go rules for Codex")
    }
}

func TestRunAddRuleQDevCreatesRuleFile(t *testing.T) {
    tempDir := t.TempDir()
    // AGENTS.md present
    if err := os.WriteFile(filepath.Join(tempDir, "AGENTS.md"), []byte("# AI Agents Configuration\n"), 0644); err != nil {
        t.Fatalf("failed to create AGENTS.md: %v", err)
    }
    // Set enabled agent to qdev
    cfg := &config.ProjectConfig{EnabledAgents: []string{"qdev"}}
    if err := config.SaveProjectConfig(tempDir, cfg); err != nil {
        t.Fatalf("failed to save project config: %v", err)
    }

    if err := RunAddRule("go", tempDir, false); err != nil {
        t.Fatalf("RunAddRule failed: %v", err)
    }

    // Expect .amazonq/rules/go.md
    path := filepath.Join(tempDir, ".amazonq", "rules", "go.md")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        t.Fatalf("Q Dev rule file not created: %s", path)
    }
}

func TestValidateAndNormalizeLanguage(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"go", "go", false},
		{"Go", "go", false},
		{"GO", "go", false},
		{"typescript", "typescript", false},
		{"ts", "typescript", false},
		{"python", "python", false},
		{"py", "python", false},
		{"docker", "docker", false},
		{"react", "react", false},
		{"javascript", "typescript", false}, // alias
		{"js", "typescript", false},         // alias
		{"golang", "go", false},             // alias
		{"unsupported", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := validateAndNormalizeLanguage(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %s, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s for input %s", tt.expected, result, tt.input)
			}
		})
	}
}

func TestCreateInstructionsDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_instructions_dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	instructionsDir := filepath.Join(tempDir, ".github", "instructions")

	// Test dry run
	err = createInstructionsDirectory(instructionsDir, true)
	if err != nil {
		t.Errorf("Dry run failed: %v", err)
	}

	// Verify directory doesn't exist after dry run
	if _, err := os.Stat(instructionsDir); !os.IsNotExist(err) {
		t.Errorf("Directory should not exist after dry run")
	}

	// Test actual creation
	err = createInstructionsDirectory(instructionsDir, false)
	if err != nil {
		t.Errorf("Directory creation failed: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(instructionsDir); os.IsNotExist(err) {
		t.Errorf("Directory was not created")
	}
}

func TestCreateRuleFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_rule_file")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the directory first
	instructionsDir := filepath.Join(tempDir, ".github", "instructions")
	err = os.MkdirAll(instructionsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create instructions directory: %v", err)
	}

	ruleFilePath := filepath.Join(instructionsDir, "go.instructions.md")
	testContent := "# Go Rules\nTest content for Go rules"

	// Test dry run
	err = createRuleFile(ruleFilePath, testContent, true)
	if err != nil {
		t.Errorf("Dry run failed: %v", err)
	}

	// Verify file doesn't exist after dry run
	if _, err := os.Stat(ruleFilePath); !os.IsNotExist(err) {
		t.Errorf("File should not exist after dry run")
	}

	// Test actual creation
	err = createRuleFile(ruleFilePath, testContent, false)
	if err != nil {
		t.Errorf("File creation failed: %v", err)
	}

	// Verify file exists and has correct content
	if _, err := os.Stat(ruleFilePath); os.IsNotExist(err) {
		t.Errorf("File was not created")
	}

	content, err := os.ReadFile(ruleFilePath)
	if err != nil {
		t.Errorf("Failed to read created file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", testContent, string(content))
	}
}
