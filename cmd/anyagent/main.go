package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/shibukawa/anyagent/internal/commands"
	"github.com/shibukawa/anyagent/internal/config"
)

// CLI represents the command line interface structure
type CLI struct {
	Init   InitCmd   `cmd:"" help:"Setup template editing environment and launch VSCode"`
	Sync   SyncCmd   `cmd:"" help:"Initialize/sync anyagent configuration for a project"`
	Add    AddCmd    `cmd:"" help:"Add additional configurations to the project"`
	Remove RemoveCmd `cmd:"" help:"Remove configurations from the project"`
	List   ListCmd   `cmd:"" help:"List configuration status for the project"`
	Switch SwitchCmd `cmd:"" help:"Switch active AI agent for the project"`
}

// InitCmd represents the init command (template editing environment)
type InitCmd struct {
	DryRun    bool `help:"Show what would be done without actually doing it" short:"n"`
	Force     bool `help:"Force reset all templates to original versions" short:"f"`
	HardReset bool `help:"(deprecated) Same as --force" hidden:""`
}

// SyncCmd represents the sync command (project initialization/sync)
type SyncCmd struct {
	ProjectDir string   `arg:"" optional:"" help:"Project directory (default: current directory)"`
	Agents     []string `help:"AI agents to configure (copilot,qdev,claude,gemini,codex)" short:"a"`
	DryRun     bool     `help:"Show what would be done without actually doing it" short:"n"`
	Force      bool     `help:"Force re-distribute user templates to .anyagent (overwrite if exists)" short:"f"`
}

// AddCmd represents the add command with subcommands
type AddCmd struct {
	Rule    AddRuleCmd    `cmd:"" help:"Add language-specific rules to the project"`
	Command AddCommandCmd `cmd:"" help:"Add VS Code Copilot prompt commands to the project"`
	Mcp     AddMCPCmd     `cmd:"" help:"Add MCP server definition and project wiring"`
}

// RemoveCmd represents the remove command with subcommands
type RemoveCmd struct {
	Rule    RemoveRuleCmd    `cmd:"" help:"Remove language-specific rules from the project"`
	Command RemoveCommandCmd `cmd:"" help:"Remove VS Code Copilot prompt commands from the project"`
}

// ListCmd represents the list command with subcommands
type ListCmd struct {
	Rule    ListRuleCmd    `cmd:"" help:"List language-specific rules status"`
	Command ListCommandCmd `cmd:"" help:"List VS Code Copilot prompt commands status"`
}

// AddRuleCmd represents the add rule subcommand
type AddRuleCmd struct {
	Language   string `arg:"" help:"Language or technology for the rule (e.g., go, typescript, docker)"`
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
	DryRun     bool   `help:"Show what would be done without actually doing it" short:"n"`
}

// AddCommandCmd represents the add command subcommand
type AddCommandCmd struct {
	Command    string `arg:"" optional:"" help:"Command name to add (e.g., create-readme, editorconfig)"`
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
	DryRun     bool   `help:"Show what would be done without actually doing it" short:"n"`
	List       bool   `help:"List available commands" short:"l"`
	Global     bool   `help:"Install to user-global location when applicable (Q Dev, Codex)"`
}

// AddMCPCmd represents the add mcp subcommand
type AddMCPCmd struct {
	Name       string `arg:"" help:"MCP server name (e.g., postgres, filesystem)"`
	Cmd        string `help:"Command to launch the MCP server" required:""`
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
	DryRun     bool   `help:"Show what would be done without actually doing it" short:"n"`
	Global     bool   `help:"Install to user-global location when applicable (Codex)"`
}

// RemoveRuleCmd represents the remove rule subcommand
type RemoveRuleCmd struct {
	Language   string `arg:"" help:"Language or technology for the rule (e.g., go, typescript, docker)"`
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
	DryRun     bool   `help:"Show what would be done without actually doing it" short:"n"`
}

// RemoveCommandCmd represents the remove command subcommand
type RemoveCommandCmd struct {
	Command    string `arg:"" help:"Command name to remove (e.g., create-readme, editorconfig)"`
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
	DryRun     bool   `help:"Show what would be done without actually doing it" short:"n"`
}

// ListRuleCmd represents the list rule subcommand
type ListRuleCmd struct {
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
}

// ListCommandCmd represents the list command subcommand
type ListCommandCmd struct {
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
}

// SwitchCmd represents the switch command
type SwitchCmd struct {
	ProjectDir string `help:"Project directory (default: current directory)" short:"d"`
	Agent      string `arg:"" help:"Target agent (copilot,qdev,claude,gemini,codex)"`
	DryRun     bool   `help:"Show what would be done without actually doing it" short:"n"`
}

// Run executes the init command (template editing environment)
func (cmd *InitCmd) Run() error {
	// Get user config directory
	userConfigDir, err := config.GetUserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}

	fmt.Printf("Using config directory: %s\n", userConfigDir)

	// Run the edit template functionality
	return commands.RunEditTemplate(userConfigDir, cmd.DryRun, cmd.Force || cmd.HardReset)
}

// Run executes the sync command (project initialization/sync)
func (cmd *SyncCmd) Run() error {
	return commands.RunSyncWithOptions(cmd.ProjectDir, cmd.Agents, cmd.DryRun, cmd.Force)
}

// Run executes the add rule command
func (cmd *AddRuleCmd) Run() error {
	return commands.RunAddRule(cmd.Language, cmd.ProjectDir, cmd.DryRun)
}

// Run executes the add command subcommand
func (cmd *AddCommandCmd) Run() error {
	if cmd.List {
		return commands.ListAvailableCommands()
	}

	if cmd.Command == "" {
		return commands.ListAvailableCommands()
	}

	return commands.RunAddCommand(cmd.Command, cmd.ProjectDir, cmd.DryRun, cmd.Global)
}

// Run executes the add mcp subcommand
func (cmd *AddMCPCmd) Run() error {
	return commands.RunAddMCP(cmd.Name, cmd.Cmd, cmd.ProjectDir, cmd.DryRun, cmd.Global)
}

// Run executes the remove rule command
func (cmd *RemoveRuleCmd) Run() error {
	return commands.RunRemoveRule(cmd.Language, cmd.ProjectDir, cmd.DryRun)
}

// Run executes the remove command subcommand
func (cmd *RemoveCommandCmd) Run() error {
	return commands.RunRemoveCommand(cmd.Command, cmd.ProjectDir, cmd.DryRun)
}

// Run executes the list rule command
func (cmd *ListRuleCmd) Run() error {
	return commands.RunListRules(cmd.ProjectDir)
}

// Run executes the list command subcommand
func (cmd *ListCommandCmd) Run() error {
	return commands.RunListCommands(cmd.ProjectDir)
}

// Run executes the switch command
func (cmd *SwitchCmd) Run() error {
	return commands.RunSwitch(cmd.ProjectDir, cmd.Agent, cmd.DryRun)
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("anyagent"),
		kong.Description("AI agent configuration management tool"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	err := ctx.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
