package test

import (
	"os"
	"testing"

	"github.com/caarloshenriq/forge-cli/internal"
)

func TestPrintChangelogVersion(t *testing.T) {
	content := `
# Changelog

### Version 1.0.0 - 2025-01-01

### Features
- feat: add login

### Fixes
- fix: resolve crash on load

### Other commits
- refactor: simplify config
`

	err := os.WriteFile("CHANGELOG.md", []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create changelog: %v", err)
	}
	defer os.Remove("CHANGELOG.md")

	features, fixes, others := internal.GetChangelogVersion("1.0.0")

	if len(features) != 1 || features[0] != "feat: add login" {
		t.Errorf("expected feature not found")
	}
	if len(fixes) != 1 || fixes[0] != "fix: resolve crash on load" {
		t.Errorf("expected fix not found")
	}
	if len(others) != 1 || others[0] != "refactor: simplify config" {
		t.Errorf("expected other commit not found")
	}
}

func TestChangelogVersionExists(t *testing.T) {
	mockChangelog := `
# Changelog

### Version 1.2.0 - 2025-06-21

### Features
- feat: test feature
`

	err := os.WriteFile("CHANGELOG.md", []byte(mockChangelog), 0644)
	if err != nil {
		t.Fatalf("failed to write changelog file: %v", err)
	}
	defer os.Remove("CHANGELOG.md")

	tests := []struct {
		version   string
		expected  bool
	}{
		{"1.2.0", true},
		{"1.0.0", false},
	}

	for _, tt := range tests {
		result := internal.ChangelogVersionExists(tt.version)
		if result != tt.expected {
			t.Errorf("changelogVersionExists(%q) = %v; want %v", tt.version, result, tt.expected)
		}
	}
}
