package changelog

import (
	"github.com/AlecAivazis/survey/v2"
)

func AskForVersion() string {
	var version string
	survey.AskOne(&survey.Input{
		Message: "Enter the version for this changelog (e.g., 1.0.0):",
	}, &version, survey.WithValidator(survey.Required))
	return version
}

func AskForDate(defaultDate string) string {
	var date string
	survey.AskOne(&survey.Input{
		Message: "Enter the date for this changelog:",
		Default: defaultDate,
	}, &date)
	return date
}

func ConfirmAppend(version string) bool {
	var choice string
	survey.AskOne(&survey.Select{
		Message: "⚠️ Version already exists. What do you want to do?",
		Options: []string{"Add new commits to this version", "Cancel and input a new version"},
	}, &choice)
	return choice == "Add new commits to this version"
}
