package validator

import (
	"testing"
)

func TestValidateBranchName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid kebab-case", "feature-branch", false},
		{"valid kebab", "user-auth", false},
		{"valid with numbers", "feature-123", false},
		{"valid nested", "feature/my-feature", false},
		{"uppercase", "UserAuth", true},
		{"empty string", "", true},
		{"too long", "this-is-a-very-long-branch-name-that-exceeds-the-maximum-allowed-length-limit", true},
		{"contains spaces", "feature branch", true},
		{"contains special chars", "feature@branch", true},
		{"double slash", "feature//branch", true},
		{"ends with slash", "feature/", true},
		{"starts with slash", "/feature", true},
		{"contains backslash", "feature\\branch", true},
		{"reserved name master", "master", true},
		{"reserved name main", "main", true},
		{"reserved name develop", "develop", true},
		{"reserved name HEAD", "HEAD", true},
		// Security pattern tests
		{"semicolon", "feature;rm", true},
		{"pipe", "feature|cat", true},
		{"ampersand", "feature&echo", true},
		{"greater than", "feature>file", true},
		{"less than", "feature<file", true},
		{"asterisk", "feature*", true},
		{"question mark", "feature?", true},
		{"bracket open", "feature[", true},
		{"bracket close", "feature]", true},
		{"tilde", "feature~1", true},
		{"caret", "feature^1", true},
		{"colon", "feature:branch", true},
		{"tab", "feature\tbranch", true},
		{"newline", "feature\nbranch", true},
		{"carriage return", "feature\rbranch", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBranchName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBranchName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestSuggestBranchName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"spaces to hyphens", "my feature branch", "my-feature-branch"},
		{"uppercase to lowercase", "MyFeature", "myfeature"},
		{"underscores to hyphens", "my_feature", "my-feature"},
		{"remove special chars", "feature@#$branch", "featurebranch"},
		{"collapse multiple hyphens", "feature---branch", "feature-branch"},
		{"trim hyphens", "-feature-branch-", "feature-branch"},
		{"truncate long", "this-is-a-very-long-branch-name-that-exceeds-the-maximum-limit", "this-is-a-very-long-branch-name-that-exceeds-the-m"},
		{"empty result", "@#$", ""},
		{"reserved name", "MASTER", "master-branch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SuggestBranchName(tt.input)
			if got != tt.want {
				t.Errorf("SuggestBranchName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
