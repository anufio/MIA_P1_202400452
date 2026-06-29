package services

//disk_read_service.go: contiene la lógica para manejar la lectura de discos, incluyendo la validación de la existencia del disco y la obtención de información del MBR.

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"MIA_P2_202400452/internal/disk"
)

func (s *DiskService) ListDisks() ([]DiskInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var result []DiskInfo

	err := filepath.WalkDir(s.root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		if !validDiskExtension(path) {
			return nil
		}

		info, err := s.diskInfoFromPath(path)
		if err == nil {
			result = append(result, info)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func (s *DiskService) diskInfoFromPath(path string) (DiskInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return DiskInfo{}, err
	}

	info := DiskInfo{
		ID:   EncodeDiskID(path),
		Name: filepath.Base(path),
		Path: path,
		Size: stat.Size(),
		Fit:  "FF",
	}

	mbr, err := disk.ReadMBR(path)
	if err == nil {
		info.Fit = byteToFit(mbr.DskFit[0])
		info.CreatedAt = cleanName(mbr.MbrFechaCreacion[:])
	}

	return info, nil
}

func (s *DiskService) absolute(path string) string {
	path = strings.TrimSpace(path)

	if path == "" {
		return ""
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return abs
}

func (s *DiskService) resolveDiskPath(path string, name string) (string, error) {
	path = strings.TrimSpace(path)
	name = strings.TrimSpace(name)

	if path != "" {
		return s.absolute(path), nil
	}

	if name == "" {
		name = "disco.dsk"
	}

	if filepath.Ext(name) == "" {
		name += ".dsk"
	}

	return filepath.Join(s.root, name), nil
}
