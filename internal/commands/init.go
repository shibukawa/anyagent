package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/anyagent/internal/config"
)

// AIAgent represents a supported AI agent
type AIAgent struct {
	Name         string
	DisplayName  string
	ConfigPath   string
	NeedsSymlink bool
}

// SupportedAgents defines the list of supported AI agents
var SupportedAgents = []AIAgent{
	{
		Name:         "copilot",
		DisplayName:  "GitHub Copilot",
		ConfigPath:   ".github/copilot-instructions.md",
		NeedsSymlink: true,
	},
	{
		Name:         "qdev",
		DisplayName:  "Amazon Q Developer",
		ConfigPath:   ".amazonq/config.json",
		NeedsSymlink: true,
	},
	{
		Name:         "claude",
		DisplayName:  "Claude Code",
		ConfigPath:   ".claude/config.json",
		NeedsSymlink: true,
	},
	{
		Name:         "junie",
		DisplayName:  "IntelliJ IDEA Junie",
		ConfigPath:   ".junie/settings.json",
		NeedsSymlink: true,
	},
	{
		Name:         "gemini",
		DisplayName:  "Gemini Code",
		ConfigPath:   ".gemini/config.json",
		NeedsSymlink: true,
	},
}

// InitParams holds the parameters for project initialization
type InitParams struct {
	ProjectName        string
	ProjectDescription string
	DynamicParameters  map[string]string // For template parameters
	SelectedAgents     []AIAgent
	ProjectDir         string
}

// RunInit executes the init command functionality
func RunInit(projectDir string, agentNames []string, dryRun bool) error {
	return RunInitWithParams(projectDir, agentNames, "", "", dryRun)
}

// RunInitWithParams executes the init command with predefined parameters (for testing)
func RunInitWithParams(projectDir string, agentNames []string, projectName, projectDesc string, dryRun bool) error {
	fmt.Printf("Initializing anyagent configuration for project...\n")

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

	// Initialize parameters
	params := &InitParams{
		ProjectDir: projectDir,
	}

	// Get AI agents selection
	if len(agentNames) > 0 {
		// Use provided agent names
		selectedAgents, err := validateAgentNames(agentNames)
		if err != nil {
			return fmt.Errorf("invalid agent names: %w", err)
		}
		params.SelectedAgents = selectedAgents
	} else {
		// Run wizard to select agents
		selectedAgents, err := selectAgentsWizard()
		if err != nil {
			return fmt.Errorf("agent selection failed: %w", err)
		}
		params.SelectedAgents = selectedAgents
	}

	// Get project parameters (including dynamic template parameters)
	if projectName != "" && projectDesc != "" {
		// Use provided parameters (for testing)
		params.ProjectName = projectName
		params.ProjectDescription = projectDesc
		params.DynamicParameters = map[string]string{
			"PROJECT_NAME":        projectName,
			"PROJECT_DESCRIPTION": projectDesc,
		}
	} else {
		if err := getProjectParametersWithTemplate(params); err != nil {
			return fmt.Errorf("failed to get project parameters: %w", err)
		}
	}

	// Create AGENTS.md file
	if err := createAgentsFile(params, dryRun); err != nil {
		return fmt.Errorf("failed to create AGENTS.md: %w", err)
	}

	// Create symlinks for selected agents
	if err := createAgentSymlinks(params, dryRun); err != nil {
		return fmt.Errorf("failed to create agent symlinks: %w", err)
	}

	// Save project configuration
	if err := saveProjectConfig(params, dryRun); err != nil {
		return fmt.Errorf("failed to save project configuration: %w", err)
	}

	fmt.Printf("✅ Project initialization completed successfully\n")
	return nil
}

// validateAgentNames validates the provided agent names
func validateAgentNames(agentNames []string) ([]AIAgent, error) {
	var selectedAgents []AIAgent

	for _, name := range agentNames {
		found := false
		for _, agent := range SupportedAgents {
			if agent.Name == name {
				selectedAgents = append(selectedAgents, agent)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("unsupported agent: %s", name)
		}
	}

	return selectedAgents, nil
}

// selectAgentsWizard runs an interactive wizard to select AI agents
func selectAgentsWizard() ([]AIAgent, error) {
	fmt.Printf("\nSelect AI agents to configure (enter numbers separated by spaces, e.g., '1 3 5'):\n")

	for i, agent := range SupportedAgents {
		fmt.Printf("  %d. %s\n", i+1, agent.DisplayName)
	}

	fmt.Printf("Enter your selection: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("no agents selected")
	}

	// Parse selections
	selections := strings.Fields(input)
	var selectedAgents []AIAgent

	for _, sel := range selections {
		var index int
		if _, err := fmt.Sscanf(sel, "%d", &index); err != nil {
			return nil, fmt.Errorf("invalid selection: %s", sel)
		}

		if index < 1 || index > len(SupportedAgents) {
			return nil, fmt.Errorf("selection out of range: %d", index)
		}

		selectedAgents = append(selectedAgents, SupportedAgents[index-1])
	}

	if len(selectedAgents) == 0 {
		return nil, fmt.Errorf("no valid agents selected")
	}

	return selectedAgents, nil
}

// getProjectParameters prompts for project parameters
func getProjectParameters(params *InitParams) error {
	reader := bufio.NewReader(os.Stdin)

	// Get project name
	fmt.Printf("\nEnter project name: ")
	projectName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read project name: %w", err)
	}
	params.ProjectName = strings.TrimSpace(projectName)

	if params.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	// Get project description
	fmt.Printf("Enter project description: ")
	projectDesc, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read project description: %w", err)
	}
	params.ProjectDescription = strings.TrimSpace(projectDesc)

	if params.ProjectDescription == "" {
		return fmt.Errorf("project description is required")
	}

	return nil
}

// getProjectParametersWithTemplate gets project parameters and scans template for dynamic parameters
func getProjectParametersWithTemplate(params *InitParams) error {
	reader := bufio.NewReader(os.Stdin)

	// Get basic project parameters first
	// Get project name
	fmt.Printf("\nEnter project name: ")
	projectName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read project name: %w", err)
	}
	params.ProjectName = strings.TrimSpace(projectName)

	if params.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	// Get project description
	fmt.Printf("Enter project description: ")
	projectDesc, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read project description: %w", err)
	}
	params.ProjectDescription = strings.TrimSpace(projectDesc)

	if params.ProjectDescription == "" {
		return fmt.Errorf("project description is required")
	}

	// Extract dynamic parameters from template
	template := config.GetAGENTSTemplate()
	templateParams := config.ExtractTemplateParameters(template)

	// Initialize dynamic parameters map
	params.DynamicParameters = make(map[string]string)

	// Always include basic parameters
	params.DynamicParameters["PROJECT_NAME"] = params.ProjectName
	params.DynamicParameters["PROJECT_DESCRIPTION"] = params.ProjectDescription

	// Get additional template parameters
	if len(templateParams) > 2 { // More than just PROJECT_NAME and PROJECT_DESCRIPTION
		fmt.Printf("\nThe template requires additional parameters:\n")

		for _, paramName := range templateParams {
			// Skip basic parameters we already have
			if paramName == "PROJECT_NAME" || paramName == "PROJECT_DESCRIPTION" {
				continue
			}

			fmt.Printf("Enter %s: ", paramName)
			paramValue, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read parameter %s: %w", paramName, err)
			}

			paramValue = strings.TrimSpace(paramValue)
			if paramValue == "" {
				fmt.Printf("Warning: %s is empty, will leave placeholder in template\n", paramName)
			}

			params.DynamicParameters[paramName] = paramValue
		}
	}

	return nil
}

// createAgentsFile creates the AGENTS.md file with populated parameters
func createAgentsFile(params *InitParams, dryRun bool) error {
	// Get the template content
	template := config.GetAGENTSTemplate()

	// Replace placeholders using dynamic parameters
	content := config.ReplaceTemplateParameters(template, params.DynamicParameters)

	agentsPath := filepath.Join(params.ProjectDir, "AGENTS.md")

	if dryRun {
		fmt.Printf("[DRY RUN] Would create AGENTS.md at: %s\n", agentsPath)
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

	fmt.Printf("📄 Creating AGENTS.md...\n")
	if err := os.WriteFile(agentsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}

	return nil
}

// createAgentSymlinks creates symlinks for selected agents
func createAgentSymlinks(params *InitParams, dryRun bool) error {
	agentsPath := filepath.Join(params.ProjectDir, "AGENTS.md")

	for _, agent := range params.SelectedAgents {
		if !agent.NeedsSymlink {
			continue
		}

		// Create the directory for the agent config if it doesn't exist
		configDir := filepath.Dir(filepath.Join(params.ProjectDir, agent.ConfigPath))
		if err := os.MkdirAll(configDir, 0755); err != nil && !dryRun {
			return fmt.Errorf("failed to create directory %s: %w", configDir, err)
		}

		symlinkPath := filepath.Join(params.ProjectDir, agent.ConfigPath)

		if dryRun {
			fmt.Printf("[DRY RUN] Would create symlink: %s -> AGENTS.md\n", symlinkPath)
			continue
		}

		// Remove existing file/symlink if it exists
		if _, err := os.Lstat(symlinkPath); err == nil {
			if err := os.Remove(symlinkPath); err != nil {
				return fmt.Errorf("failed to remove existing file %s: %w", symlinkPath, err)
			}
		}

		// Create relative symlink
		relPath, err := filepath.Rel(configDir, agentsPath)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		fmt.Printf("🔗 Creating symlink for %s: %s -> %s\n", agent.DisplayName, agent.ConfigPath, relPath)
		if err := os.Symlink(relPath, symlinkPath); err != nil {
			return fmt.Errorf("failed to create symlink %s: %w", symlinkPath, err)
		}
	}

	return nil
}

// saveProjectConfig saves the project configuration to .anyagent.yaml
func saveProjectConfig(params *InitParams, dryRun bool) error {
	// Prepare agent names
	var agentNames []string
	for _, agent := range params.SelectedAgents {
		agentNames = append(agentNames, agent.Name)
	}

	// Create project configuration
	projectConfig := &config.ProjectConfig{
		ProjectName:        params.ProjectName,
		ProjectDescription: params.ProjectDescription,
		InstalledRules:     []string{},
	}

	if dryRun {
		fmt.Printf("[DRY RUN] Would save project configuration to: %s\n", config.GetProjectConfigPath(params.ProjectDir))
		fmt.Printf("[DRY RUN] Configuration content:\n")
		fmt.Printf("  Project: %s\n", projectConfig.ProjectName)
		fmt.Printf("  Description: %s\n", projectConfig.ProjectDescription)
		fmt.Printf("  Installed Rules: %v\n", projectConfig.InstalledRules)
		return nil
	}

	fmt.Printf("💾 Saving project configuration...\n")
	if err := config.SaveProjectConfig(params.ProjectDir, projectConfig); err != nil {
		return fmt.Errorf("failed to save project configuration: %w", err)
	}

	fmt.Printf("✅ Project configuration saved to .anyagent.yaml\n")
	return nil
}
