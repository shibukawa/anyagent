package commands

import (
	"github.com/shibukawa/anyagent/internal/config"
)

// shouldCreateCopilotRuleFiles returns true when we should create
// Copilot-specific rule instruction files under .github/instructions.
//
// Rules:
// - If project config exists and enabled agent includes "copilot", return true
// - If project config exists and enabled agent is "codex" (and not copilot), return false
// - If no config or no enabled agents recorded, default to true (backward compatible)
func shouldCreateCopilotRuleFiles(projectDir string) bool {
	cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	if err != nil || len(cfg.EnabledAgents) == 0 {
		// Unknown agent selection → keep previous behavior
		return true
	}
	for _, a := range cfg.EnabledAgents {
		if a == "copilot" {
			return true
		}
	}
	// Explicitly selected agents but none is Copilot → don't create external rule files
	return false
}

// shouldCreateCopilotCommandFiles determines if Copilot command files should be created
func shouldCreateCopilotCommandFiles(projectDir string) bool {
	return shouldCreateCopilotRuleFiles(projectDir)
}

// shouldCreateQDevCommandFiles returns true if Amazon Q Developer is the selected agent
func shouldCreateQDevCommandFiles(projectDir string) bool {
	cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	if err != nil || len(cfg.EnabledAgents) == 0 {
		return false
	}
	for _, a := range cfg.EnabledAgents {
		if a == "qdev" {
			return true
		}
	}
	return false
}

// shouldCreateQDevRuleFiles returns true if Q Developer is the selected agent
func shouldCreateQDevRuleFiles(projectDir string) bool {
	return shouldCreateQDevCommandFiles(projectDir)
}

// shouldCreateCodexCommandFiles returns true if Codex is the selected agent
func shouldCreateCodexCommandFiles(projectDir string) bool {
	cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	if err != nil || len(cfg.EnabledAgents) == 0 {
		return false
	}
	for _, a := range cfg.EnabledAgents {
		if a == "codex" {
			return true
		}
	}
	return false
}

// selectedAgent returns the first enabled agent name from project config (single-agent expected).
func selectedAgent(projectDir string) string {
	cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	if err != nil {
		return ""
	}
	if len(cfg.EnabledAgents) == 0 {
		return ""
	}
	return cfg.EnabledAgents[0]
}

// shouldCreateClaudeCommandFiles returns true if Claude Code is the selected agent
func shouldCreateClaudeCommandFiles(projectDir string) bool {
	cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	if err != nil || len(cfg.EnabledAgents) == 0 {
		return false
	}
	for _, a := range cfg.EnabledAgents {
		if a == "claude" {
			return true
		}
	}
	return false
}

// shouldCreateGeminiCommandFiles returns true if Gemini Code is the selected agent
func shouldCreateGeminiCommandFiles(projectDir string) bool {
	cfg, err := config.LoadProjectConfig(config.GetProjectConfigPath(projectDir))
	if err != nil || len(cfg.EnabledAgents) == 0 {
		return false
	}
	for _, a := range cfg.EnabledAgents {
		if a == "gemini" {
			return true
		}
	}
	return false
}
