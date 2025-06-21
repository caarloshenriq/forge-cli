package test

import (
	"strings"
	"testing"
	"text/template"

	"github.com/caarloshenriq/forge-cli/internal"
)

func TestReadmeTemplateGeneration(t *testing.T) {
	data := internal.ReadmeTemplateData{
		ProjectName:  "Forge CLI",
		Description:  "CLI to generate boilerplate documentation",
		Installation: "go install github.com/your/repo",
		Usage:        "forge-cli generate",
		License:      "MIT",
	}

	const expectedContent = `# Forge CLI

CLI to generate boilerplate documentation

## Installation

go install github.com/your/repo

## Usage

forge-cli generate

## License

MIT
	`

	const readmeTemplate = `# {{.ProjectName}}

{{.Description}}

## Installation

{{.Installation}}

## Usage

{{.Usage}}

## License

{{.License}}`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		t.Fatalf("Error parsing template: %v", err)
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, data)
	if err != nil {
		t.Fatalf("Error executing template: %v", err)
	}

	result := builder.String()

	if strings.TrimSpace(result) != strings.TrimSpace(expectedContent) {
		t.Errorf("Generated content doesn't match expected.\nExpected:\n%s\n\nGot:\n%s", expectedContent, result)
	}
}
