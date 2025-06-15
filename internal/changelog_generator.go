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
	var cmd *exec.Cmd
	if lastCommit != "" {
		cmd = exec.Command("git", "log", fmt.Sprintf("%s..HEAD", lastCommit), "--pretty=format:%H|%s")
	} else {
		cmd = exec.Command("git", "log", "--pretty=format:%H|%s")
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(output), "\n")
	return lines, nil
}

func GenerateChangelog() {
	var version, date string
	defaultDate := time.Now().Format("2006-01-02")

	survey.AskOne(&survey.Input{
		Message: "Enter the version for this changelog (e.g., 1.0.0):",
	}, &version, survey.WithValidator(survey.Required))

	survey.AskOne(&survey.Input{
		Message: "Enter the date for this changelog:",
		Default: defaultDate,
	}, &date)

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

	var selectedFeatures, selectedFixes, selectedOthers []string

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
	var finalContent strings.Builder

	if len(existing) == 0 {
		// arquivo ainda não existe
		finalContent.WriteString(changelogHeader)
		finalContent.WriteString("\n")
		finalContent.WriteString(newEntry.String())
	} else {
		content := string(existing)
		parts := strings.SplitN(content, "### Version", 2)

		finalContent.WriteString(parts[0])
		finalContent.WriteString(newEntry.String())
		if len(parts) == 2 {
			finalContent.WriteString("\n### Version")
			finalContent.WriteString(parts[1])
		}
	}

	file, err := os.Create("CHANGELOG.md")
	if err != nil {
		fmt.Println("❌ Error creating CHANGELOG.md:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(finalContent.String())
	if err != nil {
		fmt.Println("❌ Error writing changelog:", err)
		return
	}
	writer.Flush()

	fmt.Println("✅ CHANGELOG.md generated successfully!")
}
