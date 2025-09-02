package config

import (
	"regexp"
	"sort"
)

// ExtractTemplateParameters extracts all {{PARAMETER}} placeholders from a template string
func ExtractTemplateParameters(template string) []string {
	// Regular expression to match {{PARAMETER_NAME}} patterns
	re := regexp.MustCompile(`\{\{([A-Z][A-Z0-9_]*)\}\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	// Use a map to avoid duplicates
	paramMap := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			paramMap[match[1]] = true
		}
	}

	// Convert to sorted slice for consistent ordering
	var params []string
	for param := range paramMap {
		params = append(params, param)
	}
	sort.Strings(params)

	return params
}

// ReplaceTemplateParameters replaces all {{PARAMETER}} placeholders with values from the map
func ReplaceTemplateParameters(template string, parameters map[string]string) string {
	re := regexp.MustCompile(`\{\{([A-Z][A-Z0-9_]*)\}\}`)

	return re.ReplaceAllStringFunc(template, func(match string) string {
		// Extract parameter name from {{PARAM_NAME}}
		paramMatch := re.FindStringSubmatch(match)
		if len(paramMatch) > 1 {
			paramName := paramMatch[1]
			if value, exists := parameters[paramName]; exists {
				return value
			}
		}
		// Return original placeholder if no replacement found
		return match
	})
}
