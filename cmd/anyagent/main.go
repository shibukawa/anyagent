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
	EditTemplate EditTemplateCmd `cmd:"" help:"Setup template editing environment and launch VSCode"`
	Init         InitCmd         `cmd:"" help:"Initialize anyagent configuration for a project"`
	Add          AddCmd          `cmd:"" help:"Add additional configurations to the project"`
	Remove       RemoveCmd       `cmd:"" help:"Remove configurations from the project"`
	List         ListCmd         `cmd:"" help:"List configuration status for the project"`
}

// EditTemplateCmd represents the edit-template command
type EditTemplateCmd struct {
	DryRun    bool `help:"Show what would be done without actually doing it" short:"n"`
	HardReset bool `help:"Force reset all templates to original versions" short:"r"`
}

// InitCmd represents the init command
type InitCmd struct {
	ProjectDir string   `arg:"" optional:"" help:"Project directory (default: current directory)"`
	Agents     []string `help:"AI agents to configure (copilot,qdev,claude,junie,gemini)" short:"a"`
	DryRun     bool     `help:"Show what would be done without actually doing it" short:"n"`
}

// AddCmd represents the add command with subcommands
type AddCmd struct {
	Rule    AddRuleCmd    `cmd:"" help:"Add language-specific rules to the project"`
	Command AddCommandCmd `cmd:"" help:"Add VS Code Copilot prompt commands to the project"`
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

// Run executes the edit-template command
func (cmd *EditTemplateCmd) Run() error {
	// Get user config directory
	userConfigDir, err := config.GetUserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}

	fmt.Printf("Using config directory: %s\n", userConfigDir)

	// Run the edit template functionality
	return commands.RunEditTemplate(userConfigDir, cmd.DryRun, cmd.HardReset)
}

// Run executes the init command
func (cmd *InitCmd) Run() error {
	return commands.RunInit(cmd.ProjectDir, cmd.Agents, cmd.DryRun)
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

	return commands.RunAddCommand(cmd.Command, cmd.ProjectDir, cmd.DryRun)
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
