# AnyAgent Template Configuration

Welcome to your anyagent template editing environment! This directory contains the templates for managing AI agent settings across your projects.

## Directory Structure

```
├── AGENTS.md                 # anyagent project configuration (for this template repo)
├── templates/                # 🎯 Main templates (this is what gets deployed to projects)
│   ├── AGENTS.md.tmpl        # Project configuration template with placeholders
│   ├── mcp.yaml              # MCP server definitions template
│   └── commands/             # AI instruction templates
│       ├── general.md        # General AI assistant instructions
│       ├── coding.md         # Coding-specific instructions
│       └── project-specific.md # Project-specific instructions
└── Agent directories/        # ⚙️ Development helpers (not deployed to projects)
    ├── .github/              # Makes this repo easy to edit with GitHub Copilot
    ├── .amazonq/             # Makes this repo easy to edit with Amazon Q
    ├── .claude/              # Makes this repo easy to edit with Claude
    ├── .junie/               # Makes this repo easy to edit with IntelliJ IDEA Junie
    └── .gemini/              # Makes this repo easy to edit with Gemini Code
```

**Important**: The main deliverable is the `templates/` directory. The agent directories (`.github/`, `.amazonq/`, etc.) are only development helpers to make editing these templates easier with AI assistants - they are not copied to your projects.

## Quick Start

1. **Edit Templates**: Modify files in the `templates/` directory to customize how projects are configured
2. **Customize Instructions**: Edit files in `templates/commands/` to define how AI assistants should behave  
3. **Test Changes**: Use `anyagent init` in a test project to see how your templates work

## Key Files to Edit

- **`templates/AGENTS.md.tmpl`**: Main project configuration template with placeholders
- **`templates/commands/coding.md`**: Instructions for coding-related AI assistance
- **`templates/commands/general.md`**: General AI assistant behavior guidelines
- **`templates/mcp.yaml`**: Model Context Protocol server configurations

## Using Placeholders

When editing `AGENTS.md.tmpl`, you can use these placeholders:
- `{{PROJECT_NAME}}` - Will be replaced with the actual project name
- `{{PROJECT_DESCRIPTION}}` - Project description
- `{{PRIMARY_LANGUAGE}}` - Main programming language
- `{{CODE_STYLE}}` - Coding style guidelines

## Next Steps

1. **Focus on Templates**: Customize the files in `templates/` directory to match your development preferences
2. **Test Your Changes**: Run `anyagent init` in a test project to see how your templates work
3. **Use AI Helpers**: The agent directories in this repo help you edit templates with AI assistants

---

💡 **Tip**: This directory is itself an anyagent project, so you can use AI assistants to help edit these templates! The agent directories (`.github/`, `.amazonq/`, etc.) exist solely to make this template editing experience better.
