package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
)

type ReadmeTemplateData struct {
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

func readFreeformContent(prompt string) string {
	fmt.Println(prompt)
	fmt.Println("(Type ':done' on a new line to finish)")

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

func GenerateReadmeFromScratch() {
	var fileName string = "README"
	content := readFreeformContent("Write your README content from scratch:")
	survey.AskOne(&survey.Input{
		Message: "file name (default README):",
		Default: "README",
		}, &fileName)
	
	file, err := os.Create(fileName + ".md")
	if err != nil {
		fmt.Println("Error creating "+fileName+".md:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to "+fileName+".md:", err)
		return
	}

	fmt.Println(fileName+".md created successfully from scratch!")
}
