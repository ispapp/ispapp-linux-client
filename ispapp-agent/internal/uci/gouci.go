package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Tree is the interface for UCI configuration operations
type Tree interface {
	Get(config, section, option string) ([]string, bool)
	Set(config, section, option string, value ...string) bool
	Del(config, section, option string) bool
	AddSection(config, section, sectionType string) error
	Commit() error
	Revert()
}

// uciTree implements the Tree interface
type uciTree struct {
	rootPath string
	configs  map[string]*uciConfig
	changes  bool
}

// uciConfig represents a single UCI configuration file
type uciConfig struct {
	name     string
	sections map[string]*uciSection
	modified bool
}

// uciSection represents a section in a UCI config
type uciSection struct {
	name        string
	sectionType string
	options     map[string][]string
	modified    bool
}

// NewTree creates a new UCI tree with the given root path
func NewTree(rootPath string) Tree {
	return &uciTree{
		rootPath: rootPath,
		configs:  make(map[string]*uciConfig),
	}
}

// Get returns the values for an option in a section of a config
func (t *uciTree) Get(config, section, option string) ([]string, bool) {
	// Load config if needed
	if _, exists := t.configs[config]; !exists {
		if err := t.loadConfig(config); err != nil {
			return nil, false
		}
	}

	// Check if config exists
	cfg, exists := t.configs[config]
	if !exists {
		return nil, false
	}

	// Check if section exists
	sec, exists := cfg.sections[section]
	if !exists {
		return nil, false
	}

	// Check if option exists
	values, exists := sec.options[option]
	return values, exists
}

// Set sets a value for an option in a section of a config
func (t *uciTree) Set(config, section, option string, values ...string) bool {
	// Load config if needed
	if _, exists := t.configs[config]; !exists {
		if err := t.loadConfig(config); err != nil {
			// If config doesn't exist, create it
			t.configs[config] = &uciConfig{
				name:     config,
				sections: make(map[string]*uciSection),
				modified: true,
			}
		}
	}

	// Get or create config
	cfg := t.configs[config]

	// Check if section exists, create it if not
	if _, exists := cfg.sections[section]; !exists {
		cfg.sections[section] = &uciSection{
			name:     section,
			options:  make(map[string][]string),
			modified: true,
		}
		cfg.modified = true
	}

	// Get section
	sec := cfg.sections[section]

	// Set option value
	sec.options[option] = values
	sec.modified = true
	cfg.modified = true
	t.changes = true

	return true
}

// Del deletes an option from a section
func (t *uciTree) Del(config, section, option string) bool {
	// Load config if needed
	if _, exists := t.configs[config]; !exists {
		if err := t.loadConfig(config); err != nil {
			return false
		}
	}

	// Check if config exists
	cfg, exists := t.configs[config]
	if !exists {
		return false
	}

	// Check if section exists
	sec, exists := cfg.sections[section]
	if !exists {
		return false
	}

	// Delete option
	if _, exists := sec.options[option]; exists {
		delete(sec.options, option)
		sec.modified = true
		cfg.modified = true
		t.changes = true
		return true
	}

	return false
}

// AddSection adds a new section to a config
func (t *uciTree) AddSection(config, section, sectionType string) error {
	// Load config if needed
	if _, exists := t.configs[config]; !exists {
		if err := t.loadConfig(config); err != nil {
			// If config doesn't exist, create it
			t.configs[config] = &uciConfig{
				name:     config,
				sections: make(map[string]*uciSection),
				modified: true,
			}
		}
	}

	// Get config
	cfg := t.configs[config]

	// Create section if it doesn't exist
	if _, exists := cfg.sections[section]; !exists {
		cfg.sections[section] = &uciSection{
			name:        section,
			sectionType: sectionType,
			options:     make(map[string][]string),
			modified:    true,
		}
		cfg.modified = true
		t.changes = true
		return nil
	}

	return fmt.Errorf("section already exists")
}

// Commit writes all changes to disk
func (t *uciTree) Commit() error {
	if !t.changes {
		return nil
	}

	for _, cfg := range t.configs {
		if cfg.modified {
			if err := t.saveConfig(cfg); err != nil {
				return err
			}
		}
	}

	t.changes = false
	return nil
}

// Revert discards all changes
func (t *uciTree) Revert() {
	t.configs = make(map[string]*uciConfig)
	t.changes = false
}

// loadConfig loads a UCI config from disk
func (t *uciTree) loadConfig(config string) error {
	configPath := filepath.Join(t.rootPath, config)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	cfg := &uciConfig{
		name:     config,
		sections: make(map[string]*uciSection),
	}

	var currentSection *uciSection
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "config ") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}
			sectionType := parts[1]
			sectionName := strings.Trim(parts[2], "'\"")

			currentSection = &uciSection{
				name:        sectionName,
				sectionType: sectionType,
				options:     make(map[string][]string),
			}
			cfg.sections[sectionName] = currentSection
		} else if strings.HasPrefix(line, "option ") && currentSection != nil {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}
			optionName := parts[1]
			optionValue := strings.Trim(strings.Join(parts[2:], " "), "'\"")
			currentSection.options[optionName] = []string{optionValue}
		} else if strings.HasPrefix(line, "list ") && currentSection != nil {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}
			optionName := parts[1]
			optionValue := strings.Trim(strings.Join(parts[2:], " "), "'\"")

			if existing, ok := currentSection.options[optionName]; ok {
				currentSection.options[optionName] = append(existing, optionValue)
			} else {
				currentSection.options[optionName] = []string{optionValue}
			}
		}
	}

	t.configs[config] = cfg
	return nil
}

// saveConfig saves a UCI config to disk
func (t *uciTree) saveConfig(cfg *uciConfig) error {
	configPath := filepath.Join(t.rootPath, cfg.name)

	// Create a temporary content buffer
	var content strings.Builder

	for _, section := range cfg.sections {
		if section.sectionType != "" {
			fmt.Fprintf(&content, "config '%s' '%s'\n", section.sectionType, section.name)
		} else {
			fmt.Fprintf(&content, "config '%s'\n", section.name)
		}

		for option, values := range section.options {
			if len(values) == 1 {
				fmt.Fprintf(&content, "\toption '%s' '%s'\n", option, values[0])
			} else {
				for _, value := range values {
					fmt.Fprintf(&content, "\tlist '%s' '%s'\n", option, value)
				}
			}
		}
		content.WriteString("\n")
	}

	// Ensure the directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write to file
	return ioutil.WriteFile(configPath, []byte(content.String()), 0644)
}
