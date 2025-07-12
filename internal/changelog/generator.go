package changelog

import (
	"fmt"
	"time"

	"github.com/caarloshenriq/forge-cli/types"
	"github.com/caarloshenriq/forge-cli/utils"
)

func GenerateChangelog(fromTag string, toTag string) {
	utils.ClearScreen()
	defaultDate := time.Now().Format("2006-01-02")
	var version string
	var entry types.ChangelogSection

	if fromTag != "" {
		version = toTag
		entry.Version = version
		entry.Date = AskForDate(defaultDate)
	} else {
		for {
			version = AskForVersion()
			if ChangelogVersionExists(version) {
				if ConfirmAppend(version) {
					existingFeatures, existingFixes, existingOthers := GetChangelogVersion(version)
					entry = types.ChangelogSection{
						Version:  version,
						Date:     defaultDate,
						Features: existingFeatures,
						Fixes:    existingFixes,
						Others:   existingOthers,
						Append:   true,
					}
					break
				}
			} else {
				entry.Version = version
				entry.Date = AskForDate(defaultDate)
				break
			}
		}
	}

	var logLines []string
	var err error

	if fromTag != "" {
		logLines, err = GetGitLogBetweenTags(fromTag, toTag)
	} else {
		lastCommit := getLastCommitFromChangelog()
		logLines, err = getGitLogSince(lastCommit)
	}

	if err != nil {
		fmt.Println("❌ Error getting git log:", err)
		return
	}

	entry = classifyCommits(entry, logLines)
	entry = PromptUserSelection(entry)
	entry.Hash = GetLatestHash(logLines)

	RenderChangelog(entry)
}
