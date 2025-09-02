package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

// Embedded template files
//
//go:embed templates/AGENTS.md.tmpl
var agentsTemplate string

//go:embed templates/mcp.yaml
var mcpTemplate string

//go:embed templates/general.md
var generalCommandsTemplate string

//go:embed templates/coding.md
var codingCommandsTemplate string

//go:embed templates/project-specific.md
var projectSpecificCommandsTemplate string

//go:embed templates/anyagent-AGENTS.md
var anyagentAGENTSContent string

//go:embed templates/README.md
var readmeTemplate string

//go:embed templates/extra_rules/go.md
var goRulesTemplate string

//go:embed templates/extra_rules/ts.md
var tsRulesTemplate string

//go:embed templates/extra_rules/docker.md
var dockerRulesTemplate string

//go:embed templates/extra_rules/python.md
var pythonRulesTemplate string

//go:embed templates/extra_rules/react.md
var reactRulesTemplate string

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
	templateFiles := map[string]string{
		"templates/AGENTS.md.tmpl":               GetAGENTSTemplate(),
		"templates/mcp.yaml":                     getMCPTemplate(),
		"templates/commands/general.md":          getGeneralCommandsTemplate(),
		"templates/commands/coding.md":           getCodingCommandsTemplate(),
		"templates/commands/project-specific.md": getProjectSpecificCommandsTemplate(),
		"templates/extra_rules/go.md":            getGoRulesTemplate(),
		"templates/extra_rules/ts.md":            getTsRulesTemplate(),
		"templates/extra_rules/docker.md":        getDockerRulesTemplate(),
		"templates/extra_rules/python.md":        getPythonRulesTemplate(),
		"templates/extra_rules/react.md":         getReactRulesTemplate(),
		"README.md":                              getReadmeTemplate(),
	}

	for filePath, content := range templateFiles {
		fullPath := filepath.Join(baseDir, filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
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
	return agentsTemplate
}

func getMCPTemplate() string {
	return mcpTemplate
}

func getGeneralCommandsTemplate() string {
	return generalCommandsTemplate
}

func getCodingCommandsTemplate() string {
	return codingCommandsTemplate
}

func getProjectSpecificCommandsTemplate() string {
	return projectSpecificCommandsTemplate
}

func getAnyagentAGENTSContent() string {
	return anyagentAGENTSContent
}

func getReadmeTemplate() string {
	return readmeTemplate
}

func getGoRulesTemplate() string {
	return goRulesTemplate
}

func getTsRulesTemplate() string {
	return tsRulesTemplate
}

func getDockerRulesTemplate() string {
	return dockerRulesTemplate
}

func getPythonRulesTemplate() string {
	return pythonRulesTemplate
}

func getReactRulesTemplate() string {
	return reactRulesTemplate
}

// Public template getters for external use

// GetGoExtraRuleTemplate returns the Go extra rules template
func GetGoExtraRuleTemplate() string {
	return goRulesTemplate
}

// GetTSExtraRuleTemplate returns the TypeScript extra rules template
func GetTSExtraRuleTemplate() string {
	return tsRulesTemplate
}

// GetDockerExtraRuleTemplate returns the Docker extra rules template
func GetDockerExtraRuleTemplate() string {
	return dockerRulesTemplate
}

// GetPythonExtraRuleTemplate returns the Python extra rules template
func GetPythonExtraRuleTemplate() string {
	return pythonRulesTemplate
}

// GetReactExtraRuleTemplate returns the React extra rules template
func GetReactExtraRuleTemplate() string {
	return reactRulesTemplate
}
