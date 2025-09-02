package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/anyagent/internal/config"
)

// RunAddCommand executes the add command functionality
func RunAddCommand(command, projectDir string, dryRun bool) error {
	fmt.Printf("Adding %s command to project...\n", command)

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
		return fmt.Errorf("project is not initialized with anyagent. Run 'anyagent init' first")
	}

	// Validate command name
	if err := validateCommand(command); err != nil {
		return err
	}

	// Get the command template content
	commandContent, err := getCommandTemplate(command)
	if err != nil {
		return fmt.Errorf("failed to get command template: %w", err)
	}

	// Create GitHub Copilot prompts directory
	copilotPromptsDir := filepath.Join(projectDir, ".github", "prompts")
	if err := createPromptsDirectory(copilotPromptsDir, dryRun); err != nil {
		return fmt.Errorf("failed to create prompts directory: %w", err)
	}

	// Create the Copilot command file
	commandFilePath := filepath.Join(copilotPromptsDir, fmt.Sprintf("%s.prompt.md", command))
	if err := createCommandFile(commandFilePath, commandContent, dryRun); err != nil {
		return fmt.Errorf("failed to create command file: %w", err)
	}

	// Create Amazon Q Developer prompts directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not get home directory for Amazon Q Developer prompts: %v\n", err)
	} else {
		qdevPromptsDir := filepath.Join(homeDir, ".aws", "amazonq", "prompts")
		if err := createPromptsDirectory(qdevPromptsDir, dryRun); err != nil {
			fmt.Printf("⚠️  Warning: Could not create Amazon Q Developer prompts directory: %v\n", err)
		} else {
			// Convert command name: replace hyphens and underscores with spaces
			qdevCommandName := strings.ReplaceAll(strings.ReplaceAll(command, "-", " "), "_", " ")
			qdevCommandFilePath := filepath.Join(qdevPromptsDir, fmt.Sprintf("%s.md", qdevCommandName))

			// Create content without YAML frontmatter for Amazon Q Developer
			qdevContent := removeYAMLFrontmatter(commandContent)

			if err := createCommandFile(qdevCommandFilePath, qdevContent, dryRun); err != nil {
				fmt.Printf("⚠️  Warning: Could not create Amazon Q Developer command file: %v\n", err)
			} else {
				fmt.Printf("📄 Amazon Q Developer prompt created: ~/.aws/amazonq/prompts/%s.md\n", qdevCommandName)
			}
		}
	}

	fmt.Printf("✅ %s command added successfully\n", command)
	fmt.Printf("💡 Use '/prompt %s' in VS Code Copilot Chat to activate this command\n", command)
	fmt.Printf("💡 Use '@%s' in Amazon Q Developer Chat to activate this command\n", strings.ReplaceAll(strings.ReplaceAll(command, "-", " "), "_", " "))
	return nil
}

// validateCommand validates the command name and checks if it's available
func validateCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	// Check if command contains invalid characters
	if strings.ContainsAny(command, "/\\<>:\"|?*") {
		return fmt.Errorf("command name contains invalid characters: %s", command)
	}

	// Get available commands
	availableCommands, err := config.GetAvailableCommands()
	if err != nil {
		return fmt.Errorf("failed to get available commands: %w", err)
	}

	// Check if command is available
	for _, available := range availableCommands {
		if command == available {
			return nil
		}
	}

	return fmt.Errorf("command '%s' is not available. Available commands: %s", command, strings.Join(availableCommands, ", "))
}

// getCommandTemplate retrieves the template content for the specified command
func getCommandTemplate(command string) (string, error) {
	return config.GetCommandTemplate(command)
}

// createPromptsDirectory creates the .github/prompts directory
func createPromptsDirectory(dir string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would create directory: %s\n", dir)
		return nil
	}

	fmt.Printf("📁 Creating prompts directory: %s\n", dir)
	return os.MkdirAll(dir, 0755)
}

// createCommandFile creates the command prompt file
func createCommandFile(filePath, content string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY RUN] Would create command file: %s\n", filePath)
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

	fmt.Printf("📄 Creating command file: %s\n", filePath)
	return os.WriteFile(filePath, []byte(content), 0644)
}

// ListAvailableCommands displays all available commands
func ListAvailableCommands() error {
	commands, err := config.GetAvailableCommands()
	if err != nil {
		return fmt.Errorf("failed to get available commands: %w", err)
	}

	if len(commands) == 0 {
		fmt.Println("No commands available.")
		return nil
	}

	fmt.Println("Available commands:")
	for _, command := range commands {
		fmt.Printf("  • %s\n", command)
	}
	fmt.Printf("\nUsage: anyagent add command <command-name>\n")
	fmt.Printf("After adding, use '/prompt <command-name>' in VS Code Copilot Chat\n")

	return nil
}

// removeYAMLFrontmatter removes YAML frontmatter from the content
func removeYAMLFrontmatter(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		// No YAML frontmatter found
		return content
	}

	// Find the end of frontmatter
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			// Found end of frontmatter, return content after it
			if i+1 < len(lines) {
				return strings.Join(lines[i+1:], "\n")
			}
			return ""
		}
	}

	// No end of frontmatter found, return original content
	return content
}
