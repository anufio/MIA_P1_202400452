package reports

//rep.go: contiene la lógica para generar reportes del sistema, incluyendo reportes de disco, bitmap, árbol de directorios y journaling.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"MIA_P2_202400452/internal/disk"
)

type ReportRequest struct {
	Type      string         `json:"type"`
	DiskPath  string         `json:"diskPath"`
	Path      string         `json:"path"`
	FilePath  string         `json:"filePath"`
	Partition disk.Partition `json:"-"`
}

type ReportResult struct {
	Type   string `json:"type"`
	Path   string `json:"path"`
	Format string `json:"format"`
}

func GenerateReport(req ReportRequest) (ReportResult, error) {
	reportType := NormalizeReportType(req.Type)

	if strings.TrimSpace(req.DiskPath) == "" {
		return ReportResult{}, fmt.Errorf("debe indicar la ruta del disco")
	}

	if strings.TrimSpace(req.Path) == "" {
		return ReportResult{}, fmt.Errorf("debe indicar la ruta de salida del reporte")
	}

	outputPath := strings.TrimSpace(req.Path)

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return ReportResult{}, fmt.Errorf("error al crear carpeta del reporte: %v", err)
	}

	switch reportType {
	case "mbr":
		if err := ReportMBR(req.DiskPath, outputPath); err != nil {
			return ReportResult{}, err
		}

	case "disk":
		if err := ReportDISK(req.DiskPath, outputPath); err != nil {
			return ReportResult{}, err
		}

	case "inode":
		if err := ReportINODE(req.DiskPath, req.Partition, outputPath); err != nil {
			return ReportResult{}, err
		}

	case "block":
		if err := ReportBLOCK(req.DiskPath, req.Partition, outputPath); err != nil {
			return ReportResult{}, err
		}

	case "bm_inode":
		if err := ReportBMNode(req.DiskPath, req.Partition, outputPath); err != nil {
			return ReportResult{}, err
		}

	case "bm_block":
		if err := ReportBMBlock(req.DiskPath, req.Partition, outputPath); err != nil {
			return ReportResult{}, err
		}

	case "sb":
		if err := ReportSB(req.DiskPath, req.Partition, outputPath); err != nil {
			return ReportResult{}, err
		}

	case "file":
		if strings.TrimSpace(req.FilePath) == "" {
			return ReportResult{}, fmt.Errorf("debe indicar filePath para el reporte file")
		}

		if err := ReportFILE(req.DiskPath, req.Partition, outputPath, req.FilePath); err != nil {
			return ReportResult{}, err
		}

	case "ls":
		filePath := strings.TrimSpace(req.FilePath)
		if filePath == "" {
			filePath = "/"
		}

		if err := ReportLS(req.DiskPath, req.Partition, outputPath, filePath); err != nil {
			return ReportResult{}, err
		}

	case "tree":
		if err := ReportTREE(req.DiskPath, req.Partition, outputPath); err != nil {
			return ReportResult{}, err
		}

	default:
		return ReportResult{}, fmt.Errorf("tipo de reporte '%s' no reconocido", req.Type)
	}

	return ReportResult{
		Type:   reportType,
		Path:   outputPath,
		Format: detectReportFormat(outputPath),
	}, nil
}

func NormalizeReportType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))

	switch value {
	case "bm_bloc":
		return "bm_block"
	case "bitmap_inodos":
		return "bm_inode"
	case "bitmap_bloques":
		return "bm_block"
	case "inodos":
		return "inode"
	case "bloques":
		return "block"
	default:
		return value
	}
}

func detectReportFormat(path string) string {
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
