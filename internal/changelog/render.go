package changelog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/caarloshenriq/forge-cli/types"

)

const changelogHeader = `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
`

func RenderChangelog(entry types.ChangelogSection) {
	var newEntry strings.Builder
	newEntry.WriteString(fmt.Sprintf("### Version %s - %s\n\n", entry.Version, entry.Date))

	if len(entry.Features) > 0 {
		newEntry.WriteString("### Features\n")
		for _, feat := range entry.Features {
			newEntry.WriteString("- " + feat + "\n")
		}
		newEntry.WriteString("\n")
	}

	if len(entry.Fixes) > 0 {
		newEntry.WriteString("### Fixes\n")
		for _, fix := range entry.Fixes {
			newEntry.WriteString("- " + fix + "\n")
		}
		newEntry.WriteString("\n")
	}

	if len(entry.Others) > 0 {
		newEntry.WriteString("### Other commits\n")
		for _, other := range entry.Others {
			newEntry.WriteString("- " + other + "\n")
		}
		newEntry.WriteString("\n")
	}

	if entry.Hash != "" {
		newEntry.WriteString(fmt.Sprintf("<!-- last-commit: %s -->\n", entry.Hash))
	}

	existing, _ := os.ReadFile("CHANGELOG.md")
	content := string(existing)
	var finalContent string

	if len(existing) == 0 {
		finalContent = changelogHeader + "\n" + newEntry.String()
	} else if entry.Append {
		lines := strings.Split(content, "\n")
		var builder strings.Builder
		header := fmt.Sprintf("### Version %s ", entry.Version)
		skipping := false

		for i := 0; i < len(lines); i++ {
			line := lines[i]

			if strings.HasPrefix(line, "### Version ") {
				if strings.HasPrefix(line, header) {
					skipping = true
					builder.WriteString(newEntry.String())
					builder.WriteString("\n")
					for i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "### Version ") {
						i++
					}
					continue
				}
			}

			if !skipping {
				builder.WriteString(line + "\n")
			} else if strings.HasPrefix(line, "### Version ") {
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
