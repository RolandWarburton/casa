package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreatePathDir(path string, params ...bool) error {
	debug := len(params) > 0 && params[0]
	homeDir, _ := os.UserHomeDir()
	var expandedPath string

	switch path[0] {
	case '~':
		expandedPath = filepath.Join(homeDir, path[1:])
	case '/':
		expandedPath = path
	default:
		return fmt.Errorf("invalid path: %s", path)
	}

	// Check if the directory already exists and skip it
	_, err := os.Stat(expandedPath)
	if err != nil || os.IsNotExist(err) {
		if debug {
			fmt.Printf("Directory already exists: %s\n", expandedPath)
		}
		return nil
	}

	// Create the directory
	err = os.MkdirAll(expandedPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("Created directory: %s\n", expandedPath)
	return nil
}
