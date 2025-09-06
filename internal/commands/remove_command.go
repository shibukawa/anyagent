package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/anyagent/internal/config"
)

// RunRemoveCommand executes the remove command functionality
func RunRemoveCommand(command, projectDir string, dryRun bool) error {
	fmt.Printf("Removing %s command from project...\n", command)

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

	// Validate command name
	if command == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	// Check if Copilot command file exists
	copilotCommandFilePath := filepath.Join(projectDir, ".github", "prompts", fmt.Sprintf("%s.prompt.md", command))
	copilotExists := false
	if _, err := os.Stat(copilotCommandFilePath); err == nil {
		copilotExists = true
	}

	// Check if Amazon Q Developer command file exists
	qdevExists := false
	qdevCommandFilePath := ""
	homeDir, err := os.UserHomeDir()
	if err == nil {
		qdevCommandName := strings.ReplaceAll(strings.ReplaceAll(command, "-", " "), "_", " ")
		qdevCommandFilePath = filepath.Join(homeDir, ".aws", "amazonq", "prompts", fmt.Sprintf("%s.md", qdevCommandName))
		if _, err := os.Stat(qdevCommandFilePath); err == nil {
			qdevExists = true
		}
	}

	// Check if Codex command file exists
	codexExists := false
	codexCommandFilePath := ""
	if homeDir == "" {
		homeDir, _ = os.UserHomeDir()
	}
	if homeDir != "" {
		codexCommandFilePath = filepath.Join(homeDir, ".codex", "prompts", fmt.Sprintf("%s.md", command))
		if _, err := os.Stat(codexCommandFilePath); err == nil {
			codexExists = true
		}
	}

	// Check if Claude command file exists
	claudeExists := false
	claudeCommandFilePath := filepath.Join(projectDir, ".claude", "commands", fmt.Sprintf("%s.md", command))
	if _, err := os.Stat(claudeCommandFilePath); err == nil {
		claudeExists = true
	}

	// Check if Gemini command file exists
	geminiExists := false
	geminiCommandFilePath := filepath.Join(projectDir, ".gemini", "commands", fmt.Sprintf("%s.toml", command))
	if _, err := os.Stat(geminiCommandFilePath); err == nil {
		geminiExists = true
	}

	if !copilotExists && !qdevExists && !codexExists && !claudeExists && !geminiExists {
		return fmt.Errorf("command '%s' is not installed", command)
	}

	// Remove Copilot command file
	if copilotExists {
		if err := removeCommandFile(copilotCommandFilePath, "VS Code Copilot", dryRun); err != nil {
			return fmt.Errorf("failed to remove VS Code Copilot command file: %w", err)
		}
	}

	// Remove Amazon Q Developer command file
	if qdevExists {
		if err := removeCommandFile(qdevCommandFilePath, "Amazon Q Developer", dryRun); err != nil {
			fmt.Printf("⚠️  Warning: Could not remove Amazon Q Developer command file: %v\n", err)
		}
	}

	// Remove Codex command file
	if codexExists {
		if err := removeCommandFile(codexCommandFilePath, "Codex", dryRun); err != nil {
			fmt.Printf("⚠️  Warning: Could not remove Codex command file: %v\n", err)
		}
	}

	// Remove Claude command file
	if claudeExists {
		if err := removeCommandFile(claudeCommandFilePath, "Claude Code", dryRun); err != nil {
			fmt.Printf("⚠️  Warning: Could not remove Claude command file: %v\n", err)
		}
	}

	// Remove Gemini command file
	if geminiExists {
		if err := removeCommandFile(geminiCommandFilePath, "Gemini Code", dryRun); err != nil {
			fmt.Printf("⚠️  Warning: Could not remove Gemini command file: %v\n", err)
		}
	}

	fmt.Printf("✅ %s command removed successfully\n", command)

	// Update project config to remove the command from installed_commands
	if err := removeInstalledCommandFromConfig(projectDir, command, dryRun); err != nil {
		fmt.Printf("⚠️  Warning: Failed to update project config when removing '%s': %v\n", command, err)
	}
	return nil
}

// RunListCommands executes the list commands command functionality
func RunListCommands(projectDir string) error {
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

	fmt.Printf("Command status for project: %s\n\n", projectDir)

	// Check if project is initialized
	agentsPath := filepath.Join(projectDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		fmt.Println("❌ Project is not initialized with anyagent")
		return nil
	}

	// Get available commands from templates
	availableCommands, err := config.GetAvailableCommands()
	if err != nil {
		return fmt.Errorf("failed to get available commands: %w", err)
	}

	if len(availableCommands) == 0 {
		fmt.Println("No commands available.")
		return nil
	}

	promptsDir := filepath.Join(projectDir, ".github", "prompts")

	// Check each available command
	fmt.Println("Available commands:")
	installedCount := 0
	for _, command := range availableCommands {
		commandFilePath := filepath.Join(promptsDir, fmt.Sprintf("%s.prompt.md", command))
		if _, err := os.Stat(commandFilePath); err == nil {
			fmt.Printf("  ✅ %s (installed)\n", command)
			installedCount++
		} else {
			fmt.Printf("  ⬜ %s (not installed)\n", command)
		}
	}

	fmt.Printf("\nSummary: %d/%d commands installed\n", installedCount, len(availableCommands))

	if installedCount == 0 {
		fmt.Println("\n💡 Use 'anyagent add command <command-name>' to install commands")
	}

	return nil
}

// removeCommandFile removes a command file
func removeCommandFile(filePath, agentType string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would remove %s command file: %s\n", agentType, filePath)
		return nil
	}

	fmt.Printf("🗑️  Removing %s command file: %s\n", agentType, filePath)
	return os.Remove(filePath)
}

// removeInstalledCommandFromConfig removes a command entry from .anyagent.yaml
func removeInstalledCommandFromConfig(projectDir, command string, dryRun bool) error {
	configPath := config.GetProjectConfigPath(projectDir)
	projectConfig, err := config.LoadProjectConfig(configPath)
	if err != nil {
		return err
	}

	// Filter out the command
	var newCommands []string
	removed := false
	for _, c := range projectConfig.InstalledCommands {
		if c != command {
			newCommands = append(newCommands, c)
		} else {
			removed = true
		}
	}
	if !removed {
		return nil
	}
	projectConfig.InstalledCommands = newCommands
	if dryRun {
		fmt.Printf("[DRY RUN] Would remove command '%s' from .anyagent.yaml\n", command)
		return nil
	}
	return projectConfig.Save(configPath)
}
