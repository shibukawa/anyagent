package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/shibukawa/anyagent/internal/config"
)

// RunEditTemplate executes the edit-template command functionality
func RunEditTemplate(configDir string, dryRun bool, hardReset bool) error {
	if hardReset {
		fmt.Printf("Hard reset mode: Resetting all templates to original versions...\n")
	} else {
		fmt.Printf("Setting up anyagent template editing environment...\n")
	}

	// Hard reset mode: force recreate everything
	if hardReset {
		fmt.Printf("Performing hard reset of template environment...\n")
		if err := performHardReset(configDir); err != nil {
			return fmt.Errorf("failed to perform hard reset: %w", err)
		}
		fmt.Printf("✅ Template environment reset to original state\n")
	} else {
		// Check if config directory exists
		if !config.CheckUserConfigExists(configDir) {
			fmt.Printf("Creating new template environment at: %s\n", configDir)
			// Create the configuration directory and all necessary components
			if err := setupNewTemplateEnvironment(configDir); err != nil {
				return fmt.Errorf("failed to setup template environment: %w", err)
			}
			fmt.Printf("✅ Template environment created successfully\n")
		} else {
			fmt.Printf("Found existing template environment at: %s\n", configDir)
			// Validate existing environment and update if necessary
			if !ValidateTemplateEnvironment(configDir) {
				fmt.Printf("Updating incomplete template environment...\n")
				if err := updateTemplateEnvironment(configDir); err != nil {
					return fmt.Errorf("failed to update template environment: %w", err)
				}
				fmt.Printf("✅ Template environment updated successfully\n")
			} else {
				fmt.Printf("✅ Template environment is up to date\n")
			}
		}
	}

	// Launch VSCode if not in dry run mode
	if !dryRun {
		if err := LaunchVSCode(configDir, false); err != nil {
			return fmt.Errorf("failed to launch VSCode: %w", err)
		}
	} else {
		if err := LaunchVSCode(configDir, true); err != nil {
			return fmt.Errorf("failed to launch VSCode: %w", err)
		}
	}

	return nil
}

// ValidateTemplateEnvironment checks if the template environment is complete and valid
func ValidateTemplateEnvironment(configDir string) bool {
	// Check if config directory exists
	if !config.CheckUserConfigExists(configDir) {
		return false
	}

	// Check for required template structure
	requiredPaths := []string{
		"templates",
		"templates/commands",
		"templates/extra_rules",
		"templates/AGENTS.md.tmpl",
		"templates/mcp.yaml",
		"templates/commands/general.md",
		"templates/commands/coding.md",
		"templates/commands/project-specific.md",
		"templates/extra_rules/go.md",
		"templates/extra_rules/ts.md",
		"templates/extra_rules/docker.md",
		"templates/extra_rules/python.md",
		"templates/extra_rules/react.md",
		"README.md",
		"AGENTS.md",
		".amazonq",
		".claude",
		".junie",
		".gemini",
	}

	for _, path := range requiredPaths {
		fullPath := filepath.Join(configDir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return false
		}
	}

	// Check for required symbolic links
	requiredSymlinks := map[string]string{
		".amazonq/rules/AGENTS.md": "../../AGENTS.md",
		".claude/AGENTS.md":        "../AGENTS.md",
		".junie/AGENTS.md":         "../AGENTS.md",
		".gemini/AGENTS.md":        "../AGENTS.md",
		"CLAUDE.md":                "AGENTS.md", // Project root CLAUDE.md
	}

	for symlinkPath, expectedTarget := range requiredSymlinks {
		fullPath := filepath.Join(configDir, symlinkPath)
		if target, err := os.Readlink(fullPath); err != nil || target != expectedTarget {
			return false
		}
	}

	return true
}

// LaunchVSCode launches Visual Studio Code with the specified directory
func LaunchVSCode(configDir string, dryRun bool) error {
	// README file to open actively
	readmeFile := filepath.Join(configDir, "README.md")

	if dryRun {
		fmt.Printf("[DRY RUN] Would launch VSCode with directory: %s and open README.md\n", configDir)
		return nil
	}

	// Try different VSCode executable names based on platform
	var vscodeCommands []string

	switch runtime.GOOS {
	case "darwin": // macOS
		vscodeCommands = []string{
			"code",
			"code-insiders",
			"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
			"/Applications/Visual Studio Code - Insiders.app/Contents/Resources/app/bin/code",
		}
	case "windows":
		vscodeCommands = []string{
			"code.cmd",
			"code",
			"code-insiders.cmd",
			"code-insiders",
		}
	default: // Linux and others
		vscodeCommands = []string{
			"code",
			"code-insiders",
			"/usr/bin/code",
			"/snap/bin/code",
		}
	}

	var cmd *exec.Cmd
	var cmdName string

	for _, cmdName = range vscodeCommands {
		if _, err := exec.LookPath(cmdName); err == nil {
			// Open folder and README file - folder first, then the file to make it active
			cmd = exec.Command(cmdName, configDir, readmeFile)
			break
		}
	}

	if cmd == nil {
		return fmt.Errorf("VSCode executable not found. Please ensure VSCode is installed and available in PATH.\nTried: %v", vscodeCommands)
	}

	fmt.Printf("Opening with VSCode...\n")

	// Start VSCode in the background
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start VSCode: %w", err)
	}

	fmt.Printf("✅ VSCode launched successfully\n")
	return nil
}

// setupNewTemplateEnvironment creates a complete new template environment
func setupNewTemplateEnvironment(configDir string) error {
	fmt.Printf("📁 Creating configuration directory...\n")
	// Create user config directory
	if err := config.CreateUserConfigDir(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	fmt.Printf("📂 Creating template structure...\n")
	// Create template structure
	if err := config.CreateTemplateStructure(configDir); err != nil {
		return fmt.Errorf("failed to create template structure: %w", err)
	}

	fmt.Printf("📄 Creating template files...\n")
	// Create template files
	if err := config.CreateTemplateFiles(configDir); err != nil {
		return fmt.Errorf("failed to create template files: %w", err)
	}

	fmt.Printf("⚙️  Creating anyagent project configuration...\n")
	// Create anyagent project configuration
	if err := config.CreateAnyagentProject(configDir); err != nil {
		return fmt.Errorf("failed to create anyagent project: %w", err)
	}

	return nil
}

// updateTemplateEnvironment updates an existing template environment
func updateTemplateEnvironment(configDir string) error {
	// Ensure template structure exists
	if err := config.CreateTemplateStructure(configDir); err != nil {
		return fmt.Errorf("failed to update template structure: %w", err)
	}

	// Add only missing template files (do not overwrite existing files unless --force)
	if err := config.CreateTemplateFilesIfMissing(configDir); err != nil {
		return fmt.Errorf("failed to add missing template files: %w", err)
	}

	// Ensure anyagent project configuration exists
	agentsFile := filepath.Join(configDir, "AGENTS.md")
	if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
		if err := config.CreateAnyagentProject(configDir); err != nil {
			return fmt.Errorf("failed to create anyagent project: %w", err)
		}
	}

	return nil
}

// printTemplateInfo prints helpful information about the template environment

// performHardReset performs a complete reset of the template environment
func performHardReset(configDir string) error {
	// Remove existing directory if it exists
	if config.CheckUserConfigExists(configDir) {
		fmt.Printf("🗑️  Removing existing template environment...\n")
		if err := os.RemoveAll(configDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	// Create fresh template environment
	fmt.Printf("🔄 Creating fresh template environment...\n")
	return setupNewTemplateEnvironment(configDir)
}
