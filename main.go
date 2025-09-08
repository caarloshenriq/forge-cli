package main

import (
	"fmt"
	"os"

	"github.com/caarloshenriq/forge-cli/cmd"
	"github.com/caarloshenriq/forge-cli/internal"
)

func main() {
	for {
		selected := internal.RunTUI()

		switch selected {
		case "Generate Changelog":
			if cmd.ChangelogCmd.RunE != nil {
				if err := cmd.ChangelogCmd.RunE(cmd.ChangelogCmd, nil); err != nil {
					fmt.Fprintln(os.Stderr, "changelog error:", err)
				}
			} else if cmd.ChangelogCmd.Run != nil {
				cmd.ChangelogCmd.Run(cmd.ChangelogCmd, nil)
			}
			return

		case "Generate README":
			if cmd.ReadmeCmd.RunE != nil {
				if err := cmd.ReadmeCmd.RunE(cmd.ReadmeCmd, nil); err != nil {
					fmt.Fprintln(os.Stderr, "readme error:", err)
				}
			} else if cmd.ReadmeCmd.Run != nil {
				cmd.ReadmeCmd.Run(cmd.ReadmeCmd, nil)
			} else {
				if err := internal.RunReadmeTUI(); err != nil {
					fmt.Fprintln(os.Stderr, "readme error:", err)
				}
			}
			return

		case "Help":
			cmd.RunHelpTUI()

		default:
			return
		}
	}
}
