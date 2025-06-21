package internal

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chzyer/readline"
)

type ReadmeTemplateData struct {
	ProjectName  string
	Description  string
	Installation string
	Usage        string
	License      string
}

func readMultilineInput(prompt string) string {
	var response string
	survey.AskOne(&survey.Multiline{
		Message: prompt,
	}, &response)
	return response
}

func GenerateREADMEFromTemplate() {
	data := ReadmeTemplateData{}

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



func GenerateReadmeFromScratch() {
	var fileName string

	fmt.Println("Write your README content from scratch:")
	fmt.Println("(Use ↑ ↓ to navagate. Type ':done' on a new line to finish)")

	rl, err := readline.New("> ")
	if err != nil {
		fmt.Println("❌ Failed to initialize readline:", err)
		return
	}
	defer rl.Close()

	var lines []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == ":done" {
			break
		}
		lines = append(lines, line)
	}
	content := strings.Join(lines, "\n")

	survey.AskOne(&survey.Input{
		Message: "file name (default README):",
		Default: "README",
	}, &fileName)

	file, err := os.Create(fileName + ".md")
	if err != nil {
		fmt.Println("❌ Error creating "+fileName+".md:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("❌ Error writing to "+fileName+".md:", err)
		return
	}

	fmt.Println("✅ " + fileName + ".md created successfully from scratch!")
}