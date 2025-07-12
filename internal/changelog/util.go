package changelog

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarloshenriq/forge-cli/types"
)

func GetLatestHash(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	parts := strings.SplitN(lines[0], "|", 2)
	if len(parts) < 1 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}

func PromptUserSelection(entry types.ChangelogSection) types.ChangelogSection {
	var selected []string

	if len(entry.Features) > 0 {
		survey.AskOne(&survey.MultiSelect{
			Message: "Select features to include:",
			Options: entry.Features,
			Default: entry.Features,
		}, &selected)
		entry.Features = selected
	}

	if len(entry.Fixes) > 0 {
		selected = nil
		survey.AskOne(&survey.MultiSelect{
			Message: "Select fixes to include:",
			Options: entry.Fixes,
			Default: entry.Fixes,
		}, &selected)
		entry.Fixes = selected
	}

	if len(entry.Others) > 0 {
		selected = nil
		survey.AskOne(&survey.MultiSelect{
			Message: "Select other commits to include:",
			Options: entry.Others,
			Default: entry.Others,
		}, &selected)
		entry.Others = selected
	}

	return entry
}
