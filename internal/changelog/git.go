package changelog

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/caarloshenriq/forge-cli/types"
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

func classifyCommits(entry types.ChangelogSection, lines []string) types.ChangelogSection {
	featPattern := regexp.MustCompile(`^feat(\(.+\))?:`)
	fixPattern := regexp.MustCompile(`^fix(\(.+\))?:`)

	for i, line := range lines {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		hash := strings.TrimSpace(parts[0])
		msg := strings.TrimSpace(parts[1])

		if i == 0 {
			entry.LastCommitHash = hash
		}

		switch {
		case featPattern.MatchString(msg):
			entry.Features = append(entry.Features, msg)
		case fixPattern.MatchString(msg):
			entry.Fixes = append(entry.Fixes, msg)
		default:
			entry.Others = append(entry.Others, msg)
		}
	}
	return entry
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

func ParseGitLogs(lines []string) (features, fixes, others []string, latestHash string) {
	for i, line := range lines {
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
		case strings.HasPrefix(msg, "feat"):
			features = append(features, msg)
		case strings.HasPrefix(msg, "fix"):
			fixes = append(fixes, msg)
		default:
			others = append(others, msg)
		}
	}
	return
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

func GetGitLogBetweenTags(fromTag, toTag string) ([]string, error) {
	if toTag == "" {
		toTag = "HEAD"
	}

	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", fromTag, toTag), "--pretty=format:%H|%s", "--no-merges", "--date-order")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}