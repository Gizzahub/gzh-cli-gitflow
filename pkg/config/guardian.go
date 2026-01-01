package config

import (
	"fmt"
	"regexp"
)

// GuardianConfig defines Guardian mode settings
type GuardianConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Mode     string        `yaml:"mode"` // "strict" or "permissive"
	Naming   NamingRule    `yaml:"naming"`
	Workflow WorkflowRules `yaml:"workflow"`
}

// NamingRule defines naming constraints for branches
type NamingRule struct {
	Pattern   string   `yaml:"pattern"`    // Regex pattern
	MaxLength int      `yaml:"max_length"` // Maximum branch name length
	Forbidden []string `yaml:"forbidden"`  // Forbidden words/patterns

	// compiled is the cached regex pattern
	compiled *regexp.Regexp
}

// WorkflowRules defines workflow constraints
type WorkflowRules struct {
	RequireCleanTree     bool `yaml:"require_clean_tree"`
	RequireUpToDate      bool `yaml:"require_up_to_date"`
	PreventDirectPush    bool `yaml:"prevent_direct_push"`
	RequireLinearHistory bool `yaml:"require_linear_history"`
}

// getPattern returns the compiled regex pattern, compiling and caching it if necessary
func (nr *NamingRule) getPattern() (*regexp.Regexp, error) {
	if nr.compiled != nil {
		return nr.compiled, nil
	}
	if nr.Pattern == "" {
		return nil, nil
	}
	re, err := regexp.Compile(nr.Pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid naming pattern: %w", err)
	}
	nr.compiled = re
	return re, nil
}

// Validate validates a branch name against the naming rule
func (nr *NamingRule) Validate(name string) error {
	// Check length
	if nr.MaxLength > 0 && len(name) > nr.MaxLength {
		return fmt.Errorf("branch name exceeds maximum length of %d characters", nr.MaxLength)
	}

	// Check pattern using cached compiled regex
	if nr.Pattern != "" {
		re, err := nr.getPattern()
		if err != nil {
			return err
		}
		if re != nil && !re.MatchString(name) {
			return fmt.Errorf("branch name does not match required pattern: %s", nr.Pattern)
		}
	}

	// Check forbidden words
	for _, forbidden := range nr.Forbidden {
		pattern := fmt.Sprintf("(?i)%s", forbidden) // case-insensitive
		matched, err := regexp.MatchString(pattern, name)
		if err != nil {
			continue
		}
		if matched {
			return fmt.Errorf("branch name contains forbidden word: %s", forbidden)
		}
	}

	return nil
}

// ValidateBranchName validates a full branch name (with prefix) against Guardian rules
func (gc *GuardianConfig) ValidateBranchName(fullName string, prefix string) error {
	if !gc.Enabled {
		return nil
	}

	// Extract the name part (remove prefix)
	name := fullName
	if len(prefix) > 0 && len(fullName) > len(prefix) {
		name = fullName[len(prefix):]
	}

	return gc.Naming.Validate(name)
}
