# AnyAgent Template Configuration

Welcome to your anyagent template editing environment! This directory contains the templates for managing AI agent settings across your projects.

## Directory Structure

```
├── AGENTS.md                 # anyagent project configuration (for this template repo)
├── templates/                # 🎯 Main templates (this is what gets deployed to projects)
│   ├── AGENTS.md.tmpl        # Project configuration template with placeholders
│   ├── mcp.yaml              # MCP server definitions template
│   └── commands/             # AI instruction templates (add/remove freely)
└── Agent directories/        # ⚙️ Development helpers (not deployed to projects)
    ├── .github/              # Makes this repo easy to edit with GitHub Copilot
    ├── .amazonq/             # Makes this repo easy to edit with Amazon Q
    ├── .claude/              # Makes this repo easy to edit with Claude
    └── .gemini/              # Makes this repo easy to edit with Gemini Code
```

**Important**: The main deliverable is the `templates/` directory. The agent directories (`.github/`, `.amazonq/`, etc.) are only development helpers to make editing these templates easier with AI assistants - they are not copied to your projects.

## Quick Start

1. **Edit Templates**: Modify files in the `templates/` directory to customize how projects are configured
2. **Customize Instructions**: Edit files in `templates/commands/` to define how AI assistants should behave  
3. **Test Changes**: Use `anyagent sync` in a test project to see how your templates work

## Key Files to Edit

- **`templates/AGENTS.md.tmpl`**: Main project configuration template with placeholders
- **`templates/commands/`**: Arbitrary command templates (name and count are up to you)
- **`templates/extra_rules/`**: Arbitrary stack-specific rules (name and count are up to you)
- **`templates/mcp.yaml`**: Model Context Protocol server configurations (optional)

## Command Template Guide

- Location: `templates/commands/<command-name>.md`
- File name: use kebab-case (e.g., `create-readme.md`, `editorconfig.md`)
- Structure:
  1. Optional YAML frontmatter (used by GitHub Copilot prompts)
  2. Main markdown body (instructions executed by the assistant)

Example:
```
---
title: "Create Project README"
description: "Generate a comprehensive README based on repo content"
tags: ["documentation", "readme"]
---

# Create README

Goal: Generate a README.md for this repository.

Steps:
- Inspect top-level files, modules, and main entry points
- Summarize purpose, setup, usage, and configuration
- Include examples and common tasks
```

How anyagent uses it:
- `anyagent add command <name>`
  - VS Code Copilot: creates `.github/prompts/<name>.prompt.md` with the same content
  - Amazon Q Developer: creates `~/.aws/amazonq/prompts/<name with spaces>.md` with YAML frontmatter removed
- Listing: `anyagent add command -l` shows available templates based on files in `templates/commands`

Codex integration:
- Codex (chatgpt.com/codex) supports slash commands via `~/.codex/prompts/<name>.md`.
- `anyagent add command <name>` installs Codex prompts when Codex is the selected agent.
- YAML frontmatter is removed for Codex prompts (Markdown body only).
- Extra rules continue to be merged into `AGENTS.md`.

## Using Placeholders

When editing `AGENTS.md.tmpl`, you can use these placeholders:
- `{{PROJECT_NAME}}` - Replaced with the actual project name
- `{{PROJECT_DESCRIPTION}}` - Project description
- `{{PRIMARY_LANGUAGE}}` - Main programming language
- `{{CODE_STYLE}}` - Coding style guidelines
- `{{EXTRA_RULES}}` - Special placeholder. At generation time, anyagent concatenates all installed extra rule documents and injects the merged content here. Project-local templates under `.anyagent/extra_rules/` are used if present; otherwise embedded defaults are used. Agents like Codex rely on this region to consume stack-specific rules.

Amazon Q Developer integration:
- Reads rules from project `.amazonq/rules/`
- Link `AGENTS.md` as `.amazonq/rules/AGENTS.md` and place per-language extra rules there

## Next Steps

1. **Focus on Templates**: Customize the files in `templates/` directory to match your development preferences
2. **Test Your Changes**: Run `anyagent sync` in a test project to see how your templates work
3. **Use AI Helpers**: The agent directories in this repo help you edit templates with AI assistants

---

💡 **Tip**: This directory is itself an anyagent project, so you can use AI assistants to help edit these templates! The agent directories (`.github/`, `.amazonq/`, etc.) exist solely to make this template editing experience better.
Gemini integration:
- Reads `AGENTS.md` directly (no symlink required)
- Extra rules are merged into `AGENTS.md`
