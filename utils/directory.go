package utils

import (
	"os"
	"path"
)

func GetCurrentDir() (string, error) {
	current, err := os.Executable()

	if err != nil {
		return "", err
	}

	return path.Dir(current), nil
}
