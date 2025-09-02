package config

import (
	"embed"
	"fmt"
	"io/fs"
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
