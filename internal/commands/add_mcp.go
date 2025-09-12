package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/shibukawa/anyagent/internal/config"
)

// RunAddMCP adds/updates an MCP server definition for this project, records it in .anyagent.yaml,
// writes/updates project mcp.yaml, and creates agent-specific symlinks to the project file.
func RunAddMCP(name, cmdline, projectDir string, dryRun bool, global bool) error {
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

	// Create agent-specific MCP config files based on enabled agent(s)
	if err := ensureMCPFilesForEnabledAgents(projectDir, dryRun); err != nil {
		return err
	}

	// If --global and Codex is selected, write to ~/.codex/config.toml directly for this server
	if global {
		if selectedAgent(projectDir) == "codex" {
			if err := updateCodexMCPMapp(map[string]string{name: cmdline}, dryRun); err != nil {
				return err
			}
		}
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

// ensureMCPFilesForEnabledAgents writes agent-specific MCP config files
func ensureMCPFilesForEnabledAgents(projectDir string, dryRun bool) error {
	cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	if err != nil {
		return nil
	}
	if len(cfg.EnabledAgents) == 0 {
		return nil
	}
	for _, a := range cfg.EnabledAgents {
		if err := ensureMCPFilesForAgent(a, projectDir, cfg.MCPServers, dryRun); err != nil {
			return err
		}
	}
	return nil
}

func ensureMCPFilesForAgent(agentName, projectDir string, servers map[string]string, dryRun bool) error {
	switch agentName {
	case "copilot":
		// VS Code Copilot: .vscode/mcp.json
		path := filepath.Join(projectDir, ".vscode", "mcp.json")
		if dryRun {
			fmt.Printf("[DRY RUN] Would write Copilot MCP config: %s\n", path)
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create .vscode: %w", err)
		}
		data, err := buildCopilotMCPJSON(servers)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
		fmt.Printf("📄 Copilot MCP config generated: %s\n", path)
		return nil
	case "qdev":
		// Amazon Q Developer: .amazonq/mcp.json
		path := filepath.Join(projectDir, ".amazonq", "mcp.json")
		if dryRun {
			fmt.Printf("[DRY RUN] Would write Q Dev MCP config: %s\n", path)
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create .amazonq: %w", err)
		}
		data, err := buildCopilotMCPJSON(servers)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
		fmt.Printf("📄 Q Dev MCP config generated: %s\n", path)
		return nil
	case "claude":
		path := filepath.Join(projectDir, ".claude", "mcp.yaml")
		return writeMCPYAML(path, servers, dryRun)
	case "junie":
		path := filepath.Join(projectDir, ".junie", "mcp.yaml")
		return writeMCPYAML(path, servers, dryRun)
	case "gemini":
		path := filepath.Join(projectDir, ".gemini", "mcp.yaml")
		return writeMCPYAML(path, servers, dryRun)
	case "codex":
		// Do not modify global config automatically; warn if missing and suggest --global
		missing := missingCodexMCPServers(servers)
		if len(missing) > 0 {
			fmt.Printf("⚠️  Some Codex MCP servers are not installed globally: %v\n", missing)
			fmt.Printf("   Enable with: anyagent add mcp <name> --global\n")
		}
		return nil
	default:
		return nil
	}
}

func buildCopilotMCPJSON(servers map[string]string) ([]byte, error) {
	type serverJSON struct {
		Command string            `json:"command"`
		Args    []string          `json:"args,omitempty"`
		Env     map[string]string `json:"env,omitempty"`
	}
	type mcpJSON struct {
		MCPServers map[string]serverJSON `json:"mcpServers"`
	}
	out := mcpJSON{MCPServers: map[string]serverJSON{}}
	for name, cmdline := range servers {
		fields := strings.Fields(cmdline)
		if len(fields) == 0 {
			continue
		}
		sj := serverJSON{Command: fields[0]}
		if len(fields) > 1 {
			sj.Args = fields[1:]
		}
		out.MCPServers[name] = sj
	}
	return json.MarshalIndent(out, "", "  ")
}

func writeMCPYAML(path string, servers map[string]string, dryRun bool) error {
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
	if dryRun {
		fmt.Printf("[DRY RUN] Would write MCP config: %s\n", path)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create dir for MCP: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	fmt.Printf("📄 MCP config generated: %s\n", path)
	return nil
}

// updateCodexMCPConfig writes/updates MCP servers into ~/.codex/config.toml
func updateCodexMCPConfig(servers map[string]string, dryRun bool) error {
	if len(servers) == 0 {
		return nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	codexDir := filepath.Join(homeDir, ".codex")
	codexFile := filepath.Join(codexDir, "config.toml")
	if dryRun {
		fmt.Printf("[DRY RUN] Would update Codex MCP config: %s\n", codexFile)
		return nil
	}
	if err := os.MkdirAll(codexDir, 0755); err != nil {
		return fmt.Errorf("failed to create .codex directory: %w", err)
	}
	var content string
	if b, err := os.ReadFile(codexFile); err == nil {
		content = string(b)
	}
	// Upsert each server section
	for name, cmdline := range servers {
		fields := strings.Fields(cmdline)
		if len(fields) == 0 {
			continue
		}
		cmd := tomlQuote(fields[0])
		var args []string
		for _, a := range fields[1:] {
			args = append(args, tomlQuote(a))
		}
		section := fmt.Sprintf("[mcp_servers.%s]\ncommand = %s\nargs = [%s]\n", name, cmd, strings.Join(args, ", "))
		content = upsertTomlSection(content, "mcp_servers."+name, section)
	}
	return os.WriteFile(codexFile, []byte(content), 0644)
}

// backwards-compatible helper to update a single server map
func updateCodexMCPMapp(servers map[string]string, dryRun bool) error {
	return updateCodexMCPConfig(servers, dryRun)
}

// missingCodexMCPServers returns names that are not present in ~/.codex/config.toml
func missingCodexMCPServers(servers map[string]string) []string {
	var missing []string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Cannot check; assume missing
		for name := range servers {
			missing = append(missing, name)
		}
		return missing
	}
	path := filepath.Join(homeDir, ".codex", "config.toml")
	b, err := os.ReadFile(path)
	if err != nil {
		for name := range servers {
			missing = append(missing, name)
		}
		return missing
	}
	content := string(b)
	for name := range servers {
		header := "[mcp_servers." + name + "]"
		if !strings.Contains(content, header) {
			missing = append(missing, name)
		}
	}
	return missing
}

func tomlQuote(s string) string {
	// TOML basic string quoting
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return "\"" + s + "\""
}

// upsertTomlSection replaces or appends a TOML section by name
func upsertTomlSection(content, sectionName, sectionBody string) string {
	header := "[" + sectionName + "]"
	if content == "" {
		return sectionBody + "\n"
	}
	lines := strings.Split(content, "\n")
	var out []string
	replaced := false
	i := 0
	for i < len(lines) {
		line := lines[i]
		if strings.TrimSpace(line) == header {
			// Skip until next section header or EOF
			// Replace with new section body
			// Move forward to next header
			// First, append sectionBody (without trailing newline split)
			out = append(out, strings.Split(strings.TrimRight(sectionBody, "\n"), "\n")...)
			replaced = true
			i++
			for i < len(lines) {
				if strings.HasPrefix(strings.TrimSpace(lines[i]), "[") {
					break
				}
				i++
			}
			continue
		}
		out = append(out, line)
		i++
	}
	if !replaced {
		if len(out) > 0 && strings.TrimSpace(out[len(out)-1]) != "" {
			out = append(out, "")
		}
		out = append(out, strings.Split(strings.TrimRight(sectionBody, "\n"), "\n")...)
		out = append(out, "")
	}
	return strings.Join(out, "\n")
}
