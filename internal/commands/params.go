package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/shibukawa/anyagent/internal/config"
)

// ensureTemplateParameters ensures that all placeholders required by the effective
// AGENTS.md template have corresponding values in the project config. When run in
// interactive mode (non-dry-run), it prompts the user for any missing values and
// saves them to .anyagent/config.yaml. In dry-run, it only reports missing keys.
func ensureTemplateParameters(projectDir string, cfg *config.ProjectConfig, dryRun bool) error {
	// Resolve effective template using unified precedence
	tpl, _ := config.ResolveTemplateContent(projectDir, "AGENTS.md.tmpl", func() (string, error) {
		return config.GetAGENTSTemplate(), nil
	})

	// Extract placeholders from template
	placeholders := config.ExtractTemplateParameters(tpl)
	if len(placeholders) == 0 {
		return nil
	}

	// Build current parameter map (merge of saved parameters and basic fields)
	current := map[string]string{}
	for k, v := range cfg.Parameters {
		current[k] = v
	}
	if cfg.ProjectName != "" {
		current["PROJECT_NAME"] = cfg.ProjectName
	}
	if cfg.ProjectDescription != "" {
		current["PROJECT_DESCRIPTION"] = cfg.ProjectDescription
	}

	// Compute missing keys (exclude special placeholders)
	var missing []string
	for _, key := range placeholders {
		if key == "EXTRA_RULES" {
			continue
		}
		if _, ok := current[key]; !ok {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)

	if len(missing) == 0 {
		return nil
	}

	// Detect interactivity (TTY-like stdin)
	fi, _ := os.Stdin.Stat()
	interactive := (fi.Mode() & os.ModeCharDevice) != 0

	if dryRun || !interactive {
		fmt.Printf("[DRY RUN] Missing template parameters: %s\n", strings.Join(missing, ", "))
		return nil
	}

	// Interactive prompt for missing values
	reader := bufio.NewReader(os.Stdin)
	for _, key := range missing {
		fmt.Printf("Enter %s: ", key)
		v, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				v = ""
			} else {
				return fmt.Errorf("failed to read parameter %s: %w", key, err)
			}
		}
		v = strings.TrimSpace(v)

		// Empty input is allowed (placeholder remains), but we do not save empty to avoid re-prompt loops
		if v == "" {
			fmt.Printf("Warning: %s is empty, will leave placeholder in template\n", key)
			continue
		}

		switch key {
		case "PROJECT_NAME":
			cfg.ProjectName = v
		case "PROJECT_DESCRIPTION":
			cfg.ProjectDescription = v
		default:
			if cfg.Parameters == nil {
				cfg.Parameters = map[string]string{}
			}
			cfg.Parameters[key] = v
		}
	}

	// Persist updated config
	if err := config.SaveProjectConfig(projectDir, cfg); err != nil {
		return fmt.Errorf("failed to save updated project parameters: %w", err)
	}
	return nil
}
