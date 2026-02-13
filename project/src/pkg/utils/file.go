// Package utils
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var PathSeparator string = fmt.Sprintf("%c", os.PathSeparator)

// GetExecutableDir returns the directory where the executable is located
func GetExecutableDir() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve any symlinks
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	return filepath.Dir(executable), nil
}

// GetProjectRoot returns the project root directory
// In development, it finds the root by looking for go.mod
// In production, it uses the executable directory
func GetProjectRoot() (string, error) {
	// First, try to find go.mod (development mode)
	if root, err := findProjectRootByGoMod(); err == nil {
		return root, nil
	}

	// If go.mod not found, assume we're in production and use executable directory
	return GetExecutableDir()
}

// findProjectRootByGoMod finds the project root by looking for go.mod file
func findProjectRootByGoMod() (string, error) {
	// Start from the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// If that doesn't work, try from the caller's directory
	if !fileExists(filepath.Join(wd, "go.mod")) {
		_, filename, _, ok := runtime.Caller(1)
		if ok {
			wd = filepath.Dir(filename)
		}
	}

	// Walk up the directory tree looking for go.mod
	for {
		if fileExists(filepath.Join(wd, "go.mod")) {
			return wd, nil
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			// Reached the root directory
			break
		}
		wd = parent
	}

	return "", fmt.Errorf("go.mod not found")
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	return fileExists(filename)
}

// fileExists is a helper function to check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadFile reads the entire file and returns its contents
func ReadFile(filename string) ([]byte, error) {
	if !FileExists(filename) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	return os.ReadFile(filename)
}

// GetConfigPath returns the path to the config directory based on the current environment
func GetConfigPath(service string) (string, error) {
	projectRoot, err := GetProjectRoot()
	if err != nil {
		return "", fmt.Errorf("failed to get project root: %w", err)
	}

	// Check if we're in development mode (go.mod exists in project root)
	if fileExists(filepath.Join(projectRoot, "go.mod")) {
		// Development mode: use src/services/main/configs
		configPath := filepath.Join(projectRoot, "src", "services", service, "configs")
		if !dirExists(configPath) {
			return "", fmt.Errorf("config directory not found: %s", configPath)
		}
		return configPath, nil
	}

	// Production mode: look for configs directory relative to executable
	// Try different possible locations
	possiblePaths := []string{
		filepath.Join(projectRoot, "configs"),                            // Same level as executable
		filepath.Join(projectRoot, "src", "services", "main", "configs"), // Full path structure
	}

	for _, path := range possiblePaths {
		if dirExists(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("config directory not found in any of the expected locations: %s", strings.Join(possiblePaths, ", "))
}

// dirExists checks if a directory exists
func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func GetFilePath(pathFromRoot []string) (string, error) {
	root, err := GetProjectRoot()
	if err != nil {
		return "", err
	}
	root += PathSeparator
	for i, x := range pathFromRoot {
		if i < len(pathFromRoot)-1 {
			root += x + PathSeparator
		} else {
			root += x
		}
	}
	return root, nil
}

func GetDefaultLogsFileName() (string, error) {
	path, err := GetFilePath([]string{"logs", "server.log"})
	if err != nil {
		return "", err
	}
	return path, nil
}
