package cmd

import (
	"os/exec"
)

func BuildCmd(path string) *exec.Cmd {
	return exec.Command("cmd.exe", "/c", path)
}
