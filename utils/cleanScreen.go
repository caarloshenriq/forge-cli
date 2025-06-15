package utils

import (
	"os"
	"os/exec"
)


func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}