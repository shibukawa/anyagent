package config

import (
	"testing"
)

func TestExtractTemplateParameters(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected []string
	}{
		{
			name:     "no parameters",
			template: "Hello World",
			expected: []string{},
		},
		{
			name:     "single parameter",
			template: "Project: {{PROJECT_NAME}}",
			expected: []string{"PROJECT_NAME"},
		},
		{
			name:     "multiple parameters",
			template: "Project: {{PROJECT_NAME}}\nDescription: {{PROJECT_DESCRIPTION}}\nLanguage: {{PRIMARY_LANGUAGE}}",
			expected: []string{"PRIMARY_LANGUAGE", "PROJECT_DESCRIPTION", "PROJECT_NAME"}, // sorted
		},
		{
			name:     "duplicate parameters",
			template: "{{PROJECT_NAME}} - {{PROJECT_NAME}} project",
			expected: []string{"PROJECT_NAME"},
		},
		{
			name:     "mixed case and numbers",
			template: "{{TEAM_NAME}} {{VERSION_2}} {{API_KEY_V3}}",
			expected: []string{"API_KEY_V3", "TEAM_NAME", "VERSION_2"},
		},
		{
			name:     "invalid patterns should be ignored",
			template: "{{lowercase}} {{123INVALID}} {{_INVALID}} {{PROJECT-NAME}}",
			expected: []string{}, // all should be ignored due to validation rules
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTemplateParameters(tt.template)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d parameters, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			for i, param := range result {
				if param != tt.expected[i] {
					t.Errorf("Expected parameter %d to be '%s', got '%s'", i, tt.expected[i], param)
				}
			}
		})
	}
}

func TestReplaceTemplateParameters(t *testing.T) {
	tests := []struct {
		name       string
		template   string
		parameters map[string]string
		expected   string
	}{
		{
			name:       "no parameters to replace",
			template:   "Hello World",
			parameters: map[string]string{"PROJECT_NAME": "test"},
			expected:   "Hello World",
		},
		{
			name:     "single parameter replacement",
			template: "Project: {{PROJECT_NAME}}",
			parameters: map[string]string{
				"PROJECT_NAME": "MyProject",
			},
			expected: "Project: MyProject",
		},
		{
			name:     "multiple parameter replacement",
			template: "Project: {{PROJECT_NAME}}\nDescription: {{PROJECT_DESCRIPTION}}",
			parameters: map[string]string{
				"PROJECT_NAME":        "MyProject",
				"PROJECT_DESCRIPTION": "A test project",
			},
			expected: "Project: MyProject\nDescription: A test project",
		},
		{
			name:     "missing parameter leaves placeholder",
			template: "Project: {{PROJECT_NAME}}\nTeam: {{TEAM_NAME}}",
			parameters: map[string]string{
				"PROJECT_NAME": "MyProject",
			},
			expected: "Project: MyProject\nTeam: {{TEAM_NAME}}",
		},
		{
			name:     "duplicate parameters replaced consistently",
			template: "{{PROJECT_NAME}} - {{PROJECT_NAME}} project",
			parameters: map[string]string{
				"PROJECT_NAME": "TestApp",
			},
			expected: "TestApp - TestApp project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceTemplateParameters(tt.template, tt.parameters)
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
