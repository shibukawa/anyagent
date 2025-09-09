package config

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Embedded template files root (deploy layout)
//
//go:embed configsrc
var templatesFS embed.FS

// (individual file embeds removed; templatesFS now holds the entire tree)

// GetUserConfigDir returns the user configuration directory for anyagent
func GetUserConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userConfigDir, "anyagent"), nil
}

// CreateUserConfigDir creates the user configuration directory
func CreateUserConfigDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// CreateTemplateStructure creates the template directory structure
func CreateTemplateStructure(baseDir string) error {
	templatesDir := filepath.Join(baseDir, "templates")
	commandsDir := filepath.Join(templatesDir, "commands")
	extraRulesDir := filepath.Join(templatesDir, "extra_rules")

	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(commandsDir, 0755); err != nil {
		return err
	}

	return os.MkdirAll(extraRulesDir, 0755)
}

// CreateTemplateFiles creates the default template files
func CreateTemplateFiles(baseDir string) error {
	return fs.WalkDir(templatesFS, "configsrc", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// create corresponding directory under baseDir for this entry
			if path == "configsrc" {
				return nil
			}
			rel := strings.TrimPrefix(path, "configsrc/")
			return os.MkdirAll(filepath.Join(baseDir, rel), 0755)
		}
		b, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}
		rel := path[len("configsrc/"):]
		outPath := filepath.Join(baseDir, rel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(outPath, b, 0644)
	})
}

// CreateTemplateFilesIfMissing creates default template files only when they don't exist.
// Existing files are preserved and not overwritten.
func CreateTemplateFilesIfMissing(baseDir string) error {
	return fs.WalkDir(templatesFS, "configsrc", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel := path[len("configsrc/"):]
		outPath := filepath.Join(baseDir, rel)
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}
			b, err := templatesFS.ReadFile(path)
			if err != nil {
				return err
			}
			if err := os.WriteFile(outPath, b, 0644); err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateAnyagentProject creates the anyagent configuration project
func CreateAnyagentProject(baseDir string) error {
	// Ensure the base directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	// Create AGENTS.md for anyagent configuration project
	agentsContent := getAnyagentAGENTSContent()
	agentsPath := filepath.Join(baseDir, "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte(agentsContent), 0644); err != nil {
		return err
	}

	// Create agent directories with symbolic links to AGENTS.md
	agentConfigs := map[string]string{
		".github":  "copilot-instructions.md",
		".amazonq": "rules/AGENTS.md",
		".claude":  "AGENTS.md",
		".junie":   "AGENTS.md",
		".gemini":  "AGENTS.md",
	}

	for agentDir, symlinkFile := range agentConfigs {
		dirPath := filepath.Join(baseDir, agentDir)

		// Create directory structure
		if agentDir == ".amazonq" {
			// Create rules subdirectory for Amazon Q
			rulesDir := filepath.Join(dirPath, "rules")
			if err := os.MkdirAll(rulesDir, 0755); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return err
			}
		}

		// Create symbolic link to AGENTS.md
		symlinkPath := filepath.Join(dirPath, symlinkFile)
		agentsRelativePath := "../AGENTS.md"
		if agentDir == ".amazonq" {
			agentsRelativePath = "../../AGENTS.md"
		}

		// Remove existing file/link if it exists
		_ = os.Remove(symlinkPath) // Ignore error if file doesn't exist

		// Create symbolic link
		if err := os.Symlink(agentsRelativePath, symlinkPath); err != nil {
			return fmt.Errorf("failed to create symbolic link %s: %w", symlinkPath, err)
		}
	}

	// Create CLAUDE.md symbolic link at project root for Claude
	claudeSymlinkPath := filepath.Join(baseDir, "CLAUDE.md")
	_ = os.Remove(claudeSymlinkPath) // Remove existing file/link if it exists
	if err := os.Symlink("AGENTS.md", claudeSymlinkPath); err != nil {
		return fmt.Errorf("failed to create CLAUDE.md symbolic link: %w", err)
	}

	return nil
}

// CheckUserConfigExists checks if the user configuration directory exists
func CheckUserConfigExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}

// Template content functions

func GetAGENTSTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/AGENTS.md.tmpl")
	return string(b)
}
func getMCPTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/mcp.yaml")
	return string(b)
}
func getGeneralCommandsTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/commands/general.md")
	return string(b)
}
func getCodingCommandsTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/commands/coding.md")
	return string(b)
}
func getProjectSpecificCommandsTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/commands/project-specific.md")
	return string(b)
}

func getAnyagentAGENTSContent() string {
	b, _ := templatesFS.ReadFile("configsrc/anyagent-AGENTS.md")
	return string(b)
}
func getReadmeTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/README.md")
	return string(b)
}
func getGoRulesTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/extra_rules/go.md")
	return string(b)
}
func getTsRulesTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/extra_rules/ts.md")
	return string(b)
}
func getDockerRulesTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/extra_rules/docker.md")
	return string(b)
}
func getPythonRulesTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/extra_rules/python.md")
	return string(b)
}
func getReactRulesTemplate() string {
	b, _ := templatesFS.ReadFile("configsrc/templates/extra_rules/react.md")
	return string(b)
}

// Public template getters for external use

// GetGoExtraRuleTemplate returns the Go extra rules template
func GetGoExtraRuleTemplate() string { return getGoRulesTemplate() }

// GetTSExtraRuleTemplate returns the TypeScript extra rules template
func GetTSExtraRuleTemplate() string { return getTsRulesTemplate() }

// GetDockerExtraRuleTemplate returns the Docker extra rules template
func GetDockerExtraRuleTemplate() string { return getDockerRulesTemplate() }

// GetPythonExtraRuleTemplate returns the Python extra rules template
func GetPythonExtraRuleTemplate() string { return getPythonRulesTemplate() }

// GetReactExtraRuleTemplate returns the React extra rules template
func GetReactExtraRuleTemplate() string { return getReactRulesTemplate() }
