package utils

import (
	"fmt"
	"os"
	"strings"
)

func CheckCancel(input string) {
	if strings.TrimSpace(input) == ":cancel" {
		fmt.Println("❌ Operation cancelled by user.")
		os.Exit(0)
	}
}