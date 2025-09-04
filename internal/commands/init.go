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
        ConfigPath:   ".amazonq/rules/AGENTS.md",
        NeedsSymlink: true,
    },
	{
		Name:         "claude",
		DisplayName:  "Claude Code",
		ConfigPath:   "CLAUDE.md",
		NeedsSymlink: true,
	},
    {
        Name:         "gemini",
        DisplayName:  "Gemini Code",
        ConfigPath:   "",
        NeedsSymlink: false,
    },
    {
        Name:         "codex",
        DisplayName:  "ChatGPT Codex",
        ConfigPath:   "",
        NeedsSymlink: false,
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

// RunSync executes the sync command functionality
// - If project is already initialized, it updates AGENTS.md from the latest template using stored parameters
// - If --agents is specified, it reconfigures enabled agents (removing deselected agent artifacts)
// - If project is not initialized, it behaves like RunInit
func RunSync(projectDir string, agentNames []string, dryRun bool) error {
    fmt.Printf("Synchronizing anyagent configuration for project...\n")

    // Determine project directory
    if projectDir == "" {
        var err error
        projectDir, err = os.Getwd()
        if err != nil {
            return fmt.Errorf("failed to get current directory: %w", err)
        }
    }
    if _, err := os.Stat(projectDir); os.IsNotExist(err) {
        return fmt.Errorf("project directory does not exist: %s", projectDir)
    }

    configPath := config.GetProjectConfigPath(projectDir)
    projectConfig, err := config.LoadProjectConfig(configPath)
    if err != nil {
        return fmt.Errorf("failed to load project config: %w", err)
    }

    // If AGENTS.md doesn't exist, treat as first-time initialization
    agentsPath := filepath.Join(projectDir, "AGENTS.md")
    if _, err := os.Stat(agentsPath); os.IsNotExist(err) && len(projectConfig.Parameters) == 0 && projectConfig.ProjectName == "" {
        // Fallback to init flow
        return RunInit(projectDir, agentNames, dryRun)
    }

    // Determine new selection of agents
    var selectedAgents []AIAgent
    if len(agentNames) > 0 {
        selectedAgents, err = validateAgentNames(agentNames)
        if err != nil {
            return fmt.Errorf("invalid agent names: %w", err)
        }
    } else if len(projectConfig.EnabledAgents) > 0 {
        selectedAgents = agentsFromNames(projectConfig.EnabledAgents)
    } else {
        selectedAgents, err = selectAgentsWizard()
        if err != nil {
            return fmt.Errorf("agent selection failed: %w", err)
        }
    }

    // Compute agents to remove (previous - new)
    prevAgents := map[string]bool{}
    for _, a := range projectConfig.EnabledAgents {
        prevAgents[a] = true
    }
    newAgents := map[string]bool{}
    var newAgentNames []string
    for _, a := range selectedAgents {
        newAgents[a.Name] = true
        newAgentNames = append(newAgentNames, a.Name)
    }

    var removedAgents []string
    for a := range prevAgents {
        if !newAgents[a] {
            removedAgents = append(removedAgents, a)
        }
    }

    // List installed commands for Copilot and Q Developer
    installedCopilotCmds, installedQDevCmds, err := listInstalledCommands(projectDir)
    if err != nil {
        return err
    }

    // Rules from project config
    installedRules := projectConfig.InstalledRules

    // Remove artifacts for deselected agents
    for _, agent := range removedAgents {
        if err := removeAgentArtifacts(agent, projectDir, installedCopilotCmds, installedQDevCmds, installedRules, dryRun); err != nil {
            return fmt.Errorf("failed to remove artifacts for %s: %w", agent, err)
        }
    }

    // Regenerate AGENTS.md using stored parameters and rules
    if dryRun {
        fmt.Printf("[DRY RUN] Would regenerate AGENTS.md using stored parameters and rules\n")
    } else {
        if err := projectConfig.RegenerateAgentsFileAt(projectDir); err != nil {
            return fmt.Errorf("failed to regenerate AGENTS.md: %w", err)
        }
        fmt.Printf("📄 AGENTS.md regenerated from latest template\n")
    }

    // Recreate symlinks for selected agents
    if err := createAgentSymlinks(&InitParams{ProjectDir: projectDir, SelectedAgents: selectedAgents}, dryRun); err != nil {
        return fmt.Errorf("failed to create agent symlinks: %w", err)
    }

    // Reinstall commands for the selected agent from project config (info kept even if agent changes)
    if len(selectedAgents) == 1 {
        if err := reinstallCommandsForAgent(selectedAgents[0].Name, projectDir, projectConfig.InstalledCommands, dryRun); err != nil {
            fmt.Printf("⚠️  Warning: Failed to reinstall commands for agent %s: %v\n", selectedAgents[0].Name, err)
        }
    }

    // Update and save project configuration
    if !dryRun {
        projectConfig.EnabledAgents = newAgentNames
        if err := config.SaveProjectConfig(projectDir, projectConfig); err != nil {
            return fmt.Errorf("failed to save project configuration: %w", err)
        }
    } else {
        fmt.Printf("[DRY RUN] Would save enabled agents to .anyagent.yaml: %v\n", newAgentNames)
    }

    fmt.Printf("✅ Project synchronization completed successfully\n")
    return nil
}

// listInstalledCommands returns installed command names for Copilot (project) and Q Developer (home)
func listInstalledCommands(projectDir string) ([]string, []string, error) {
    available, err := config.GetAvailableCommands()
    if err != nil {
        return nil, nil, fmt.Errorf("failed to get available commands: %w", err)
    }

    var copilot []string
    promptsDir := filepath.Join(projectDir, ".github", "prompts")
    for _, c := range available {
        p := filepath.Join(promptsDir, fmt.Sprintf("%s.prompt.md", c))
        if _, err := os.Stat(p); err == nil {
            copilot = append(copilot, c)
        }
    }

    var qdev []string
    homeDir, err := os.UserHomeDir()
    if err == nil {
        qdevDir := filepath.Join(homeDir, ".aws", "amazonq", "prompts")
        for _, c := range available {
            qname := strings.ReplaceAll(strings.ReplaceAll(c, "-", " "), "_", " ")
            p := filepath.Join(qdevDir, fmt.Sprintf("%s.md", qname))
            if _, err := os.Stat(p); err == nil {
                qdev = append(qdev, c)
            }
        }
    }

    return copilot, qdev, nil
}

// removeAgentArtifacts removes symlinks and agent-specific files for a deselected agent
func removeAgentArtifacts(agentName, projectDir string, copilotCmds, qdevCmds, rules []string, dryRun bool) error {
    switch agentName {
    case "copilot":
        // Remove symlink
        symlinkPath := filepath.Join(projectDir, ".github", "copilot-instructions.md")
        if err := removePath(symlinkPath, "GitHub Copilot symlink", dryRun); err != nil {
            return err
        }
        // Remove commands
        for _, cmd := range copilotCmds {
            path := filepath.Join(projectDir, ".github", "prompts", fmt.Sprintf("%s.prompt.md", cmd))
            if err := removePath(path, fmt.Sprintf("Copilot command '%s'", cmd), dryRun); err != nil {
                return err
            }
        }
        // Remove rule instruction files
        for _, r := range rules {
            path := filepath.Join(projectDir, ".github", "instructions", fmt.Sprintf("%s.instructions.md", r))
            if err := removePath(path, fmt.Sprintf("Copilot rule '%s'", r), dryRun); err != nil {
                return err
            }
        }
    case "qdev":
        symlinkPath := filepath.Join(projectDir, ".amazonq", "rules", "AGENTS.md")
        if err := removePath(symlinkPath, "Amazon Q Developer symlink", dryRun); err != nil {
            return err
        }
        homeDir, err := os.UserHomeDir()
        if err == nil {
            for _, cmd := range qdevCmds {
                qname := strings.ReplaceAll(strings.ReplaceAll(cmd, "-", " "), "_", " ")
                path := filepath.Join(homeDir, ".aws", "amazonq", "prompts", fmt.Sprintf("%s.md", qname))
                if err := removePath(path, fmt.Sprintf("Amazon Q command '%s'", cmd), dryRun); err != nil {
                    return err
                }
            }
        }
    case "claude":
        // Remove CLAUDE.md symlink at project root
        symlinkPath := filepath.Join(projectDir, "CLAUDE.md")
        if err := removePath(symlinkPath, "Claude Code symlink", dryRun); err != nil {
            return err
        }
        // Remove Claude command files (project-local)
        configPath := config.GetProjectConfigPath(projectDir)
        projectConfig, err := config.LoadProjectConfig(configPath)
        if err == nil {
            for _, cmd := range projectConfig.InstalledCommands {
                path := filepath.Join(projectDir, ".claude", "commands", fmt.Sprintf("%s.md", cmd))
                if err := removePath(path, fmt.Sprintf("Claude command '%s'", cmd), dryRun); err != nil {
                    return err
                }
            }
        }
    case "gemini":
        symlinkPath := filepath.Join(projectDir, ".gemini", "config.json")
        if err := removePath(symlinkPath, "Gemini Code symlink", dryRun); err != nil {
            return err
        }
    case "codex":
        // Remove Codex command files based on installed commands in project config
        configPath := config.GetProjectConfigPath(projectDir)
        projectConfig, err := config.LoadProjectConfig(configPath)
        if err != nil {
            return nil // silently skip if no config
        }
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil
        }
        for _, cmd := range projectConfig.InstalledCommands {
            path := filepath.Join(homeDir, ".codex", "prompts", fmt.Sprintf("%s.md", cmd))
            if err := removePath(path, fmt.Sprintf("Codex command '%s'", cmd), dryRun); err != nil {
                return err
            }
        }
    }
    return nil
}

func removePath(path, label string, dryRun bool) error {
    if _, err := os.Lstat(path); err != nil {
        // Nothing to remove
        return nil
    }
    if dryRun {
        fmt.Printf("[DRY RUN] Would remove %s: %s\n", label, path)
        return nil
    }
    fmt.Printf("🗑️  Removing %s: %s\n", label, path)
    return os.Remove(path)
}

// agentsFromNames converts agent names to AIAgent definitions
func agentsFromNames(names []string) []AIAgent {
    var result []AIAgent
    for _, n := range names {
        for _, a := range SupportedAgents {
            if a.Name == n {
                result = append(result, a)
                break
            }
        }
    }
    return result
}

// reinstallCommandsForAgent installs command files for the selected agent using the project config list.
// For Codex and other agents that don't support external commands, do nothing.
func reinstallCommandsForAgent(agentName, projectDir string, commands []string, dryRun bool) error {
    if len(commands) == 0 {
        return nil
    }
    switch agentName {
    case "copilot":
        // Ensure prompts dir
        promptsDir := filepath.Join(projectDir, ".github", "prompts")
        if err := createPromptsDirectory(promptsDir, dryRun); err != nil {
            return err
        }
        for _, c := range commands {
            content, err := getCommandTemplate(c)
            if err != nil {
                fmt.Printf("⚠️  Warning: Command template not found for '%s': %v\n", c, err)
                continue
            }
            path := filepath.Join(promptsDir, fmt.Sprintf("%s.prompt.md", c))
            if err := createCommandFile(path, content, dryRun); err != nil {
                fmt.Printf("⚠️  Warning: Could not create Copilot command '%s': %v\n", c, err)
            }
        }
    case "qdev":
        // Create in ~/.aws/amazonq/prompts
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil // silently skip if no home dir
        }
        dir := filepath.Join(homeDir, ".aws", "amazonq", "prompts")
        if err := createPromptsDirectory(dir, dryRun); err != nil {
            return err
        }
        for _, c := range commands {
            content, err := getCommandTemplate(c)
            if err != nil {
                fmt.Printf("⚠️  Warning: Command template not found for '%s': %v\n", c, err)
                continue
            }
            content = removeYAMLFrontmatter(content)
            qname := strings.ReplaceAll(strings.ReplaceAll(c, "-", " "), "_", " ")
            path := filepath.Join(dir, fmt.Sprintf("%s.md", qname))
            if err := createCommandFile(path, content, dryRun); err != nil {
                fmt.Printf("⚠️  Warning: Could not create Amazon Q command '%s': %v\n", c, err)
            }
        }
    case "claude":
        dir := filepath.Join(projectDir, ".claude", "commands")
        if err := createPromptsDirectory(dir, dryRun); err != nil {
            return err
        }
        for _, c := range commands {
            content, err := getCommandTemplate(c)
            if err != nil {
                fmt.Printf("⚠️  Warning: Command template not found for '%s': %v\n", c, err)
                continue
            }
            content = buildClaudeCommandContent(content)
            path := filepath.Join(dir, fmt.Sprintf("%s.md", c))
            if err := createCommandFile(path, content, dryRun); err != nil {
                fmt.Printf("⚠️  Warning: Could not create Claude command '%s': %v\n", c, err)
            }
        }
    default:
        // gemini: no external command files supported
        return nil
    case "codex":
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil
        }
        dir := filepath.Join(homeDir, ".codex", "prompts")
        if err := createPromptsDirectory(dir, dryRun); err != nil {
            return err
        }
        for _, c := range commands {
            content, err := getCommandTemplate(c)
            if err != nil {
                fmt.Printf("⚠️  Warning: Command template not found for '%s': %v\n", c, err)
                continue
            }
            content = removeYAMLFrontmatter(content)
            path := filepath.Join(dir, fmt.Sprintf("%s.md", c))
            if err := createCommandFile(path, content, dryRun); err != nil {
                fmt.Printf("⚠️  Warning: Could not create Codex command '%s': %v\n", c, err)
            }
        }
    }
    return nil
}

// validateAgentNames validates the provided agent names
func validateAgentNames(agentNames []string) ([]AIAgent, error) {
    if len(agentNames) > 1 {
        return nil, fmt.Errorf("only one agent can be selected at a time")
    }

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
    fmt.Printf("\nSelect one AI agent to configure (enter a number):\n")
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
        return nil, fmt.Errorf("no agent selected")
    }
    var index int
    if _, err := fmt.Sscanf(input, "%d", &index); err != nil {
        return nil, fmt.Errorf("invalid selection: %s", input)
    }
    if index < 1 || index > len(SupportedAgents) {
        return nil, fmt.Errorf("selection out of range: %d", index)
    }
    return []AIAgent{SupportedAgents[index-1]}, nil
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

    // Also ensure MCP symlink for enabled agents to project mcp.yaml (if any)
    // This keeps MCP wiring consistent when switching agents.
    cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(params.ProjectDir))
    if err == nil && len(cfg.MCPServers) > 0 {
        for _, agent := range params.SelectedAgents {
            if err := ensureMCPSymlinkForAgent(agent.Name, params.ProjectDir, dryRun); err != nil {
                return err
            }
        }
    }

    return nil
}

// saveProjectConfig saves the project configuration to .anyagent.yaml
func saveProjectConfig(params *InitParams, dryRun bool) error {
    // Prepare agent names
    var agentNames []string
    for _, a := range params.SelectedAgents {
        agentNames = append(agentNames, a.Name)
    }

    // Create project configuration
    projectConfig := &config.ProjectConfig{
        ProjectName:        params.ProjectName,
        ProjectDescription: params.ProjectDescription,
        InstalledRules:     []string{},
        EnabledAgents:      agentNames,
        Parameters:         params.DynamicParameters,
    }

	if dryRun {
		fmt.Printf("[DRY RUN] Would save project configuration to: %s\n", config.GetProjectConfigPath(params.ProjectDir))
		fmt.Printf("[DRY RUN] Configuration content:\n")
		fmt.Printf("  Project: %s\n", projectConfig.ProjectName)
		fmt.Printf("  Description: %s\n", projectConfig.ProjectDescription)
        fmt.Printf("  Installed Rules: %v\n", projectConfig.InstalledRules)
        fmt.Printf("  Enabled Agents: %v\n", projectConfig.EnabledAgents)
        return nil
    }

	fmt.Printf("💾 Saving project configuration...\n")
	if err := config.SaveProjectConfig(params.ProjectDir, projectConfig); err != nil {
		return fmt.Errorf("failed to save project configuration: %w", err)
	}

	fmt.Printf("✅ Project configuration saved to .anyagent.yaml\n")
	return nil
}
