package utils

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// prodCurrentDir
// Get a program current directory, currently considered best practice
// With a limitation of "go run" being not work
func prodCurrentDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	realEx, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(realEx)
	if err != nil {
		return "", err
	}

	return dir, nil
}

// devCurrentDir
// Get a development current directory (Simulates prodCurrentDir but with workspace set to development folder)
// Limitations: It is currently relative to "utils" path.
func devCurrentDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return "", errors.New("Unable to get working file source directory")
	}

	dirname := filepath.Dir(filename)

	splitted := strings.Split(dirname, string(os.PathSeparator))
	dirname = strings.Join(splitted[:len(splitted)-1], string(os.PathSeparator))

	return dirname, nil
}

// GetCurrentDir
// Get an executable current directory intelligently based on application state.
// Return workspace in production / executable in development.
func GetCurrentDir() (string, error) {
	if IsProduction() {
		return prodCurrentDir()
	}

	return devCurrentDir()
}

// PublicPath
// Return an public path which smartly can guess the correct folder based on application state.
func PublicPath() string {
	currentDir, err := GetCurrentDir()

	if err != nil {
		ShouldPanic(err)
	}

	return path.Join(currentDir, "public")
}

func AssetPath(depth ...string) string {
	return path.Join(PublicPath(), "files", path.Join(depth...))
}
