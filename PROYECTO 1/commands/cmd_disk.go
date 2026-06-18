package commands

// cmd_disk.go

import (
	"MIA_P1_202400452/disk"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CmdMKDISK(params map[string]string) string {
	sizeStr, ok := params["size"]
	if !ok {
		return "Error: falta el parámetro obligatorio -size"
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		return "Error: -size debe ser un número positivo mayor que cero"
	}

	path, ok := params["path"]
	if !ok {
		return "Error: falta el parámetro obligatorio -path"
	}

	if strings.TrimSpace(path) == "" {
		return "Error: -path no puede estar vacío"
	}

	if !hasValidDiskExtension(path) {
		return "Error: el archivo debe tener extensión .mia o .dsk"
	}

	unit := "M"
	if u, ok := params["unit"]; ok {
		unit = strings.ToUpper(strings.TrimSpace(u))
		if unit != "K" && unit != "M" {
			return fmt.Sprintf("Error: unidad inválida '%s'. Use K o M", params["unit"])
		}
	}

	fit := "FF"
	if f, ok := params["fit"]; ok {
		fit = strings.ToUpper(strings.TrimSpace(f))
		if fit != "BF" && fit != "FF" && fit != "WF" {
			return fmt.Sprintf("Error: ajuste inválido '%s'. Use BF, FF o WF", params["fit"])
		}
	}

	var sizeBytes int64
	switch unit {
	case "K":
		sizeBytes = int64(size) * 1024
	case "M":
		sizeBytes = int64(size) * 1024 * 1024
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Sprintf("Error al crear directorio del disco: %v", err)
		}
	}

	if err := disk.CreateDisk(path, sizeBytes); err != nil {
		return fmt.Sprintf("Error al crear disco: %v", err)
	}

	fitChar, _ := disk.GetFitChar(fit)

	rand.Seed(time.Now().UnixNano())

	mbr := disk.MBR{
		MbrTamano:        int32(sizeBytes),
		MbrFechaCreacion: disk.GetFechaCreacion(),
		MbrDskSignature:  rand.Int31(),
		DskFit:           [1]byte{fitChar},
	}

	for i := range mbr.MbrPartitions {
		mbr.MbrPartitions[i].PartStatus = [1]byte{'0'}
		mbr.MbrPartitions[i].PartType = [1]byte{'0'}
		mbr.MbrPartitions[i].PartFit = [1]byte{'0'}
		mbr.MbrPartitions[i].PartStart = 0
		mbr.MbrPartitions[i].PartS = 0
		mbr.MbrPartitions[i].PartCorrelative = 0
		mbr.MbrPartitions[i].PartId = [4]byte{}
	}

	if err := disk.WriteMBR(path, mbr); err != nil {
		return fmt.Sprintf("Error al escribir MBR: %v", err)
	}

	return fmt.Sprintf("Disco creado exitosamente: %s (%.2f %s)", path, float64(size), unit)
}

func CmdRMDISK(params map[string]string) string {
	path, ok := params["path"]
	if !ok {
		return "Error: falta el parámetro obligatorio -path"
	}

	if strings.TrimSpace(path) == "" {
		return "Error: -path no puede estar vacío"
	}

	if !disk.DiskExists(path) {
		return fmt.Sprintf("Error: el disco '%s' no existe", path)
	}

	fmt.Printf("¿Está seguro que desea eliminar el disco '%s'? (s/n): ", path)

	var resp string
	fmt.Scanln(&resp)

	resp = strings.ToLower(strings.TrimSpace(resp))

	if resp != "s" && resp != "si" && resp != "sí" && resp != "y" && resp != "yes" {
		return "Operación cancelada."
	}

	if err := os.Remove(path); err != nil {
		return fmt.Sprintf("Error al eliminar disco: %v", err)
	}

	return fmt.Sprintf("Disco '%s' eliminado exitosamente.", path)
}

func hasValidDiskExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mia" || ext == ".dsk"
}
