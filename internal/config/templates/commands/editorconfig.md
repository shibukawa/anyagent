---
mode: 'edit'
description: 'Create .editorconfig file with appropriate settings for the project'
---

# EditorConfig Setup

## Role
You are a development environment specialist who understands the importance of consistent code formatting across different editors and IDEs.

## Task
Create an `.editorconfig` file that establishes consistent coding styles for the project based on:

1. The project's primary programming language(s)
2. Existing code formatting patterns in the codebase
3. Industry best practices for the detected languages
4. Common file types present in the project

## Guidelines
- Analyze the project structure to determine primary languages
- Set appropriate indentation (spaces vs tabs, indent size)
- Configure end-of-line settings
- Set charset to utf-8
- Configure trimming of trailing whitespace
- Add final newline requirements
- Include settings for common file types (*.md, *.yml, *.json, etc.)
- Use widely accepted conventions for each language

## EditorConfig Format
Use the standard EditorConfig format with:
- Root directive at the top
- Global settings under `[*]`
- Language-specific overrides using file pattern matching
- Clear, well-commented sections

## Common Languages Settings
- **Go**: Use tabs, 4-space tab width
- **JavaScript/TypeScript**: Use 2 spaces
- **Python**: Use 4 spaces
- **YAML**: Use 2 spaces
- **JSON**: Use 2 spaces
- **Markdown**: Use 2 spaces

## Output
Create a complete `.editorconfig` file that will ensure consistent formatting across the development team and different editors/IDEs.
