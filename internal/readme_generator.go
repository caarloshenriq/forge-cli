package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
)

type ReadmeData struct {
	ProjectName  string
	Description  string
	Installation string
	Usage        string
	License      string
}

func readMultilineInput(prompt string) string {
	fmt.Println(prompt + " (type ':done' to finish):")

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == ":done" {
			break
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func GenerateREADME() {
	data := ReadmeData{}

	survey.AskOne(&survey.Input{
		Message: "What is the project name?",
	}, &data.ProjectName, survey.WithValidator(survey.Required))

	data.Description = readMultilineInput("Enter project description")
	data.Installation = readMultilineInput("Enter installation instructions")
	data.Usage = readMultilineInput("Enter usage examples")

	survey.AskOne(&survey.Input{
		Message: "License type (default MIT):",
	}, &data.License)
	if data.License == "" {
		data.License = "MIT"
	}

	const readmeTemplate = `# {{.ProjectName}}

{{.Description}}

## Installation

{{.Installation}}

## Usage

{{.Usage}}

## License

{{.License}}
`

	file, err := os.Create("README.md")
	if err != nil {
		fmt.Println("Error creating README.md:", err)
		return
	}
	defer file.Close()

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Error writing to README.md:", err)
		return
	}

	fmt.Println("✅ README.md generated successfully!")
}
