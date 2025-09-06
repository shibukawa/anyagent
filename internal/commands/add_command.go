package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/anyagent/internal/config"
)

// RunAddCommand executes the add command functionality
func RunAddCommand(command, projectDir string, dryRun bool, global bool) error {
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
		return fmt.Errorf("project is not initialized with anyagent. Run 'anyagent sync' first")
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

	// Create Copilot prompt if Copilot is selected (or no config present)
	if shouldCreateCopilotCommandFiles(projectDir) {
		copilotPromptsDir := filepath.Join(projectDir, ".github", "prompts")
		if err := createPromptsDirectory(copilotPromptsDir, dryRun); err != nil {
			return fmt.Errorf("failed to create prompts directory: %w", err)
		}

		commandFilePath := filepath.Join(copilotPromptsDir, fmt.Sprintf("%s.prompt.md", command))
		if err := createCommandFile(commandFilePath, commandContent, dryRun); err != nil {
			return fmt.Errorf("failed to create command file: %w", err)
		}
	}

	// Create Amazon Q Developer prompt if Q Dev is selected (only with --global)
	if shouldCreateQDevCommandFiles(projectDir) {
		if global {
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
		} else {
			fmt.Printf("ℹ️  Q Dev selected: use '--global' to install to ~/.aws/amazonq/prompts\n")
		}
	}

	// Create Codex prompt if Codex is selected
	if shouldCreateCodexCommandFiles(projectDir) {
		if global {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("⚠️  Warning: Could not get home directory for Codex prompts: %v\n", err)
			} else {
				codexPromptsDir := filepath.Join(homeDir, ".codex", "prompts")
				if err := createPromptsDirectory(codexPromptsDir, dryRun); err != nil {
					fmt.Printf("⚠️  Warning: Could not create Codex prompts directory: %v\n", err)
				} else {
					codexCommandFilePath := filepath.Join(codexPromptsDir, fmt.Sprintf("%s.md", command))
					codexContent := removeYAMLFrontmatter(commandContent)
					if err := createCommandFile(codexCommandFilePath, codexContent, dryRun); err != nil {
						fmt.Printf("⚠️  Warning: Could not create Codex command file: %v\n", err)
					} else {
						fmt.Printf("📄 Codex prompt created: ~/.codex/prompts/%s.md\n", command)
					}
				}
			}
		} else {
			fmt.Printf("ℹ️  Codex selected: use '--global' to install to ~/.codex/prompts\n")
		}
	}

	// Create Claude Code prompt if Claude is selected
	if shouldCreateClaudeCommandFiles(projectDir) {
		claudeDir := filepath.Join(projectDir, ".claude", "commands")
		if err := createPromptsDirectory(claudeDir, dryRun); err != nil {
			fmt.Printf("⚠️  Warning: Could not create Claude commands directory: %v\n", err)
		} else {
			claudeCommandFilePath := filepath.Join(claudeDir, fmt.Sprintf("%s.md", command))
			// Build Claude-specific content: add YAML frontmatter with allowed-tools and description
			claudeContent := buildClaudeCommandContent(commandContent)
			if err := createCommandFile(claudeCommandFilePath, claudeContent, dryRun); err != nil {
				fmt.Printf("⚠️  Warning: Could not create Claude command file: %v\n", err)
			} else {
				fmt.Printf("📄 Claude command created: .claude/commands/%s.md\n", command)
			}
		}
	}

	// Create Gemini Code command if Gemini is selected
	if shouldCreateGeminiCommandFiles(projectDir) {
		geminiDir := filepath.Join(projectDir, ".gemini", "commands")
		if err := createPromptsDirectory(geminiDir, dryRun); err != nil {
			fmt.Printf("⚠️  Warning: Could not create Gemini commands directory: %v\n", err)
		} else {
			geminiCommandFilePath := filepath.Join(geminiDir, fmt.Sprintf("%s.toml", command))
			tomlContent := buildGeminiCommandTOML(commandContent)
			if err := createCommandFile(geminiCommandFilePath, tomlContent, dryRun); err != nil {
				fmt.Printf("⚠️  Warning: Could not create Gemini command file: %v\n", err)
			} else {
				fmt.Printf("📄 Gemini command created: .gemini/commands/%s.toml\n", command)
			}
		}
	}

	fmt.Printf("✅ %s command added successfully\n", command)
	if shouldCreateCopilotCommandFiles(projectDir) {
		fmt.Printf("💡 Use '/prompt %s' in VS Code Copilot Chat to activate this command\n", command)
	}
	if shouldCreateQDevCommandFiles(projectDir) && global {
		fmt.Printf("💡 Use '@%s' in Amazon Q Developer Chat to activate this command\n", strings.ReplaceAll(strings.ReplaceAll(command, "-", " "), "_", " "))
	}
	if shouldCreateCodexCommandFiles(projectDir) && global {
		fmt.Printf("💡 Use '/%s' in Codex to activate this command\n", command)
	}
	if shouldCreateClaudeCommandFiles(projectDir) {
		fmt.Printf("💡 Claude Code: use the command from .claude/commands/%s.md\n", command)
	}
	if shouldCreateGeminiCommandFiles(projectDir) {
		fmt.Printf("💡 Gemini Code: command saved at .gemini/commands/%s.toml\n", command)
	}

	// Track installed command in project config for future syncs (info only)
	if err := addInstalledCommandToConfig(projectDir, command, dryRun); err != nil {
		fmt.Printf("⚠️  Warning: Failed to update project config with command '%s': %v\n", command, err)
	}
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

// addInstalledCommandToConfig records the installed command into .anyagent.yaml
func addInstalledCommandToConfig(projectDir, command string, dryRun bool) error {
	configPath := config.GetProjectConfigPath(projectDir)
	projectConfig, err := config.LoadProjectConfig(configPath)
	if err != nil {
		return err
	}
	// If not present, add it
	present := false
	for _, c := range projectConfig.InstalledCommands {
		if c == command {
			present = true
			break
		}
	}
	if !present {
		projectConfig.InstalledCommands = append(projectConfig.InstalledCommands, command)
		if dryRun {
			fmt.Printf("[DRY RUN] Would record installed command '%s' into .anyagent.yaml\n", command)
			return nil
		}
		if err := projectConfig.Save(configPath); err != nil {
			return fmt.Errorf("failed to save project config: %w", err)
		}
		fmt.Printf("💾 Project config updated: recorded command '%s'\n", command)
	}
	return nil
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

// buildClaudeCommandContent wraps a command template body with Claude-specific YAML frontmatter
// including allowed-tools and description. It extracts description from the original template
// frontmatter if available.
func buildClaudeCommandContent(templateContent string) string {
	desc := extractDescriptionFromTemplate(templateContent)
	body := removeYAMLFrontmatter(templateContent)
	// Default allowed-tools empty; specialized commands may add more later.
	header := []string{
		"---",
		"allowed-tools: []",
		fmt.Sprintf("description: %s", quoteIfNeeded(desc)),
		"---",
		"",
	}
	return strings.Join(header, "\n") + body
}

// extractDescriptionFromTemplate tries to read 'description:' from a YAML frontmatter at the top.
func extractDescriptionFromTemplate(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return "" // no frontmatter
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			break
		}
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "description:") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "description:"))
			// Strip surrounding quotes if present
			v = strings.Trim(v, "'\"")
			return v
		}
	}
	return ""
}

func quoteIfNeeded(s string) string {
	if s == "" {
		return "''"
	}
	if strings.ContainsAny(s, ":#[]{}\"'\n") {
		// Use single quotes and escape existing single quotes by doubling them
		return "'" + strings.ReplaceAll(s, "'", "''") + "'"
	}
	return s
}

// buildGeminiCommandTOML constructs a TOML definition with description and prompt
func buildGeminiCommandTOML(templateContent string) string {
	desc := extractDescriptionFromTemplate(templateContent)
	prompt := removeYAMLFrontmatter(templateContent)
	// Escape description for TOML basic string
	descEsc := strings.ReplaceAll(desc, "\\", "\\\\")
	descEsc = strings.ReplaceAll(descEsc, "\"", "\\\"")
	if descEsc == "" {
		descEsc = ""
	}
	// Use TOML multiline basic string for prompt
	// Trim leading newline to keep formatting tidy
	if strings.HasPrefix(prompt, "\n") {
		prompt = prompt[1:]
	}
	return fmt.Sprintf("description = \"%s\"\nprompt = \"\"\"%s\"\"\"\n", descEsc, prompt)
}
