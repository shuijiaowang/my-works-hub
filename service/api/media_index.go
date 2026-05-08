package api

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"workhub/util"
)

// ensureAndLoadMediaIndex loads index.json; if missing, it bootstraps it from current files.
// It also reconciles missing files and appends newly found files to the end.
func ensureAndLoadMediaIndex(mediaDir string) ([]util.MediaItem, error) {
	mediaDir = strings.TrimSpace(mediaDir)
	if mediaDir == "" {
		return nil, fmt.Errorf("mediaDir is empty")
	}

	// Best-effort ensure dir exists (project might exist but dir absent)
	_ = os.MkdirAll(mediaDir, 0o755)

	items, err := util.LoadMediaIndex(mediaDir)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(mediaDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []util.MediaItem{}, nil
		}
		return nil, fmt.Errorf("read media dir failed: %w", err)
	}

	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if e.Name() == util.MediaIndexFileName {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	byFile := map[string]util.MediaItem{}
	for _, it := range items {
		if strings.TrimSpace(it.FileName) == "" {
			continue
		}
		byFile[it.FileName] = it
	}

	// Remove items whose files are gone; append new files.
	now := time.Now()
	nextOrder := len(items)

	exists := map[string]bool{}
	for _, f := range files {
		exists[f] = true
		if _, ok := byFile[f]; ok {
			continue
		}
		id := strings.TrimSuffix(f, filepath.Ext(f))
		id = util.SanitizeFileBaseName(id)
		if id == "" {
			id = f
		}
		items = append(items, util.MediaItem{
			ID:           id,
			FileName:     f,
			OriginalName: f,
			Kind:         util.DetectMediaKindByName(f),
			Order:        nextOrder,
			CreatedAt:    now,
		})
		nextOrder++
	}

	filtered := make([]util.MediaItem, 0, len(items))
	for _, it := range items {
		if strings.TrimSpace(it.FileName) == "" {
			continue
		}
		if !exists[it.FileName] {
			continue
		}
		filtered = append(filtered, it)
	}

	normalized := util.NormalizeAndSortMedia(filtered)

	// Persist reconciliation results (ignore if it fails? better return err)
	if err := util.SaveMediaIndex(mediaDir, normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}
