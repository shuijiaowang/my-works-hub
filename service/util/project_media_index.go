package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const MediaIndexFileName = "_media_index.json"

type MediaKind string

const (
	MediaKindImage MediaKind = "image"
	MediaKindVideo MediaKind = "video"
)

type MediaItem struct {
	ID           string    `json:"id"`
	FileName     string    `json:"fileName"`
	OriginalName string    `json:"originalName"`
	Kind         MediaKind `json:"kind"`
	Order        int       `json:"order"`
	CreatedAt    time.Time `json:"createdAt"`
}

type mediaIndex struct {
	Items []MediaItem `json:"items"`
}

var mediaIndexMu sync.Map // map[absIndexPath]*sync.Mutex

func mediaIndexLock(indexPath string) *sync.Mutex {
	indexPath = strings.TrimSpace(indexPath)
	if indexPath == "" {
		// should never happen
		return &sync.Mutex{}
	}
	v, _ := mediaIndexMu.LoadOrStore(indexPath, &sync.Mutex{})
	return v.(*sync.Mutex)
}

func MediaIndexPath(mediaDir string) string {
	return filepath.Join(mediaDir, MediaIndexFileName)
}

func LoadMediaIndex(mediaDir string) ([]MediaItem, error) {
	mediaDir = strings.TrimSpace(mediaDir)
	if mediaDir == "" {
		return nil, errors.New("mediaDir is empty")
	}
	indexPath := MediaIndexPath(mediaDir)

	mu := mediaIndexLock(indexPath)
	mu.Lock()
	defer mu.Unlock()

	return loadMediaIndexNoLock(mediaDir)
}

func SaveMediaIndex(mediaDir string, items []MediaItem) error {
	mediaDir = strings.TrimSpace(mediaDir)
	if mediaDir == "" {
		return errors.New("mediaDir is empty")
	}
	indexPath := MediaIndexPath(mediaDir)

	mu := mediaIndexLock(indexPath)
	mu.Lock()
	defer mu.Unlock()

	idx := mediaIndex{Items: items}
	return saveMediaIndexNoLock(indexPath, idx)
}

func loadMediaIndexNoLock(mediaDir string) ([]MediaItem, error) {
	indexPath := MediaIndexPath(mediaDir)
	b, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []MediaItem{}, nil
		}
		return nil, fmt.Errorf("read media index failed: %w", err)
	}

	var idx mediaIndex
	if err := json.Unmarshal(b, &idx); err != nil {
		return nil, fmt.Errorf("parse media index failed: %w", err)
	}
	if idx.Items == nil {
		return []MediaItem{}, nil
	}
	return idx.Items, nil
}

func saveMediaIndexNoLock(indexPath string, idx mediaIndex) error {
	if err := os.MkdirAll(filepath.Dir(indexPath), 0o755); err != nil {
		return fmt.Errorf("ensure media dir failed: %w", err)
	}

	b, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal media index failed: %w", err)
	}

	tmp := indexPath + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return fmt.Errorf("write media index tmp failed: %w", err)
	}
	if err := os.Rename(tmp, indexPath); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename media index tmp failed: %w", err)
	}
	return nil
}

func DetectMediaKindByName(name string) MediaKind {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(name)))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".bmp":
		return MediaKindImage
	case ".mp4", ".webm", ".ogg", ".mov", ".m4v":
		return MediaKindVideo
	default:
		// 默认当图片处理（前端用 <img> 失败再回退也行），但这里更保守：未知当 image
		return MediaKindImage
	}
}

func SanitizeFileBaseName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	name = filepath.Base(name)
	// allow: letters, digits, dot, dash, underscore; replace others with dash
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '.' || r == '-' || r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	out := strings.Trim(b.String(), "-._")
	return out
}

func NewMediaFileName(original string) (id string, fileName string) {
	orig := strings.TrimSpace(original)
	ext := strings.ToLower(filepath.Ext(orig))
	if ext == "" || len(ext) > 10 {
		ext = ""
	}
	id = fmt.Sprintf("%x%x", time.Now().UnixNano(), rand.Uint64())
	return id, id + ext
}

func NormalizeAndSortMedia(items []MediaItem) []MediaItem {
	out := make([]MediaItem, 0, len(items))
	seen := map[string]bool{}
	for _, it := range items {
		if strings.TrimSpace(it.ID) == "" {
			continue
		}
		if seen[it.ID] {
			continue
		}
		seen[it.ID] = true
		out = append(out, it)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Order == out[j].Order {
			return out[i].CreatedAt.Before(out[j].CreatedAt)
		}
		return out[i].Order < out[j].Order
	})
	// re-assign contiguous orders
	for i := range out {
		out[i].Order = i
	}
	return out
}
