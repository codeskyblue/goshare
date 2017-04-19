package comtool

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// Return program dir
func SelfDir() string {
	return filepath.Dir(os.Args[0])
}

// Check if file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsDir(name string) bool {
	info, err := os.Stat(name)
	return err == nil && info.IsDir()
}

func IsFile(name string) bool {
	info, err := os.Stat(name)
	return err == nil && !info.IsDir()
}

// HomeDir returns path of '~'(in Linux) on Windows,
// it returns error when the variable does not exist.
func HomeDir() (home string, err error) {
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	if len(home) == 0 {
		return "", errors.New("Cannot specify home directory because it's empty")
	}

	return home, nil
}
