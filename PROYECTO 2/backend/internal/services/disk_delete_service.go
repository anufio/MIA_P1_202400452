package services

//disk_delete_service.go: contiene la lógica para manejar la eliminación de discos, incluyendo la validación de la existencia del disco y la eliminación del archivo de disco.

import (
	"fmt"
	"os"

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

	return nil
}
