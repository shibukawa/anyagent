package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/anyagent/internal/config"
)

// RunRemoveRule executes the remove rule command functionality
func RunRemoveRule(language, projectDir string, dryRun bool) error {
	fmt.Printf("Removing %s rules from project...\n", language)

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

	// When Copilot is enabled, remove external rule file. For Codex-only projects
	// there is no external rule file, so skip this step.
	if shouldCreateCopilotRuleFiles(projectDir) {
		// Check if rule file exists
		ruleFilePath := filepath.Join(projectDir, ".github", "instructions", fmt.Sprintf("%s.instructions.md", normalizedLanguage))
		if _, err := os.Stat(ruleFilePath); os.IsNotExist(err) {
			return fmt.Errorf("rule file does not exist: %s", ruleFilePath)
		}

		// Remove the rule file
		if err := removeRuleFile(ruleFilePath, dryRun); err != nil {
			return fmt.Errorf("failed to remove rule file: %w", err)
		}
	} else if shouldCreateQDevRuleFiles(projectDir) {
		// Remove Q Developer rule file if present
		qdevRulePath := filepath.Join(projectDir, ".amazonq", "rules", fmt.Sprintf("%s.md", normalizedLanguage))
		if _, err := os.Stat(qdevRulePath); os.IsNotExist(err) {
			return fmt.Errorf("rule file does not exist: %s", qdevRulePath)
		}
		if err := removeRuleFile(qdevRulePath, dryRun); err != nil {
			return fmt.Errorf("failed to remove rule file: %w", err)
		}
	} else {
		fmt.Printf("ℹ️  Codex selected: no external rule files to remove; updating AGENTS.md only.\n")
	}

	// Update project configuration and regenerate AGENTS.md
	if !dryRun {
		if err := removeFromProjectConfigAndRegenerate(projectDir, normalizedLanguage); err != nil {
			fmt.Printf("⚠️  Warning: Failed to update configuration: %v\n", err)
		}
	}

	fmt.Printf("✅ %s rules removed successfully\n", normalizedLanguage)
	return nil
}

// RunListRules executes the list rules command functionality
func RunListRules(projectDir string) error {
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

	fmt.Printf("Rule status for project: %s\n\n", projectDir)

	// Check if project is initialized
	agentsPath := filepath.Join(projectDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		fmt.Println("❌ Project is not initialized with anyagent")
		return nil
	}

	instructionsDir := filepath.Join(projectDir, ".github", "instructions")

	// Load installed rules from config and detect selected agent
	cfg, _ := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	installed := map[string]bool{}
	for _, r := range cfg.InstalledRules {
		installed[r] = true
	}
	agent := selectedAgent(projectDir)

	// Check each supported rule
	fmt.Println("Available rules:")
	installedCount := 0
	for _, rule := range SupportedRules {
		isInstalled := installed[rule]
		// If not recorded in config, consider Copilot file presence
		if !isInstalled {
			if _, err := os.Stat(filepath.Join(instructionsDir, fmt.Sprintf("%s.instructions.md", rule))); err == nil {
				isInstalled = true
			}
		}

		if isInstalled {
			// Add hint for Codex listing
			if agent == "codex" {
				fmt.Printf("  ✅ %s (installed in AGENTS.md)\n", rule)
			} else {
				fmt.Printf("  ✅ %s (installed)\n", rule)
			}
			installedCount++
		} else {
			fmt.Printf("  ⬜ %s (not installed)\n", rule)
		}
	}

	fmt.Printf("\nSummary: %d/%d rules installed\n", installedCount, len(SupportedRules))

	if installedCount == 0 {
		fmt.Println("\n💡 Use 'anyagent add rule <language>' to install rules")
	}

	return nil
}

// removeRuleFile removes the rule instruction file
func removeRuleFile(filePath string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would remove rule file: %s\n", filePath)
		return nil
	}

	fmt.Printf("🗑️  Removing rule file: %s\n", filePath)
	return os.Remove(filePath)
}

// removeFromProjectConfigAndRegenerate removes a rule from project config and regenerates AGENTS.md
func removeFromProjectConfigAndRegenerate(projectDir, rule string) error {
	configPath := config.GetProjectConfigPath(projectDir)

	// Load existing config
	projectConfig, err := config.LoadProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Remove the rule from installed rules
	newRules := []string{}
	found := false
	for _, installedRule := range projectConfig.InstalledRules {
		if installedRule != rule {
			newRules = append(newRules, installedRule)
		} else {
			found = true
		}
	}

	if !found {
		// Rule wasn't in config, just regenerate at the project directory
		return projectConfig.RegenerateAgentsFileAt(projectDir)
	}

	// Update the rules list
	projectConfig.InstalledRules = newRules

	// Save the updated config
	if err := projectConfig.Save(configPath); err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	// Regenerate AGENTS.md at the specified project directory
	return projectConfig.RegenerateAgentsFileAt(projectDir)
}
