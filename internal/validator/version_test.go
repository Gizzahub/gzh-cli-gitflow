package validator

import "testing"

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"valid semver", "1.0.0", false},
		{"valid with large numbers", "12.34.56", false},
		{"valid zero version", "0.0.0", false},
		{"invalid with v prefix", "v1.0.0", true},
		{"invalid two parts", "1.0", true},
		{"invalid four parts", "1.0.0.0", true},
		{"invalid with dash", "1.0.0-beta", true},
		{"invalid with plus", "1.0.0+build", true},
		{"invalid with letters", "1.0.x", true},
		{"invalid negative", "-1.0.0", true},
		{"empty", "", true},
		{"only dots", "..", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}
