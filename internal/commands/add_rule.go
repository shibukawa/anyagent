package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/anyagent/internal/config"
)

// SupportedRules defines the available extra rules
var SupportedRules = []string{
	"go",
	"typescript", "ts",
	"docker",
	"python", "py",
	"react",
}

// RunAddRule executes the add rule command functionality
func RunAddRule(language, projectDir string, dryRun bool) error {
	fmt.Printf("Adding %s rules to project...\n", language)

	// Get project directory (current directory if not specified)
	if projectDir == "" {
		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Make sure the project directory exists
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectDir)
	}

	fmt.Printf("Project directory: %s\n", projectDir)

	// Check if project is initialized (has AGENTS.md)
	agentsPath := filepath.Join(projectDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		return fmt.Errorf("project is not initialized with anyagent. Run 'anyagent sync' first")
	}

	// Validate and normalize language name
	normalizedLanguage, err := validateAndNormalizeLanguage(language)
	if err != nil {
		return fmt.Errorf("unsupported language: %s. Supported: %s", language, strings.Join(SupportedRules, ", "))
	}

	// Get the rule template content
	ruleContent, err := getRuleTemplate(normalizedLanguage)
	if err != nil {
		return fmt.Errorf("failed to get rule template: %w", err)
	}

	// Create external rule files for agent-specific locations
	if shouldCreateCopilotRuleFiles(projectDir) {
		// Create GitHub Copilot instructions directory
		copilotInstructionsDir := filepath.Join(projectDir, ".github", "instructions")
		if err := createInstructionsDirectory(copilotInstructionsDir, dryRun); err != nil {
			return fmt.Errorf("failed to create instructions directory: %w", err)
		}

		// Create the rule file
		ruleFilePath := filepath.Join(copilotInstructionsDir, fmt.Sprintf("%s.instructions.md", normalizedLanguage))
		if err := createRuleFile(ruleFilePath, ruleContent, dryRun); err != nil {
			return fmt.Errorf("failed to create rule file: %w", err)
		}
	} else if shouldCreateQDevRuleFiles(projectDir) {
		// Amazon Q Developer rules directory
		qdevRulesDir := filepath.Join(projectDir, ".amazonq", "rules")
		if err := createInstructionsDirectory(qdevRulesDir, dryRun); err != nil {
			return fmt.Errorf("failed to create Q Developer rules directory: %w", err)
		}
		qdevRulePath := filepath.Join(qdevRulesDir, fmt.Sprintf("%s.md", normalizedLanguage))
		if err := createRuleFile(qdevRulePath, ruleContent, dryRun); err != nil {
			return fmt.Errorf("failed to create Q Developer rule file: %w", err)
		}
		fmt.Printf("📄 Amazon Q Developer rule created: .amazonq/rules/%s.md\n", normalizedLanguage)
	} else {
		fmt.Printf("ℹ️  Codex selected: skipping external rule files; regenerating AGENTS.md only.\n")
	}

	// Update project configuration and regenerate AGENTS.md
	if !dryRun {
		if err := updateProjectConfigAndRegenerate(projectDir, normalizedLanguage); err != nil {
			fmt.Printf("⚠️  Warning: Failed to update configuration: %v\n", err)
		}
	}

	fmt.Printf("✅ %s rules added successfully\n", normalizedLanguage)
	return nil
}

// validateAndNormalizeLanguage validates the language and returns the normalized name
func validateAndNormalizeLanguage(language string) (string, error) {
	language = strings.ToLower(language)

	// Language aliases
	aliases := map[string]string{
		"ts":         "typescript",
		"py":         "python",
		"golang":     "go",
		"javascript": "typescript", // Use typescript rules for JavaScript
		"js":         "typescript",
	}

	// Check if it's an alias
	if normalized, exists := aliases[language]; exists {
		language = normalized
	}

	// Check if the language is supported
	for _, supported := range SupportedRules {
		if language == supported {
			return language, nil
		}
	}

	return "", fmt.Errorf("unsupported language: %s", language)
}

// getRuleTemplate retrieves the template content for the specified language
func getRuleTemplate(language string) (string, error) {
	switch language {
	case "go":
		return config.GetGoExtraRuleTemplate(), nil
	case "typescript":
		return config.GetTSExtraRuleTemplate(), nil
	case "docker":
		return config.GetDockerExtraRuleTemplate(), nil
	case "python":
		return config.GetPythonExtraRuleTemplate(), nil
	case "react":
		return config.GetReactExtraRuleTemplate(), nil
	default:
		return "", fmt.Errorf("template not found for language: %s", language)
	}
}

// createInstructionsDirectory creates the .github/instructions directory
func createInstructionsDirectory(dir string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would create directory: %s\n", dir)
		return nil
	}

	fmt.Printf("📁 Creating instructions directory: %s\n", dir)
	return os.MkdirAll(dir, 0755)
}

// createRuleFile creates the rule instruction file
func createRuleFile(filePath, content string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would create rule file: %s\n", filePath)
		fmt.Printf("[DRY RUN] Content preview:\n")
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if i >= 10 {
				fmt.Printf("  ... (truncated)\n")
				break
			}
			fmt.Printf("  %s\n", line)
		}
		return nil
	}

	fmt.Printf("📄 Creating rule file: %s\n", filePath)
	return os.WriteFile(filePath, []byte(content), 0644)
}

// updateProjectConfigAndRegenerate updates the project config and regenerates AGENTS.md
func updateProjectConfigAndRegenerate(projectDir, rule string) error {
	configPath := config.GetProjectConfigPath(projectDir)

	// Load existing config or create new one
	projectConfig, err := config.LoadProjectConfig(configPath)
	if err != nil {
		// If config doesn't exist, create a new one
		projectConfig = &config.ProjectConfig{
			InstalledRules: []string{},
		}
	}

	// Check if rule is already installed
	for _, installedRule := range projectConfig.InstalledRules {
		if installedRule == rule {
			// Already installed, just regenerate
			return projectConfig.RegenerateAgentsFile()
		}
	}

	// Add the rule to installed rules
	projectConfig.InstalledRules = append(projectConfig.InstalledRules, rule)

	// Save the updated config
	if err := projectConfig.Save(configPath); err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	// Regenerate AGENTS.md at the specified project directory
	return projectConfig.RegenerateAgentsFileAt(projectDir)
}
