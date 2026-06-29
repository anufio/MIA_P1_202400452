package services

// disk_delete_service.go: elimina discos y limpia los reportes generados para que la GUI no muestre reportes viejos.

import (
	"fmt"
	"os"
	"path/filepath"

	"MIA_P2_202400452/internal/disk"
)

func (s *DiskService) DeleteDisk(path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path = s.absolute(path)

	if path == "" {
		return fmt.Errorf("debe indicar la ruta del disco")
	}

	if !disk.DiskExists(path) {
		return fmt.Errorf("el disco no existe: %s", path)
	}

	for id, mounted := range s.mounted {
		if mounted.DiskPath == path {
			delete(s.mounted, id)
		}
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error eliminando disco: %v", err)
	}

	reportsRoot := filepath.Join(filepath.Dir(s.root), "reports")

	if err := os.RemoveAll(reportsRoot); err != nil {
		return fmt.Errorf("disco eliminado, pero no se pudieron limpiar los reportes: %v", err)
	}

	if err := os.MkdirAll(reportsRoot, 0755); err != nil {
		return fmt.Errorf("disco eliminado, pero no se pudo recrear la carpeta de reportes: %v", err)
	}

	return nil
}
