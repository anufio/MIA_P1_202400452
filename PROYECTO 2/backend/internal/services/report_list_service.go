package services

//report_list_service.go: contiene la lógica para manejar la generación de reportes dentro del sistema de archivos, incluyendo la búsqueda de archivos de reporte en el directorio de reportes, la clasificación por tipo y formato, y la construcción de URLs para acceder a los reportes.

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (s *ReportService) List() ([]ReportOutput, error) {
	var result []ReportOutput

	err := filepath.WalkDir(s.reportRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".svg" && ext != ".txt" && ext != ".pdf" {
			return nil
		}

		rel, err := filepath.Rel(s.reportRoot, path)
		if err != nil {
			rel = filepath.Base(path)
		}

		result = append(result, ReportOutput{
			Type:   reportTypeFromFilename(filepath.Base(path)),
			Path:   path,
			URL:    "/reports/" + filepath.ToSlash(rel),
			Format: reportFormat(path),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(result, func(i, j int) bool {
		infoI, errI := os.Stat(result[i].Path)
		infoJ, errJ := os.Stat(result[j].Path)

		if errI == nil && errJ == nil {
			return infoI.ModTime().After(infoJ.ModTime())
		}

		return result[i].Path > result[j].Path
	})

	return result, nil
}

func reportTypeFromFilename(name string) string {
	name = strings.ToLower(name)

	known := []string{
		"bm_inode",
		"bm_block",
		"inode",
		"block",
		"tree",
		"disk",
		"mbr",
		"sb",
		"file",
		"ls",
	}

	for _, reportType := range known {
		if strings.HasPrefix(name, reportType+"_") || strings.HasPrefix(name, reportType+".") {
			return reportType
		}
	}

	return strings.TrimPrefix(filepath.Ext(name), ".")
}

func reportFormat(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg", ".png":
		return "image"
	case ".svg":
		return "svg"
	case ".txt":
		return "text"
	case ".pdf":
		return "pdf"
	default:
		return "file"
	}
}
