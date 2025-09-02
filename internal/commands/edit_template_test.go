package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/anyagent/internal/config"
)

// TestEditTemplateCommand tests the edit-template command functionality
func TestEditTemplateCommand(t *testing.T) {
	tempDir := t.TempDir()
	testConfigDir := filepath.Join(tempDir, "anyagent")

	tests := []struct {
		name          string
		configDir     string
		setupExisting bool
		wantErr       bool
	}{
		{
			name:          "create new template environment",
			configDir:     testConfigDir,
			setupExisting: false,
			wantErr:       false,
		},
		{
			name:          "existing template environment",
			configDir:     testConfigDir,
			setupExisting: true,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupExisting {
				// Pre-create the configuration directory for this test
				err := setupExistingConfig(tt.configDir)
				if err != nil {
					t.Fatalf("Failed to setup existing config: %v", err)
				}
			}

			err := RunEditTemplate(tt.configDir, true, false) // dryRun = true, hardReset = false for testing
			if (err != nil) != tt.wantErr {
				t.Errorf("RunEditTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the configuration was set up correctly
				verifyTemplateEnvironment(t, tt.configDir)
			}
		})
	}
}

// TestEditTemplateWithoutVSCode tests the edit-template command without VSCode launch
func TestEditTemplateWithoutVSCode(t *testing.T) {
	tempDir := t.TempDir()
	testConfigDir := filepath.Join(tempDir, "anyagent")

	err := RunEditTemplate(testConfigDir, true, false) // dryRun = true, hardReset = false
	if err != nil {
		t.Errorf("RunEditTemplate() with dryRun failed: %v", err)
	}

	// Verify all necessary components were created
	verifyTemplateEnvironment(t, testConfigDir)
}

// TestEditTemplateHardReset tests the edit-template command with hard reset
func TestEditTemplateHardReset(t *testing.T) {
	tempDir := t.TempDir()
	testConfigDir := filepath.Join(tempDir, "anyagent")

	// First, create an existing environment
	err := setupCompleteConfig(testConfigDir)
	if err != nil {
		t.Fatalf("Failed to setup initial config: %v", err)
	}

	// Modify one of the template files to test hard reset
	templatesDir := filepath.Join(testConfigDir, "templates")
	modifiedTemplate := filepath.Join(templatesDir, "AGENTS.md.tmpl")
	err = os.WriteFile(modifiedTemplate, []byte("# Modified template"), 0644)
	if err != nil {
		t.Fatalf("Failed to modify template: %v", err)
	}

	// Perform hard reset
	err = RunEditTemplate(testConfigDir, true, true) // dryRun = true, hardReset = true
	if err != nil {
		t.Errorf("RunEditTemplate() with hard reset failed: %v", err)
	}

	// Verify the template was reset to original
	content, err := os.ReadFile(modifiedTemplate)
	if err != nil {
		t.Fatalf("Failed to read template after reset: %v", err)
	}

	// Check that the content is no longer the modified version
	if strings.Contains(string(content), "# Modified template") {
		t.Errorf("Hard reset did not restore original template content")
	}

	// Verify all components are still present
	verifyTemplateEnvironment(t, testConfigDir)
}

// TestValidateTemplateEnvironment tests template environment validation
func TestValidateTemplateEnvironment(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		setupFunc func(string) error
		configDir string
		wantValid bool
	}{
		{
			name: "valid template environment",
			setupFunc: func(dir string) error {
				return setupCompleteConfig(dir)
			},
			configDir: filepath.Join(tempDir, "valid"),
			wantValid: true,
		},
		{
			name: "missing template files",
			setupFunc: func(dir string) error {
				return config.CreateUserConfigDir(dir) // Only create directory
			},
			configDir: filepath.Join(tempDir, "incomplete"),
			wantValid: false,
		},
		{
			name: "non-existent directory",
			setupFunc: func(dir string) error {
				return nil // Don't create anything
			},
			configDir: filepath.Join(tempDir, "nonexistent"),
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupFunc(tt.configDir)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			valid := ValidateTemplateEnvironment(tt.configDir)
			if valid != tt.wantValid {
				t.Errorf("ValidateTemplateEnvironment() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

// TestLaunchVSCode tests VSCode launching functionality (dry run)
func TestLaunchVSCode(t *testing.T) {
	tempDir := t.TempDir()
	testConfigDir := filepath.Join(tempDir, "anyagent")

	// Setup a complete configuration
	err := setupCompleteConfig(testConfigDir)
	if err != nil {
		t.Fatalf("Failed to setup config: %v", err)
	}

	tests := []struct {
		name      string
		configDir string
		dryRun    bool
		wantErr   bool
	}{
		{
			name:      "launch VSCode dry run",
			configDir: testConfigDir,
			dryRun:    true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LaunchVSCode(tt.configDir, tt.dryRun)
			if (err != nil) != tt.wantErr {
				t.Errorf("LaunchVSCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper functions

func setupExistingConfig(configDir string) error {
	if err := config.CreateUserConfigDir(configDir); err != nil {
		return err
	}

	if err := config.CreateTemplateStructure(configDir); err != nil {
		return err
	}

	// Create some existing template files
	templatesDir := filepath.Join(configDir, "templates")
	existingTemplate := filepath.Join(templatesDir, "AGENTS.md.tmpl")
	return os.WriteFile(existingTemplate, []byte("# Existing template"), 0644)
}

func setupCompleteConfig(configDir string) error {
	if err := config.CreateUserConfigDir(configDir); err != nil {
		return err
	}

	if err := config.CreateTemplateStructure(configDir); err != nil {
		return err
	}

	if err := config.CreateTemplateFiles(configDir); err != nil {
		return err
	}

	return config.CreateAnyagentProject(configDir)
}

func verifyTemplateEnvironment(t *testing.T, configDir string) {
	t.Helper()

	// Check that the config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %s", configDir)
		return
	}

	// Check template structure
	templatesDir := filepath.Join(configDir, "templates")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Errorf("Templates directory was not created: %s", templatesDir)
	}

	commandsDir := filepath.Join(templatesDir, "commands")
	if _, err := os.Stat(commandsDir); os.IsNotExist(err) {
		t.Errorf("Commands directory was not created: %s", commandsDir)
	}

	extraRulesDir := filepath.Join(templatesDir, "extra_rules")
	if _, err := os.Stat(extraRulesDir); os.IsNotExist(err) {
		t.Errorf("Extra rules directory was not created: %s", extraRulesDir)
	}

	// Check template files
	expectedFiles := []string{
		"templates/AGENTS.md.tmpl",
		"templates/mcp.yaml",
		"templates/commands/general.md",
		"templates/commands/coding.md",
		"templates/commands/project-specific.md",
		"templates/extra_rules/go.md",
		"templates/extra_rules/ts.md",
		"templates/extra_rules/docker.md",
		"templates/extra_rules/python.md",
		"templates/extra_rules/react.md",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(configDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Template file was not created: %s", file)
		}
	}

	// Check AGENTS.md for anyagent project
	agentsFile := filepath.Join(configDir, "AGENTS.md")
	if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
		t.Errorf("AGENTS.md was not created: %s", agentsFile)
	} else {
		// Verify content
		content, err := os.ReadFile(agentsFile)
		if err != nil {
			t.Errorf("Failed to read AGENTS.md: %v", err)
		} else if !strings.Contains(string(content), "anyagent") {
			t.Errorf("AGENTS.md does not contain expected content")
		}
	}

	// Check agent directories
	expectedDirs := []string{".github", ".amazonq", ".claude", ".junie", ".gemini"}
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(configDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Agent directory was not created: %s", dir)
		}
	}

	// Check symbolic links
	expectedSymlinks := map[string]string{
		".github/copilot-instructions.md": "../AGENTS.md",
		".amazonq/rules/AGENTS.md":        "../../AGENTS.md",
		".claude/AGENTS.md":               "../AGENTS.md",
		".junie/AGENTS.md":                "../AGENTS.md",
		".gemini/AGENTS.md":               "../AGENTS.md",
		"CLAUDE.md":                       "AGENTS.md", // Project root CLAUDE.md
	}

	for symlinkPath, expectedTarget := range expectedSymlinks {
		fullPath := filepath.Join(configDir, symlinkPath)
		if target, err := os.Readlink(fullPath); err != nil {
			t.Errorf("Symbolic link was not created or is not a symlink: %s, error: %v", symlinkPath, err)
		} else if target != expectedTarget {
			t.Errorf("Symbolic link %s has wrong target: got %s, want %s", symlinkPath, target, expectedTarget)
		}
	}
}
