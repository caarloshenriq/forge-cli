package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func GenerateChangelog() {
	var version string
	var date string

	fmt.Print("Enter the version for this changelog (e.g., 1.0.0): ")
	fmt.Scanln(&version)

	fmt.Print("Enter the date for this changelog (e.g., 2023-10-01): ")
	fmt.Scanln(&date)

	cmd := exec.Command("git", "log", "--pretty=format:%s")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("❌ Error running git log:", err)
		return
	}

	lines := strings.Split(string(output), "\n")

	var (
		rawFeatures []string
		rawFixes    []string
		rawOthers   []string

		featPattern = regexp.MustCompile(`^feat(\(.+\))?:`)
		fixPattern  = regexp.MustCompile(`^fix(\(.+\))?:`)
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case featPattern.MatchString(line):
			rawFeatures = append(rawFeatures, line)
		case fixPattern.MatchString(line):
			rawFixes = append(rawFixes, line)
		default:
			rawOthers = append(rawOthers, line)
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

	var changelog strings.Builder
	changelog.WriteString("# Changelog\n\n")
	changelog.WriteString("All notable changes to this project will be documented in this file.\n\n")
	changelog.WriteString("The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),\n")
	changelog.WriteString("and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n\n")

	changelog.WriteString("### Version " + version + " - " + date + "\n\n")

	if len(selectedFeatures) > 0 {
		changelog.WriteString("### Features\n")
		for _, feat := range selectedFeatures {
			changelog.WriteString("- " + feat + "\n")
		}
		changelog.WriteString("\n")
	}

	if len(selectedFixes) > 0 {
		changelog.WriteString("### Fixes\n")
		for _, fix := range selectedFixes {
			changelog.WriteString("- " + fix + "\n")
		}
		changelog.WriteString("\n")
	}

	if len(selectedOthers) > 0 {
		changelog.WriteString("### Other commits\n")
		for _, other := range selectedOthers {
			changelog.WriteString("- " + other + "\n")
		}
		changelog.WriteString("\n")
	}

	file, err := os.Create("CHANGELOG.md")
	if err != nil {
		fmt.Println("❌ Error creating CHANGELOG.md:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(changelog.String())
	if err != nil {
		fmt.Println("❌ Error writing changelog:", err)
		return
	}
	writer.Flush()

	fmt.Println("✅ CHANGELOG.md generated successfully!")
}
