package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnsureProjectDirs creates:
// - <resourcesRoot>/
// - <resourcesRoot>/<projectFolder>/
// - <resourcesRoot>/<projectFolder>/media/
// - <resourcesRoot>/<projectFolder>/zip/
func EnsureProjectDirs(resourcesRoot string, projectFolder string) error {
	resourcesRoot = strings.TrimSpace(resourcesRoot)
	projectFolder = strings.TrimSpace(projectFolder)
	if resourcesRoot == "" {
		resourcesRoot = "./resources"
	}
	if projectFolder == "" {
		return errors.New("project folder name is empty")
	}
	// Basic path traversal / separator checks for relative folder names.
	// We intentionally disallow any path separators to keep it a single folder segment.
	if strings.Contains(projectFolder, "..") ||
		strings.ContainsAny(projectFolder, `/\`) ||
		projectFolder != filepath.Base(projectFolder) {
		return fmt.Errorf("invalid project folder name: %q", projectFolder)
	}

	projectRoot := filepath.Join(resourcesRoot, projectFolder)
	if err := os.MkdirAll(filepath.Join(projectRoot, "media"), 0o755); err != nil {
		return fmt.Errorf("create media dir failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(projectRoot, "zip"), 0o755); err != nil {
		return fmt.Errorf("create zip dir failed: %w", err)
	}
	return nil
}
