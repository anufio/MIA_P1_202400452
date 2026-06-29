package services

//partition_delete_resize_service.go: contiene la lógica para manejar la eliminación y redimensionamiento de particiones dentro del sistema de archivos, incluyendo la validación de parámetros, la lectura y escritura de MBR y EBR, y la gestión de particiones primarias, extendidas y lógicas.

import (
	"fmt"
	"strings"

	"MIA_P2_202400452/internal/disk"
)

func (s *DiskService) DeletePartition(input DeletePartitionInput) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path := s.absolute(firstNonEmpty(input.DiskPath, input.Path))
	name := strings.TrimSpace(input.Name)
	deleteType := strings.ToLower(strings.TrimSpace(firstNonEmpty(input.DeleteType, input.Delete)))

	if deleteType == "" {
		deleteType = "fast"
	}

	if deleteType != "fast" && deleteType != "full" {
		return fmt.Errorf("deleteType debe ser fast o full")
	}

	if path == "" || name == "" {
		return fmt.Errorf("debe indicar ruta del disco y nombre de partición")
	}

	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return fmt.Errorf("error leyendo MBR: %v", err)
	}

	for i, partition := range mbr.MbrPartitions {
		if partition.PartStart <= 0 || partition.PartS <= 0 {
			continue
		}

		if equalName(cleanName(partition.PartName[:]), name) {
			if deleteType == "full" {
				if err := zeroRange(path, int64(partition.PartStart), int64(partition.PartS)); err != nil {
					return fmt.Errorf("error limpiando espacio de partición: %v", err)
				}
			}

			mbr.MbrPartitions[i] = disk.Partition{
				PartStatus: [1]byte{'0'},
				PartType:   [1]byte{'0'},
				PartFit:    [1]byte{'0'},
			}

			if err := disk.WriteMBR(path, mbr); err != nil {
				return fmt.Errorf("error escribiendo MBR: %v", err)
			}

			s.removeMountedByPathName(path, name)
			return nil
		}
	}

	if extended, _, found := disk.GetExtended(mbr); found {
		if err := s.deleteLogical(path, extended, name, deleteType); err != nil {
			return err
		}

		s.removeMountedByPathName(path, name)
		return nil
	}

	return fmt.Errorf("no existe una partición con el nombre '%s'", name)
}

func (s *DiskService) ResizePartition(input ResizePartitionInput) (PartitionInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path := s.absolute(firstNonEmpty(input.DiskPath, input.Path))
	name := strings.TrimSpace(input.Name)

	if path == "" || name == "" {
		return PartitionInfo{}, fmt.Errorf("debe indicar ruta del disco y nombre de partición")
	}

	unit := strings.ToUpper(strings.TrimSpace(input.Unit))
	if unit == "" {
		unit = "K"
	}

	addBytes, err := signedSizeToBytes(input.Add, unit)
	if err != nil {
		return PartitionInfo{}, err
	}

	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return PartitionInfo{}, fmt.Errorf("error leyendo MBR: %v", err)
	}

	for i, partition := range mbr.MbrPartitions {
		if partition.PartStart <= 0 || partition.PartS <= 0 {
			continue
		}

		if equalName(cleanName(partition.PartName[:]), name) {
			newSize := partition.PartS + addBytes

			if newSize <= 0 {
				return PartitionInfo{}, fmt.Errorf("el nuevo tamaño debe ser mayor que cero")
			}

			if addBytes > 0 {
				currentEnd := partition.PartStart + partition.PartS
				limit := mbr.MbrTamano

				for j, other := range mbr.MbrPartitions {
					if j == i || other.PartStart <= 0 || other.PartS <= 0 {
						continue
					}

					if other.PartStart >= currentEnd && other.PartStart < limit {
						limit = other.PartStart
					}
				}

				if currentEnd+addBytes > limit {
					return PartitionInfo{}, fmt.Errorf("no hay espacio libre contiguo suficiente")
				}
			}

			mbr.MbrPartitions[i].PartS = newSize

			if err := disk.WriteMBR(path, mbr); err != nil {
				return PartitionInfo{}, fmt.Errorf("error escribiendo MBR: %v", err)
			}

			return s.primaryPartitionInfo(path, EncodeDiskID(path), mbr.MbrPartitions[i]), nil
		}
	}

	if extended, _, found := disk.GetExtended(mbr); found {
		return s.resizeLogical(path, extended, name, addBytes)
	}

	return PartitionInfo{}, fmt.Errorf("no existe una partición con el nombre '%s'", name)
}

func (s *DiskService) deleteLogical(path string, extended disk.Partition, name string, deleteType string) error {
	records := s.readLogicalRecords(path, extended)

	targetIndex := -1

	for i, record := range records {
		if equalName(cleanName(record.ebr.PartName[:]), name) {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return fmt.Errorf("no existe una partición lógica con el nombre '%s'", name)
	}

	target := records[targetIndex]

	if deleteType == "full" {
		if err := zeroRange(path, int64(target.ebr.PartStart), int64(target.ebr.PartS)); err != nil {
			return err
		}
	}

	if targetIndex > 0 {
		previous := records[targetIndex-1]
		previous.ebr.PartNext = target.ebr.PartNext

		if err := disk.WriteEBR(path, int64(previous.position), previous.ebr); err != nil {
			return err
		}
	}

	tombstone := target.ebr
	tombstone.PartMount = [1]byte{'0'}
	tombstone.PartS = 0
	tombstone.PartName = [16]byte{}

	if targetIndex == 0 {
		tombstone.PartNext = target.ebr.PartNext
	} else {
		tombstone.PartNext = -1
	}

	return disk.WriteEBR(path, int64(target.position), tombstone)
}

func (s *DiskService) resizeLogical(path string, extended disk.Partition, name string, addBytes int32) (PartitionInfo, error) {
	records := s.readLogicalRecords(path, extended)

	for i, record := range records {
		if !equalName(cleanName(record.ebr.PartName[:]), name) {
			continue
		}

		newSize := record.ebr.PartS + addBytes
		if newSize <= 0 {
			return PartitionInfo{}, fmt.Errorf("el nuevo tamaño debe ser mayor que cero")
		}

		if addBytes > 0 {
			currentEnd := record.ebr.PartStart + record.ebr.PartS
			limit := extended.PartStart + extended.PartS

			if i < len(records)-1 {
				limit = records[i+1].position
			}

			if currentEnd+addBytes > limit {
				return PartitionInfo{}, fmt.Errorf("no hay espacio libre contiguo suficiente")
			}
		}

		record.ebr.PartS = newSize

		if err := disk.WriteEBR(path, int64(record.position), record.ebr); err != nil {
			return PartitionInfo{}, fmt.Errorf("error escribiendo EBR: %v", err)
		}

		return s.logicalPartitionInfo(path, EncodeDiskID(path), record.position, record.ebr), nil
	}

	return PartitionInfo{}, fmt.Errorf("no existe una partición lógica con el nombre '%s'", name)
}
