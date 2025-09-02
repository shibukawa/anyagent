package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the project-specific configuration stored in .anyagent.yaml
type ProjectConfig struct {
	ProjectName        string   `yaml:"project_name"`
	ProjectDescription string   `yaml:"project_description"`
	InstalledRules     []string `yaml:"installed_rules"`
}

// LoadProjectConfig loads the project configuration from .anyagent.yaml
func LoadProjectConfig(configPath string) (*ProjectConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return &ProjectConfig{
				InstalledRules: []string{},
			}, nil
		}
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}

	// Initialize slice if nil
	if config.InstalledRules == nil {
		config.InstalledRules = []string{}
	}

	return &config, nil
}

// Save saves the project configuration to the specified file
func (c *ProjectConfig) Save(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// SaveProjectConfig saves the project configuration to .anyagent.yaml
func SaveProjectConfig(projectDir string, config *ProjectConfig) error {
	configPath := filepath.Join(projectDir, ".anyagent.yaml")
	return config.Save(configPath)
}

// GetProjectConfigPath returns the path to .anyagent.yaml for the given project directory
func GetProjectConfigPath(projectDir string) string {
	return filepath.Join(projectDir, ".anyagent.yaml")
}

// RegenerateAgentsFile regenerates AGENTS.md with current project configuration
func (c *ProjectConfig) RegenerateAgentsFile() error {
	// Get the template
	agentsTemplate := GetAGENTSTemplate()

	// Collect rule content
	var extraRules []string
	for _, rule := range c.InstalledRules {
		content, err := getRuleTemplateContent(rule)
		if err != nil {
			return fmt.Errorf("failed to get content for rule %s: %w", rule, err)
		}
		extraRules = append(extraRules, content)
	}

	// Replace placeholder with actual rules
	extraRulesContent := strings.Join(extraRules, "\n\n")
	finalContent := strings.Replace(agentsTemplate, "{{EXTRA_RULES}}", extraRulesContent, 1)

	// Write to AGENTS.md
	return os.WriteFile("AGENTS.md", []byte(finalContent), 0644)
}

// getRuleTemplateContent gets the content for a specific rule
func getRuleTemplateContent(rule string) (string, error) {
	switch rule {
	case "go":
		return GetGoExtraRuleTemplate(), nil
	case "typescript":
		return GetTSExtraRuleTemplate(), nil
	case "docker":
		return GetDockerExtraRuleTemplate(), nil
	case "python":
		return GetPythonExtraRuleTemplate(), nil
	case "react":
		return GetReactExtraRuleTemplate(), nil
	default:
		return "", fmt.Errorf("unknown rule: %s", rule)
	}
}
