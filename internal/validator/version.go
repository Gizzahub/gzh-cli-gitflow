package validator

import (
	"fmt"
	"regexp"
)

var (
	// semverPattern is the strict semver pattern (X.Y.Z only)
	semverPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
)

// ValidateVersion validates a version string using strict semver format (X.Y.Z).
// Only accepts versions in the format: X.Y.Z where X, Y, Z are non-negative integers.
//
// Valid examples:
//   - 1.0.0
//   - 2.5.3
//   - 10.20.30
//
// Invalid examples:
//   - v1.0.0 (no prefix allowed)
//   - 1.0 (incomplete)
//   - 1.0.0-beta (no pre-release)
//   - 1.0.0+build (no build metadata)
func ValidateVersion(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	if !semverPattern.MatchString(version) {
		return fmt.Errorf("invalid version format (expected: X.Y.Z, e.g., 1.0.0)")
	}

	return nil
}
