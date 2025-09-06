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
	ProjectName        string            `yaml:"project_name"`
	ProjectDescription string            `yaml:"project_description"`
	InstalledRules     []string          `yaml:"installed_rules"`
	InstalledCommands  []string          `yaml:"installed_commands"`
	EnabledAgents      []string          `yaml:"enabled_agents"`
	Parameters         map[string]string `yaml:"parameters"`
	MCPServers         map[string]string `yaml:"mcp_servers"`
}

// LoadProjectConfig loads the project configuration from .anyagent.yaml
func LoadProjectConfig(configPath string) (*ProjectConfig, error) {
	// Primary: read from provided path
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Back-compat: if provided path is .anyagent/config.yaml, try legacy .anyagent.yaml
			if filepath.Base(configPath) == "config.yaml" {
				legacy := filepath.Join(filepath.Dir(filepath.Dir(configPath)), ".anyagent.yaml")
				if b, e := os.ReadFile(legacy); e == nil {
					data = b
				} else if os.IsNotExist(e) {
					// Return default config if neither exists
					return &ProjectConfig{InstalledRules: []string{}}, nil
				} else {
					return nil, fmt.Errorf("failed to read legacy project config: %w", e)
				}
			} else {
				// Return default config if file doesn't exist
				return &ProjectConfig{InstalledRules: []string{}}, nil
			}
		} else {
			return nil, fmt.Errorf("failed to read project config: %w", err)
		}
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}

	// Initialize slice if nil
	if config.InstalledRules == nil {
		config.InstalledRules = []string{}
	}
	if config.InstalledCommands == nil {
		config.InstalledCommands = []string{}
	}
	if config.EnabledAgents == nil {
		config.EnabledAgents = []string{}
	}
	if config.Parameters == nil {
		config.Parameters = map[string]string{}
	}
	if config.MCPServers == nil {
		config.MCPServers = map[string]string{}
	}

	return &config, nil
}

// Save saves the project configuration to the specified file
func (c *ProjectConfig) Save(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return os.WriteFile(configPath, data, 0644)
}

// SaveProjectConfig saves the project configuration to .anyagent.yaml
func SaveProjectConfig(projectDir string, config *ProjectConfig) error {
	configPath := GetProjectConfigPath(projectDir)
	return config.Save(configPath)
}

// GetProjectConfigPath returns the path to .anyagent.yaml for the given project directory
func GetProjectConfigPath(projectDir string) string {
	return filepath.Join(projectDir, ".anyagent", "config.yaml")
}

// RegenerateAgentsFile regenerates AGENTS.md with current project configuration
func (c *ProjectConfig) RegenerateAgentsFile() error {
	// Get the template
	agentsTemplate := GetAGENTSTemplate()
	// Prefer project .anyagent/AGENTS.md.tmpl if present
	if b, err := os.ReadFile(filepath.Join(".anyagent", "AGENTS.md.tmpl")); err == nil {
		agentsTemplate = string(b)
	}

	// Prepare parameters map ensuring required keys are present
	params := map[string]string{}
	for k, v := range c.Parameters {
		params[k] = v
	}
	if c.ProjectName != "" {
		params["PROJECT_NAME"] = c.ProjectName
	}
	if c.ProjectDescription != "" {
		params["PROJECT_DESCRIPTION"] = c.ProjectDescription
	}

	// Replace placeholders with parameters
	content := ReplaceTemplateParameters(agentsTemplate, params)

	// Collect and inject extra rule content
	var extraRules []string
	for _, rule := range c.InstalledRules {
		contentRule, err := getRuleTemplateContent(rule)
		if err != nil {
			return fmt.Errorf("failed to get content for rule %s: %w", rule, err)
		}
		extraRules = append(extraRules, contentRule)
	}
	extraRulesContent := strings.Join(extraRules, "\n\n")
	content = strings.Replace(content, "{{EXTRA_RULES}}", extraRulesContent, 1)

	// Write to AGENTS.md in current working directory
	return os.WriteFile("AGENTS.md", []byte(content), 0644)
}

// RegenerateAgentsFileAt regenerates AGENTS.md at the specified project directory
func (c *ProjectConfig) RegenerateAgentsFileAt(projectDir string) error {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	if projectDir != "" {
		_ = os.Chdir(projectDir)
	}
	return c.RegenerateAgentsFile()
}

// getRuleTemplateContent gets the content for a specific rule
func getRuleTemplateContent(rule string) (string, error) {
	// Map rule to filename
	filename := ""
	switch rule {
	case "go":
		filename = "go.md"
	case "typescript":
		filename = "ts.md"
	case "docker":
		filename = "docker.md"
	case "python":
		filename = "python.md"
	case "react":
		filename = "react.md"
	default:
		return "", fmt.Errorf("unknown rule: %s", rule)
	}
	// Prefer project-local .anyagent/extra_rules/<file>
	if b, err := os.ReadFile(filepath.Join(".anyagent", "extra_rules", filename)); err == nil {
		return string(b), nil
	}
	// Fallback to embedded templates
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
	}
	return "", fmt.Errorf("unknown rule: %s", rule)
}
