package config

import (
    "embed"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
)

//go:embed templates/commands/*
var commandsFS embed.FS

// GetCommandTemplate retrieves the template content for the specified command
func GetCommandTemplate(command string) (string, error) {
    fileName := fmt.Sprintf("templates/commands/%s.md", command)
    content, err := commandsFS.ReadFile(fileName)
    if err != nil {
        return "", fmt.Errorf("command template not found: %s", command)
    }
    return string(content), nil
}

// ResolveTemplateContent reads a template with precedence:
// 1) projectDir/.anyagent/<relPath>
// 2) <userConfigDir>/templates/<relPath>
// 3) embeddedFallback()
func ResolveTemplateContent(projectDir string, relPath string, embeddedFallback func() (string, error)) (string, error) {
    // Determine project base
    base := projectDir
    if base == "" {
        if wd, err := os.Getwd(); err == nil {
            base = wd
        }
    }
    // Project override
    if base != "" {
        if b, err := os.ReadFile(filepath.Join(base, ".anyagent", relPath)); err == nil {
            return string(b), nil
        }
    }
    // User override
    if userDir, err := GetUserConfigDir(); err == nil {
        if b, err := os.ReadFile(filepath.Join(userDir, "templates", relPath)); err == nil {
            return string(b), nil
        }
    }
    // Embedded fallback
    return embeddedFallback()
}

// GetCommandTemplateResolved resolves a command template using standard precedence.
func GetCommandTemplateResolved(projectDir, command string) (string, error) {
    rel := filepath.Join("commands", fmt.Sprintf("%s.md", command))
    return ResolveTemplateContent(projectDir, rel, func() (string, error) {
        return GetCommandTemplate(command)
    })
}

// GetAvailableCommands returns a list of available command templates
func GetAvailableCommands() ([]string, error) {
	var commands []string

	entries, err := fs.ReadDir(commandsFS, "templates/commands")
	if err != nil {
		return nil, fmt.Errorf("failed to read commands directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			// Remove .md extension to get command name
			command := strings.TrimSuffix(entry.Name(), ".md")
			commands = append(commands, command)
		}
	}

	return commands, nil
}
