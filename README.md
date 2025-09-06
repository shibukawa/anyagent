## anyagent

Unified configuration manager for multiple AI coding assistants.

## Overview

anyagent is a CLI tool that manages project configuration for multiple AI assistants, including GitHub Copilot, Amazon Q Developer, Claude Code, Gemini Code, and ChatGPT Codex. It composes a single AGENTS.md and generates agent‑specific files as needed.

## Features

- 🤖 Multi‑agent: Copilot, Q Dev, Claude, Gemini, Codex
- 📝 Single‑source config: generates agent‑specific files from one set of templates
- 🔧 Extra rules: stack‑specific rules merged into AGENTS.md via `{{EXTRA_RULES}}`
- 📋 Custom commands: per‑agent command files (global install for Q Dev/Codex)
- 🎯 Template system: reusable templates under `.anyagent/`
- 📊 Status & management: add/remove/list for rules and commands

## Installation

```bash
go install github.com/shibukawa/anyagent/cmd/anyagent@latest
```

## Quick Start

```bash
# Launch template editing environment (once)
anyagent init

# Initialize/sync project (creates .anyagent/ on first run)
anyagent sync

# Add rules
anyagent add rule go
anyagent add rule typescript

# Add commands
anyagent add command review-code

# Show status
anyagent list rule
anyagent list command
```

## Init / Sync / Switch

```bash
anyagent init                       # Launch template editing env (user side)
anyagent sync [directory]           # Distribute .anyagent/ and generate AGENTS.md
anyagent switch <agent>             # Switch active agent (links/commands refreshed)

# Options
#   --force, -f   Overwrite existing .anyagent/ on sync
#   --dry-run, -n Preview actions only
```

## Rule Management

```bash
anyagent add rule <language>
anyagent remove rule <language>
anyagent list rule
```

## Command Management

```bash
anyagent add command <name> [--global]  # Q Dev/Codex: use --global to install in user folder
anyagent remove command <name>
anyagent list command
```

## MCP Servers

```bash
anyagent add mcp <name> --cmd "<launcher and args>" [--global]
# e.g. anyagent add mcp context7 --cmd "npx -y @upstash/context7-mcp@latest"
```

## Supported AI Agents

### GitHub Copilot
- Rules: `.github/instructions/<language>.instructions.md`
- Commands: `.github/prompts/<command>.prompt.md`
- Integration: `.github/copilot-instructions.md` → `../AGENTS.md`
- Format: Markdown with YAML frontmatter
- MCP: `.vscode/mcp.json` (JSON in project)

### Amazon Q Developer
- Rules: reads `.amazonq/rules/` (link AGENTS.md as `.amazonq/rules/AGENTS.md`)
- Commands: `~/.aws/amazonq/prompts/<command>.md` (global, install via `--global`)
- Format: plain Markdown (frontmatter removed)
- MCP: `.amazonq/mcp.json` (JSON in project)

### Claude Code
- Rules: merged in AGENTS.md
- Integration: `CLAUDE.md` → `AGENTS.md`
- Commands: `.claude/commands/<command>.md` (project‑local)
- MCP: `.claude/mcp.yaml` (YAML in project)

### Gemini Code
- Rules: merged in AGENTS.md
- Commands: `.gemini/commands/<command>.toml` (TOML with `description` and `prompt`)
- Integration: no symlink (reads AGENTS.md directly)
- MCP: `.gemini/mcp.yaml` (YAML in project)

### ChatGPT Codex
- Rules: merged in AGENTS.md
- Commands: `~/.codex/prompts/<command>.md` (global, install via `--global`)
- Integration: AGENTS.md only
- MCP: `~/.codex/config.toml` (`[mcp_servers.<name>]`, install via `--global`)
  - sync/switch doesn’t overwrite global config; missing items are warned with an activation hint.

## Configuration Files

### Project config (`.anyagent/config.yaml`)

```yaml
project_name: "myproject"
project_description: "My awesome project"
installed_rules:
  - go
  - typescript
installed_commands:
  - create-readme
enabled_agents:
  - copilot
mcp_servers:
  context7: "npx -y @upstash/context7-mcp@latest"
```

### AGENTS.md (composed)
- Unified settings for all agents
- Project info and rules
- Injects concatenated extra rules at `{{EXTRA_RULES}}`
- Agents like Codex read this single file directly

## Development

### Build
```bash
go build -o anyagent cmd/anyagent/main.go
```

### Test
```bash
go test ./...
```

## License

AGPL-3.0. See `LICENSE`.

## Contributing

See README.ja.md for the Japanese version. Contribution guidelines will be added later.
