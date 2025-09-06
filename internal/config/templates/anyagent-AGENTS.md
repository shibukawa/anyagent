# AI Agents Configuration - AnyAgent Template Environment

## Project Information
- name: anyagent
- description: AI agent configuration management tool – Template Editing Environment
- version: 1.0.0

## Template Editing Guidelines

This directory is the editing environment for the template files managed by the anyagent tool.

### Files to Edit

#### 1. templates/AGENTS.md.tmpl
- Purpose: Base settings and policies applied across all projects
- Contents: Project information template and the basic structure for agent settings
- Editing policy: Write project‑specific items with `{{PLACEHOLDER}}` format

#### 2. templates/commands/
- Purpose: Define task‑specific prompts and instructions
- Structure: Arbitrary (add, remove, and name files freely per project needs)
- Editing policy: Split by task type and keep content reusable
- Note: General development rules belong in `AGENTS.md.tmpl`

#### 3. templates/extra_rules/
- Purpose: Detailed, stack‑specific rules for languages, frameworks, tools, etc.
- Structure: Arbitrary (add, remove, and name files freely per project needs)
- Editing policy: Write team‑agreed guidelines and subdivide as needed

#### 4. templates/mcp.yaml
- Purpose: Template for MCP (Model Context Protocol) server settings
- Editing policy: Define MCP settings reusable across projects

### Important Notes

⚠️ Scope of edits
- Edit only files under `templates/`
- Files under `.github/` are not targets for editing (they are agent‑specific placement locations)
- Edit only the templates that are applied to projects by running the `anyagent` commands

✅ Template hierarchy
1. **AGENTS.md.tmpl**: Base settings and general development rules shared across projects
2. **commands/**: Task‑specific prompts (coding, documentation, reviews, etc.)
3. **extra_rules/**: Detailed, stack‑specific rules (create only those you need)

### Template Parameter System

🔧 Dynamic parameters
- You can use placeholders in the form `{{PARAMETER_NAME}}` inside templates
- They are detected automatically when running `anyagent sync` and you will be prompted for values
- Entered parameters are saved to `.anyagent/config.yaml` and reused on regeneration

Example:
```markdown
# Project: {{PROJECT_NAME}}
Description: {{PROJECT_DESCRIPTION}}
Primary Language: {{PRIMARY_LANGUAGE}}
Team: {{TEAM_NAME}}
```

Supported parameter patterns:
- `{{PROJECT_NAME}}` – Project name
- `{{PROJECT_DESCRIPTION}}` – Project description
- `{{PRIMARY_LANGUAGE}}` – Primary development language
- `{{any_valid_name}}` – Custom parameters (alphanumeric and underscore)

Saving parameters:
- Entered parameters are saved per project to `.anyagent/config.yaml`
- They are reused automatically on `anyagent sync` and subsequent generations
- You can edit `.anyagent/config.yaml` manually to change parameter values

## Agent Specific Settings

### GitHub Copilot
- enabled: true
- instructions_file: .github/copilot‑instructions.md
- note: Not an editing target (managed/placed automatically by anyagent)

### Amazon Q Developer
- enabled: true
- config_file: .amazonq/config.json
- note: Not an editing target (managed/placed automatically by anyagent)

### Claude Code
- enabled: true
- config_file: .claude/config.json
- note: Not an editing target (managed/placed automatically by anyagent)

### IntelliJ IDEA Junie
- enabled: true
- config_file: .junie/settings.json
- note: Not an editing target (managed/placed automatically by anyagent)

## MCP Server Configuration

### Gemini Code
- enabled: true
- mcp_servers:
  - filesystem
  - git

## Template Development Workflow

1. Define project‑wide rules and general development guidance in **templates/AGENTS.md.tmpl**
2. Create/update task‑specific prompts in **templates/commands/**
3. Elaborate stack‑specific rules in **templates/extra_rules/**
4. Apply changes to projects with the `anyagent sync` command

Edits made in this environment are applied to real projects through the `anyagent` CLI.

