package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestUserConfigDir tests the user configuration directory functionality
func TestUserConfigDir(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "get user config directory",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserConfigDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserConfigDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("GetUserConfigDir() returned empty string")
			}
		})
	}
}

// TestCreateUserConfigDir tests creating the user configuration directory
func TestCreateUserConfigDir(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := t.TempDir()
	testConfigDir := filepath.Join(tempDir, "anyagent")

	tests := []struct {
		name    string
		dir     string
		wantErr bool
	}{
		{
			name:    "create new config directory",
			dir:     testConfigDir,
			wantErr: false,
		},
		{
			name:    "create directory in existing path",
			dir:     testConfigDir, // Same directory, should not error
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateUserConfigDir(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserConfigDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify directory was created
				if _, err := os.Stat(tt.dir); os.IsNotExist(err) {
					t.Errorf("CreateUserConfigDir() did not create directory %s", tt.dir)
				}
			}
		})
	}
}

// TestCreateTemplateStructure tests creating the template directory structure
func TestCreateTemplateStructure(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "anyagent")

	tests := []struct {
		name    string
		baseDir string
		wantErr bool
	}{
		{
			name:    "create template structure",
			baseDir: configDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateTemplateStructure(tt.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTemplateStructure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify template directories were created
				templatesDir := filepath.Join(tt.baseDir, "templates")
				commandsDir := filepath.Join(templatesDir, "commands")

				if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
					t.Errorf("CreateTemplateStructure() did not create templates directory")
				}

				if _, err := os.Stat(commandsDir); os.IsNotExist(err) {
					t.Errorf("CreateTemplateStructure() did not create commands directory")
				}
			}
		})
	}
}

// TestCreateTemplateFiles tests creating template files
func TestCreateTemplateFiles(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "anyagent")

	// First create the directory structure
	err := CreateTemplateStructure(configDir)
	if err != nil {
		t.Fatalf("Failed to create template structure: %v", err)
	}

	tests := []struct {
		name    string
		baseDir string
		wantErr bool
	}{
		{
			name:    "create template files",
			baseDir: configDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateTemplateFiles(tt.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTemplateFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify template files were created
				expectedFiles := []string{
					"templates/AGENTS.md.tmpl",
					"templates/mcp.yaml",
					"templates/commands/general.md",
					"templates/commands/coding.md",
					"templates/commands/project-specific.md",
				}

				for _, file := range expectedFiles {
					filePath := filepath.Join(tt.baseDir, file)
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						t.Errorf("CreateTemplateFiles() did not create file %s", file)
					}
				}
			}
		})
	}
}

// TestCreateAnyagentProject tests creating the anyagent configuration project
func TestCreateAnyagentProject(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "anyagent")

	tests := []struct {
		name    string
		baseDir string
		wantErr bool
	}{
		{
			name:    "create anyagent project",
			baseDir: configDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateAnyagentProject(tt.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAnyagentProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify AGENTS.md was created
				agentsFile := filepath.Join(tt.baseDir, "AGENTS.md")
				if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
					t.Errorf("CreateAnyagentProject() did not create AGENTS.md")
				}

				// Verify agent directories were created
				expectedDirs := []string{
					".amazonq",
					".claude",
					".junie",
					".gemini",
				}

				for _, dir := range expectedDirs {
					dirPath := filepath.Join(tt.baseDir, dir)
					if _, err := os.Stat(dirPath); os.IsNotExist(err) {
						t.Errorf("CreateAnyagentProject() did not create directory %s", dir)
					}
				}
			}
		})
	}
}

// TestCheckUserConfigExists tests checking if user config directory exists
func TestCheckUserConfigExists(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := filepath.Join(tempDir, "existing")
	nonExistingDir := filepath.Join(tempDir, "nonexisting")

	// Create the existing directory
	err := os.MkdirAll(existingDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name string
		dir  string
		want bool
	}{
		{
			name: "existing directory",
			dir:  existingDir,
			want: true,
		},
		{
			name: "non-existing directory",
			dir:  nonExistingDir,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckUserConfigExists(tt.dir)
			if got != tt.want {
				t.Errorf("CheckUserConfigExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
