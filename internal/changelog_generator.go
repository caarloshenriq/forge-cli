package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarloshenriq/forge-cli/utils"
)

func getLastCommitFromChangelog() string {
	utils.ClearScreen()
	file, err := os.Open("CHANGELOG.md")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "last-commit:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
		if strings.HasPrefix(line, "### Version") {
			break
		}
	}
	return ""
}

func getGitLogSince(lastCommit string) ([]string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%H|%s", "--no-merges", "--date-order")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	if lastCommit == "" {
		return lines, nil
	}

	var filtered []string
	found := false
	for _, line := range lines {
		if strings.HasPrefix(line, lastCommit) {
			found = true
			continue
		}
		if !found {
			filtered = append(filtered, line)
		}
	}

	return filtered, nil
}

func ChangelogVersionExists(version string) bool {
	file, err := os.Open("CHANGELOG.md")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	header := fmt.Sprintf("### Version %s ", version)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), header) {
			return true
		}
	}
	return false
}

func GetChangelogVersion(version string) (features, fixes, others []string) {
	file, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return
	}

	lines := strings.Split(string(file), "\n")
	header := fmt.Sprintf("### Version %s ", version)

	var (
		found   bool
		section []string
	)

	for _, line := range lines {
		if strings.HasPrefix(line, "### Version ") {
			if found {
				break
			}
			if strings.HasPrefix(line, header) {
				found = true
			}
		}
		if found {
			section = append(section, line)
		}
	}

	if len(section) == 0 {
		return
	}

	var current string
	for _, line := range section {
		switch {
		case strings.HasPrefix(line, "### Features"):
			current = "features"
		case strings.HasPrefix(line, "### Fixes"):
			current = "fixes"
		case strings.HasPrefix(line, "### Other commits"):
			current = "others"
		case strings.HasPrefix(line, "- "):
			item := strings.TrimPrefix(line, "- ")
			switch current {
			case "features":
				features = append(features, item)
			case "fixes":
				fixes = append(fixes, item)
			case "others":
				others = append(others, item)
			}
		}
	}

	return
}

func GenerateChangelog() {
	var version, date string
	defaultDate := time.Now().Format("2006-01-02")
	var selectedFeatures, selectedFixes, selectedOthers []string
	var isAppending bool

	for {
		survey.AskOne(&survey.Input{
			Message: "Enter the version for this changelog (e.g., 1.0.0):",
		}, &version, survey.WithValidator(survey.Required))

		if ChangelogVersionExists(version) {
			var choice string
			survey.AskOne(&survey.Select{
				Message: fmt.Sprintf("⚠️ Version %s already exists. What do you want to do?", version),
				Options: []string{
					"Add new commits to this version",
					"Cancel and input a new version",
				},
			}, &choice)

			if choice == "Add new commits to this version" {
				existingFeatures, existingFixes, existingOthers := GetChangelogVersion(version)
				selectedFeatures = append(selectedFeatures, existingFeatures...)
				selectedFixes = append(selectedFixes, existingFixes...)
				selectedOthers = append(selectedOthers, existingOthers...)
				isAppending = true
				break
			}
		} else {
			break
		}
	}

	if !isAppending {
		survey.AskOne(&survey.Input{
			Message: "Enter the date for this changelog:",
			Default: defaultDate,
		}, &date)
	} else {
		date = defaultDate
	}

	lastCommit := getLastCommitFromChangelog()
	logLines, err := getGitLogSince(lastCommit)
	if err != nil {
		fmt.Println("❌ Error getting git log:", err)
		return
	}

	var (
		rawFeatures []string
		rawFixes    []string
		rawOthers   []string
		latestHash  string

		featPattern = regexp.MustCompile(`^feat(\(.+\))?:`)
		fixPattern  = regexp.MustCompile(`^fix(\(.+\))?:`)
	)

	for i, line := range logLines {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		hash := strings.TrimSpace(parts[0])
		msg := strings.TrimSpace(parts[1])

		if i == 0 {
			latestHash = hash
		}

		switch {
		case featPattern.MatchString(msg):
			rawFeatures = append(rawFeatures, msg)
		case fixPattern.MatchString(msg):
			rawFixes = append(rawFixes, msg)
		default:
			rawOthers = append(rawOthers, msg)
		}
	}

	if len(rawFeatures) > 0 {
		survey.AskOne(&survey.MultiSelect{
			Message: "Select features to include:",
			Options: rawFeatures,
			Default: rawFeatures,
		}, &selectedFeatures)
	}

	if len(rawFixes) > 0 {
		survey.AskOne(&survey.MultiSelect{
			Message: "Select fixes to include:",
			Options: rawFixes,
			Default: rawFixes,
		}, &selectedFixes)
	}

	if len(rawOthers) > 0 {
		survey.AskOne(&survey.MultiSelect{
			Message: "Select other commits to include:",
			Options: rawOthers,
			Default: rawOthers,
		}, &selectedOthers)
	}

	var newEntry strings.Builder
	newEntry.WriteString(fmt.Sprintf("### Version %s - %s\n\n", version, date))

	if len(selectedFeatures) > 0 {
		newEntry.WriteString("### Features\n")
		for _, feat := range selectedFeatures {
			newEntry.WriteString("- " + feat + "\n")
		}
		newEntry.WriteString("\n")
	}

	if len(selectedFixes) > 0 {
		newEntry.WriteString("### Fixes\n")
		for _, fix := range selectedFixes {
			newEntry.WriteString("- " + fix + "\n")
		}
		newEntry.WriteString("\n")
	}

	if len(selectedOthers) > 0 {
		newEntry.WriteString("### Other commits\n")
		for _, other := range selectedOthers {
			newEntry.WriteString("- " + other + "\n")
		}
		newEntry.WriteString("\n")
	}

	if latestHash != "" {
		newEntry.WriteString(fmt.Sprintf("<!-- last-commit: %s -->\n", latestHash))
	}

	const changelogHeader = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
`

	existing, _ := os.ReadFile("CHANGELOG.md")
	content := string(existing)
	var finalContent string

	if len(existing) == 0 {
		finalContent = changelogHeader + "\n" + newEntry.String()
	} else if isAppending {
		lines := strings.Split(content, "\n")
		var builder strings.Builder
		header := fmt.Sprintf("### Version %s ", version)
		skipping := false
	
		for i := 0; i < len(lines); i++ {
			line := lines[i]
	
			if strings.HasPrefix(line, "### Version ") {
				if strings.HasPrefix(line, header) {
					// Começamos a sobrescrever a versão existente
					skipping = true
					builder.WriteString(newEntry.String())
					builder.WriteString("\n")
					// Pular linhas até o próximo bloco de versão ou EOF
					for i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "### Version ") {
						i++
					}
					continue
				}
			}
	
			// Só escreve se não estiver pulando
			if !skipping {
				builder.WriteString(line + "\n")
			} else if strings.HasPrefix(line, "### Version ") {
				// Terminamos de pular, voltar a escrever normalmente
				skipping = false
				builder.WriteString(line + "\n")
			}
		}
	
		finalContent = builder.String()
	} else {
		parts := strings.SplitN(content, "### Version", 2)
		finalContent = parts[0] + newEntry.String()
		if len(parts) == 2 {
			finalContent += "\n### Version" + parts[1]
		}
	}

	file, err := os.Create("CHANGELOG.md")
	if err != nil {
		fmt.Println("❌ Error creating CHANGELOG.md:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(finalContent)
	if err != nil {
		fmt.Println("❌ Error writing changelog:", err)
		return
	}
	writer.Flush()

	fmt.Println("✅ CHANGELOG.md generated successfully!")
}
