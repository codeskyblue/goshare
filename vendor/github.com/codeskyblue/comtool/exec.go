package comtool

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func ReadExitcode(err error) int {
	if err == nil {
		return 0
	}
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 126
}

func ShRun(name string, args ...string) error {
	if runtime.GOOS == "windows" {
		args = append([]string{"/c", name}, args...)
		name = "cmd"
	}
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}
