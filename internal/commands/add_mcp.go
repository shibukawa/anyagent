package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "gopkg.in/yaml.v3"

    "github.com/shibukawa/anyagent/internal/config"
)

// RunAddMCP adds/updates an MCP server definition for this project, records it in .anyagent.yaml,
// writes/updates project mcp.yaml, and creates agent-specific symlinks to the project file.
func RunAddMCP(name, cmdline, projectDir string, dryRun bool) error {
    if name == "" {
        return fmt.Errorf("MCP server name cannot be empty")
    }
    if strings.TrimSpace(cmdline) == "" {
        return fmt.Errorf("--cmd cannot be empty")
    }

    // Resolve project directory
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

    // Must be initialized (AGENTS.md present)
    if _, err := os.Stat(filepath.Join(projectDir, "AGENTS.md")); os.IsNotExist(err) {
        return fmt.Errorf("project is not initialized with anyagent. Run 'anyagent sync' first")
    }

    fmt.Printf("Adding MCP server '%s' to project...\n", name)

    // Update .anyagent.yaml
    cfgPath := config.GetProjectConfigPath(projectDir)
    cfg, err := config.LoadProjectConfig(cfgPath)
    if err != nil {
        return fmt.Errorf("failed to load project config: %w", err)
    }
    if cfg.MCPServers == nil {
        cfg.MCPServers = map[string]string{}
    }
    cfg.MCPServers[name] = cmdline
    if dryRun {
        fmt.Printf("[DRY RUN] Would record MCP server '%s' in .anyagent.yaml with cmd: %s\n", name, cmdline)
    } else {
        if err := cfg.Save(cfgPath); err != nil {
            return fmt.Errorf("failed to save project config: %w", err)
        }
        fmt.Printf("💾 Recorded MCP server '%s' in .anyagent.yaml\n", name)
    }

    // Ensure mcp.yaml in project root is updated
    if err := writeOrUpdateProjectMCP(projectDir, cfg.MCPServers, dryRun); err != nil {
        return err
    }

    // Create agent-specific symlink(s) to project mcp.yaml based on enabled agent(s)
    if err := ensureMCPSymlinksForEnabledAgents(projectDir, dryRun); err != nil {
        return err
    }

    fmt.Printf("✅ MCP server '%s' added/updated successfully\n", name)
    return nil
}

// writeOrUpdateProjectMCP materializes mcp.yaml aggregating servers from config.
func writeOrUpdateProjectMCP(projectDir string, servers map[string]string, dryRun bool) error {
    // Convert to YAML structure similar to templates/mcp.yaml
    type serverDef struct {
        Command string   `yaml:"command"`
        Args    []string `yaml:"args,omitempty"`
    }
    type mcpConfig struct {
        Servers map[string]serverDef `yaml:"servers"`
    }
    out := mcpConfig{Servers: map[string]serverDef{}}
    for name, cmdline := range servers {
        fields := strings.Fields(cmdline)
        if len(fields) == 0 {
            continue
        }
        sd := serverDef{Command: fields[0]}
        if len(fields) > 1 {
            sd.Args = fields[1:]
        }
        out.Servers[name] = sd
    }

    data, err := yaml.Marshal(&out)
    if err != nil {
        return fmt.Errorf("failed to marshal mcp.yaml: %w", err)
    }
    path := filepath.Join(projectDir, "mcp.yaml")
    if dryRun {
        fmt.Printf("[DRY RUN] Would write project MCP config: %s\n", path)
        return nil
    }
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("failed to write %s: %w", path, err)
    }
    fmt.Printf("📄 Project MCP config updated: mcp.yaml\n")
    return nil
}

// ensureMCPSymlinksForEnabledAgents creates symlinks from agent config folders to project mcp.yaml
func ensureMCPSymlinksForEnabledAgents(projectDir string, dryRun bool) error {
    cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
    if err != nil {
        return nil
    }
    if len(cfg.EnabledAgents) == 0 {
        return nil
    }
    for _, a := range cfg.EnabledAgents {
        if err := ensureMCPSymlinkForAgent(a, projectDir, dryRun); err != nil {
            return err
        }
    }
    return nil
}

func ensureMCPSymlinkForAgent(agentName, projectDir string, dryRun bool) error {
    var symlinkPath, targetRel string
    switch agentName {
    case "copilot":
        symlinkPath = filepath.Join(projectDir, ".github", "mcp.yaml")
        targetRel = "../mcp.yaml"
    case "qdev":
        symlinkPath = filepath.Join(projectDir, ".amazonq", "mcp.yaml")
        targetRel = "../mcp.yaml"
    case "claude":
        symlinkPath = filepath.Join(projectDir, ".claude", "mcp.yaml")
        targetRel = "../mcp.yaml"
    case "junie":
        symlinkPath = filepath.Join(projectDir, ".junie", "mcp.yaml")
        targetRel = "../mcp.yaml"
    case "gemini":
        symlinkPath = filepath.Join(projectDir, ".gemini", "mcp.yaml")
        targetRel = "../mcp.yaml"
    case "codex":
        // No project symlink for Codex (home-based prompts); not applicable
        return nil
    default:
        return nil
    }

    // Ensure directory exists
    if !dryRun {
        if err := os.MkdirAll(filepath.Dir(symlinkPath), 0755); err != nil {
            return fmt.Errorf("failed to create agent dir: %w", err)
        }
    }
    if dryRun {
        fmt.Printf("[DRY RUN] Would create symlink: %s -> %s\n", symlinkPath, targetRel)
        return nil
    }
    // Remove existing
    _ = os.Remove(symlinkPath)
    if err := os.Symlink(targetRel, symlinkPath); err != nil {
        return fmt.Errorf("failed to create symlink %s: %w", symlinkPath, err)
    }
    fmt.Printf("🔗 MCP symlink for %s: %s -> %s\n", agentName, symlinkPath, targetRel)
    return nil
}

