package storage

// storage.go

import (
	"os"
	"path/filepath"
	"strings"
)

const BaseDir = "sesiones"

func ResolvePath(rawPath string) string {
	if rawPath == "" {
		return ""
	}
	if strings.HasPrefix(rawPath, "/") {
		dir := filepath.Dir(rawPath)
		os.MkdirAll(dir, 0755)
		return rawPath
	}
	fullPath := filepath.Join(BaseDir, rawPath)
	os.MkdirAll(BaseDir, 0755)
	dir := filepath.Dir(fullPath)
	os.MkdirAll(dir, 0755)
	return fullPath
}
