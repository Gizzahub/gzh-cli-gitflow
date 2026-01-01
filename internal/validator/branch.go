package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// DefaultPattern is the default regex pattern for valid branch names (kebab-case)
	DefaultPattern = `^[a-z0-9]+(-[a-z0-9]+)*(/[a-z0-9]+(-[a-z0-9]+)*)*$`

	// MaxBranchLength is the maximum allowed length for branch names
	MaxBranchLength = 50
)

var (
	// ReservedNames are branch names that cannot be used
	ReservedNames = []string{"master", "main", "develop", "HEAD", "FETCH_HEAD", "ORIG_HEAD"}

	// ForbiddenPatterns are patterns that indicate security issues or injection attempts
	// Spec requires: "..", "/", "\\", " ", "@", "$", "`"
	// Extended for additional security coverage
	ForbiddenPatterns = []string{
		"..", // directory traversal
		"//", // double slash
		"\\", // backslash
		"@{", // git ref syntax
		"$",  // shell variable
		"`",  // command substitution
		";",  // command separator
		"|",  // pipe
		"&",  // background process
		">",  // redirection
		"<",  // redirection
		"*",  // wildcard
		"?",  // wildcard
		"[",  // character class
		"]",  // character class
		"~",  // home directory
		"^",  // git ref syntax
		":",  // git ref syntax
		" ",  // space
		"\t", // tab
		"\n", // newline
		"\r", // carriage return
	}

	defaultRegex        = regexp.MustCompile(DefaultPattern)
	multiHyphenPattern  = regexp.MustCompile(`-+`)
)

// ValidateBranchName validates a branch name according to git-flow conventions
// and security requirements.
func ValidateBranchName(name string) error {
	if name == "" {
		return errors.New("branch name cannot be empty")
	}

	if len(name) > MaxBranchLength {
		return fmt.Errorf("branch name exceeds maximum length of %d characters", MaxBranchLength)
	}

	// Check reserved names
	for _, reserved := range ReservedNames {
		if name == reserved {
			return errors.New("branch name '" + name + "' is reserved and cannot be used")
		}
	}

	// Check forbidden patterns for security
	for _, forbidden := range ForbiddenPatterns {
		if strings.Contains(name, forbidden) {
			return errors.New("branch name contains forbidden character or pattern: '" + forbidden + "'")
		}
	}

	// Check if starts or ends with slash
	if strings.HasPrefix(name, "/") {
		return errors.New("branch name cannot start with '/'")
	}
	if strings.HasSuffix(name, "/") {
		return errors.New("branch name cannot end with '/'")
	}

	// Validate against pattern
	if !defaultRegex.MatchString(name) {
		return errors.New("branch name must be in kebab-case format (lowercase letters, numbers, and hyphens)")
	}

	return nil
}

// SuggestBranchName converts an invalid branch name into a valid one by:
// - Converting to lowercase
// - Replacing underscores and spaces with hyphens
// - Removing special characters
// - Collapsing multiple hyphens
// - Trimming leading/trailing hyphens
// - Truncating if too long
func SuggestBranchName(input string) string {
	// Convert to lowercase
	result := strings.ToLower(input)

	// Replace common separators with hyphens
	result = strings.ReplaceAll(result, "_", "-")
	result = strings.ReplaceAll(result, " ", "-")

	// Remove all characters that are not alphanumeric, hyphen, or slash
	var cleaned strings.Builder
	for _, char := range result {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' || char == '/' {
			cleaned.WriteRune(char)
		}
	}
	result = cleaned.String()

	// Collapse multiple hyphens into single hyphen
	result = multiHyphenPattern.ReplaceAllString(result, "-")

	// Trim leading and trailing hyphens
	result = strings.Trim(result, "-")

	// Truncate if too long
	if len(result) > MaxBranchLength {
		result = result[:MaxBranchLength]
		result = strings.TrimRight(result, "-")
	}

	// Check if result is a reserved name
	resultLower := strings.ToLower(result)
	for _, reserved := range ReservedNames {
		if resultLower == strings.ToLower(reserved) {
			return result + "-branch"
		}
	}

	return result
}
