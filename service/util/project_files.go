package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ProjectMediaDir(resourcesRoot, projectFolder string) string {
	resourcesRoot = strings.TrimSpace(resourcesRoot)
	if resourcesRoot == "" {
		resourcesRoot = "./resources"
	}
	return filepath.Join(resourcesRoot, projectFolder, "media")
}

func ProjectZipDir(resourcesRoot, projectFolder string) string {
	resourcesRoot = strings.TrimSpace(resourcesRoot)
	if resourcesRoot == "" {
		resourcesRoot = "./resources"
	}
	return filepath.Join(resourcesRoot, projectFolder, "zip")
}

// ListFiles lists files (not dirs) directly under dir.
// If dir doesn't exist, it returns an empty list and nil error.
func ListFiles(dir string) ([]string, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return nil, fmt.Errorf("dir is empty")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("read dir %q failed: %w", dir, err)
	}

	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		// Hide internal index file from generic listing
		if e.Name() == MediaIndexFileName {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out, nil
}
