package reports

// rep.go: comando REP para generar reportes

import (
	"MIA_P1_202400452/commands"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CmdREP(params map[string]string) string {
	name, ok := params["name"]
	if !ok {
		return "Error: falta el parámetro obligatorio -name"
	}

	name = strings.ToLower(strings.TrimSpace(name))

	if name == "bm_bloc" {
		name = "bm_block"
	}

	path, ok := params["path"]
	if !ok || strings.TrimSpace(path) == "" {
		return "Error: falta el parámetro obligatorio -path"
	}

	id, ok := params["id"]
	if !ok || strings.TrimSpace(id) == "" {
		return "Error: falta el parámetro obligatorio -id"
	}

	path = strings.TrimSpace(path)
	id = strings.ToUpper(strings.TrimSpace(id))

	diskPath, partition, _, err := commands.GetMountedPartitionInfo(id)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Sprintf("Error al crear carpeta del reporte: %v", err)
	}

	switch name {
	case "mbr":
		return handleReport(func() error {
			return ReportMBR(diskPath, path)
		})

	case "disk":
		return handleReport(func() error {
			return ReportDISK(diskPath, path)
		})

	case "inode":
		return handleReport(func() error {
			return ReportINODE(diskPath, partition, path)
		})

	case "block":
		return handleReport(func() error {
			return ReportBLOCK(diskPath, partition, path)
		})

	case "bm_inode":
		return handleReport(func() error {
			return ReportBMNode(diskPath, partition, path)
		})

	case "bm_block":
		return handleReport(func() error {
			return ReportBMBlock(diskPath, partition, path)
		})

	case "sb":
		return handleReport(func() error {
			return ReportSB(diskPath, partition, path)
		})

	case "file":
		pathFile, ok := params["path_file_ls"]
		if !ok || strings.TrimSpace(pathFile) == "" {
			return "Error: falta el parámetro obligatorio -path_file_ls para reporte file"
		}

		pathFile = strings.TrimSpace(pathFile)

		return handleReport(func() error {
			return ReportFILE(diskPath, partition, path, pathFile)
		})

	case "ls":
		pathFile := strings.TrimSpace(params["path_file_ls"])
		if pathFile == "" {
			pathFile = "/"
		}

		return handleReport(func() error {
			return ReportLS(diskPath, partition, path, pathFile)
		})

	case "tree":
		return handleReport(func() error {
			return ReportTREE(diskPath, partition, path)
		})

	default:
		return fmt.Sprintf("Error: tipo de reporte '%s' no reconocido", name)
	}
}

func handleReport(fn func() error) string {
	if err := fn(); err != nil {
		return fmt.Sprintf("Error generando reporte: %v", err)
	}

	return "Reporte generado exitosamente"
}
