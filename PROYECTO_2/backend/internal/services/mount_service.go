package services

//mount_service.go: contiene la lógica para manejar el montaje y desmontaje de particiones dentro del sistema de archivos, incluyendo la lectura y escritura de MBR y EBR, la gestión de particiones montadas y la generación de identificadores únicos.

import (
	"fmt"
	"sort"
	"strings"

	"MIA_P2_202400452/internal/disk"
)

func (s *DiskService) Mount(input MountInput) (MountedInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path := s.absolute(firstNonEmpty(input.DiskPath, input.Path))
	name := strings.TrimSpace(input.Name)

	if path == "" || name == "" {
		return MountedInfo{}, fmt.Errorf("debe indicar ruta del disco y nombre de partición")
	}

	for _, mounted := range s.mounted {
		if mounted.DiskPath == path && equalName(mounted.Name, name) {
			return mounted, nil
		}
	}

	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return MountedInfo{}, fmt.Errorf("error leyendo MBR: %v", err)
	}

	for i, partition := range mbr.MbrPartitions {
		if partition.PartStart <= 0 || partition.PartS <= 0 {
			continue
		}

		partName := cleanName(partition.PartName[:])

		if !equalName(partName, name) {
			continue
		}

		if partition.PartType[0] == 'E' {
			return MountedInfo{}, fmt.Errorf("las particiones extendidas no se montan")
		}

		n := partition.PartCorrelative
		if n <= 0 {
			n = int32(i + 1)
		}

		id := s.buildMountID(path, n)

		partition.PartStatus = [1]byte{'1'}
		partition.PartId = copyID4(id)
		mbr.MbrPartitions[i] = partition

		if err := disk.WriteMBR(path, mbr); err != nil {
			return MountedInfo{}, fmt.Errorf("error actualizando MBR: %v", err)
		}

		mounted := MountedInfo{
			ID:        id,
			DiskPath:  path,
			Name:      partName,
			Type:      string(partition.PartType[0]),
			Start:     partition.PartStart,
			Size:      partition.PartS,
			Index:     i,
			IsLogical: false,
			EBRStart:  -1,
		}

		s.mounted[id] = mounted
		return mounted, nil
	}

	if extended, _, found := disk.GetExtended(mbr); found {
		ebr, ebrStart, found := disk.FindLogicalByName(path, extended, name)
		if found {
			id := s.buildMountID(path, int32(len(s.mounted)+1))

			ebr.PartMount = [1]byte{'1'}

			if err := disk.WriteEBR(path, ebrStart, ebr); err != nil {
				return MountedInfo{}, fmt.Errorf("error actualizando EBR: %v", err)
			}

			mounted := MountedInfo{
				ID:        id,
				DiskPath:  path,
				Name:      cleanName(ebr.PartName[:]),
				Type:      "L",
				Start:     ebr.PartStart,
				Size:      ebr.PartS,
				Index:     -1,
				IsLogical: true,
				EBRStart:  ebrStart,
			}

			s.mounted[id] = mounted
			return mounted, nil
		}
	}

	return MountedInfo{}, fmt.Errorf("no existe una partición primaria o lógica con el nombre '%s'", name)
}

func (s *DiskService) Unmount(input UnmountInput) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := strings.ToUpper(strings.TrimSpace(input.ID))

	mounted, ok := s.mounted[id]
	if !ok {
		return fmt.Errorf("el ID '%s' no está montado", id)
	}

	if mounted.IsLogical {
		ebr, err := disk.ReadEBR(mounted.DiskPath, mounted.EBRStart)
		if err == nil {
			ebr.PartMount = [1]byte{'0'}
			_ = disk.WriteEBR(mounted.DiskPath, mounted.EBRStart, ebr)
		}
	} else {
		mbr, err := disk.ReadMBR(mounted.DiskPath)
		if err == nil && mounted.Index >= 0 && mounted.Index < len(mbr.MbrPartitions) {
			mbr.MbrPartitions[mounted.Index].PartStatus = [1]byte{'0'}
			mbr.MbrPartitions[mounted.Index].PartId = [4]byte{}
			_ = disk.WriteMBR(mounted.DiskPath, mbr)
		}
	}

	delete(s.mounted, id)

	return nil
}

func (s *DiskService) ListMounted() []MountedInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var result []MountedInfo

	for _, mounted := range s.mounted {
		result = append(result, mounted)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func (s *DiskService) GetMounted(id string) (MountedInfo, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id = strings.ToUpper(strings.TrimSpace(id))
	mounted, ok := s.mounted[id]

	return mounted, ok
}

func (s *DiskService) buildMountID(path string, partitionNumber int32) string {
	letter := s.diskLetter(path)

	return strings.ToUpper(fmt.Sprintf("52%d%s", partitionNumber, letter))
}

func (s *DiskService) diskLetter(path string) string {
	if letter, ok := s.diskLetters[path]; ok {
		return letter
	}

	letter := string(rune('A' + s.nextLetter))
	s.diskLetters[path] = letter
	s.nextLetter++

	return letter
}

func (s *DiskService) removeMountedByPathName(path string, name string) {
	for id, mounted := range s.mounted {
		if mounted.DiskPath == path && equalName(mounted.Name, name) {
			delete(s.mounted, id)
		}
	}
}
