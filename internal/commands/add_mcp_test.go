package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/anyagent/internal/config"
)

func TestRunAddMCP_CreatesProjectMCPAndAgentFiles(t *testing.T) {
	dir := t.TempDir()
	// Initialize project: AGENTS.md + config with copilot enabled
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# AGENTS"), 0644); err != nil {
		t.Fatalf("failed to write AGENTS.md: %v", err)
	}
	if err := config.SaveProjectConfig(dir, &config.ProjectConfig{EnabledAgents: []string{"copilot"}}); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	if err := RunAddMCP("postgres", "npx -y @modelcontextprotocol/server-postgres postgresql://localhost/mydb", dir, false, false); err != nil {
		t.Fatalf("RunAddMCP failed: %v", err)
	}

	// Project mcp.yaml should exist
	if _, err := os.Stat(filepath.Join(dir, "mcp.yaml")); os.IsNotExist(err) {
		t.Fatalf("mcp.yaml was not created in project root")
	}
	// Copilot config should exist at .vscode/mcp.json
	if _, err := os.Stat(filepath.Join(dir, ".vscode", "mcp.json")); os.IsNotExist(err) {
		t.Fatalf(".vscode/mcp.json was not created for Copilot")
	}
}
