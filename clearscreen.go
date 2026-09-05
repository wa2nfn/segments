package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ClearScreen clears the terminal screen based on the operating system.
func ClearScreen() {
	switch runtime.GOOS {
	case "windows":

		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":

		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println()
	}
}
