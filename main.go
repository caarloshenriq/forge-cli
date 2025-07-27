package main

import (
	"github.com/caarloshenriq/forge-cli/cmd"
	"github.com/caarloshenriq/forge-cli/internal"
)

func main() {
	for {
		selected := internal.RunTUI()

		switch selected {
		case "Generate Changelog":
			// Lógica real aqui
			return
		case "Generate README":
			// Lógica real aqui
			return
		case "Help":
			cmd.RunHelpTUI()
		default:
			return
		}
	}
}

