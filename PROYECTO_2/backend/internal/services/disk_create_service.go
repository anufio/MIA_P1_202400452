package services

//disk_create_service.go: contiene la lógica para manejar la creación de discos, incluyendo la validación de parámetros, la creación del archivo de disco y la escritura del MBR.

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"MIA_P2_202400452/internal/disk"
)

func (s *DiskService) CreateDisk(input CreateDiskInput) (DiskInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path, err := s.resolveDiskPath(input.Path, input.Name)
	if err != nil {
		return DiskInfo{}, err
	}

	if !validDiskExtension(path) {
		return DiskInfo{}, fmt.Errorf("el archivo de disco debe tener extensión .dsk o .mia")
	}

	if disk.DiskExists(path) {
		return DiskInfo{}, fmt.Errorf("el disco ya existe: %s", path)
	}

	unit := strings.ToUpper(strings.TrimSpace(input.Unit))
	if unit == "" {
		unit = "M"
	}

	sizeBytes, err := positiveSizeToBytes(input.Size, unit)
	if err != nil {
		return DiskInfo{}, err
	}

	fit := strings.ToUpper(strings.TrimSpace(input.Fit))
	if fit == "" {
		fit = "FF"
	}

	fitByte, err := disk.GetFitChar(fit)
	if err != nil {
		return DiskInfo{}, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return DiskInfo{}, fmt.Errorf("error creando carpeta del disco: %v", err)
	}

	if err := disk.CreateDisk(path, int64(sizeBytes)); err != nil {
		return DiskInfo{}, err
	}

	rand.Seed(time.Now().UnixNano())

	mbr := disk.MBR{
		MbrTamano:        sizeBytes,
		MbrFechaCreacion: disk.GetFechaCreacion(),
		MbrDskSignature:  rand.Int31(),
		DskFit:           [1]byte{fitByte},
	}

	for i := range mbr.MbrPartitions {
		mbr.MbrPartitions[i] = disk.Partition{
			PartStatus: [1]byte{'0'},
			PartType:   [1]byte{'0'},
			PartFit:    [1]byte{'0'},
		}
	}

	if err := disk.WriteMBR(path, mbr); err != nil {
		return DiskInfo{}, fmt.Errorf("error escribiendo MBR: %v", err)
	}

	return s.diskInfoFromPath(path)
}
