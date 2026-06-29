package services

//partition_read_service.go: contiene la lógica para manejar la lectura de particiones dentro del sistema de archivos, incluyendo la obtención de información de particiones primarias, extendidas y lógicas, así como la verificación de su estado y formato.

import (
	"fmt"
	"sort"

	"MIA_P2_202400452/internal/disk"
)

type logicalRecord struct {
	position int32
	ebr      disk.EBR
}

func (s *DiskService) ListPartitions(path string) ([]PartitionInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path = s.absolute(path)

	if !disk.DiskExists(path) {
		return nil, fmt.Errorf("el disco no existe: %s", path)
	}

	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo MBR: %v", err)
	}

	diskID := EncodeDiskID(path)

	var result []PartitionInfo

	for _, partition := range mbr.MbrPartitions {
		if partition.PartStart <= 0 || partition.PartS <= 0 {
			continue
		}

		result = append(result, s.primaryPartitionInfo(path, diskID, partition))
	}

	if extended, _, found := disk.GetExtended(mbr); found {
		records := s.readLogicalRecords(path, extended)

		for _, record := range records {
			result = append(result, s.logicalPartitionInfo(path, diskID, record.position, record.ebr))
		}
	}

	return result, nil
}

func (s *DiskService) primaryPartitionInfo(path string, diskID string, partition disk.Partition) PartitionInfo {
	id := cleanName(partition.PartId[:])
	name := cleanName(partition.PartName[:])

	if id == "" {
		id = EncodeDiskID(path + "::" + name)
	}

	return PartitionInfo{
		ID:          id,
		DiskID:      diskID,
		DiskPath:    path,
		Name:        name,
		Type:        string(partition.PartType[0]),
		Fit:         byteToFit(partition.PartFit[0]),
		Start:       partition.PartStart,
		Size:        partition.PartS,
		Correlative: partition.PartCorrelative,
		Mounted:     partition.PartStatus[0] == '1',
		Formatted:   isFormatted(path, partition.PartStart),
	}
}

func (s *DiskService) logicalPartitionInfo(path string, diskID string, ebrPosition int32, ebr disk.EBR) PartitionInfo {
	name := cleanName(ebr.PartName[:])
	id := EncodeDiskID(fmt.Sprintf("%s::%s::%d", path, name, ebrPosition))

	return PartitionInfo{
		ID:        id,
		DiskID:    diskID,
		DiskPath:  path,
		Name:      name,
		Type:      "L",
		Fit:       byteToFit(ebr.PartFit[0]),
		Start:     ebr.PartStart,
		Size:      ebr.PartS,
		Mounted:   ebr.PartMount[0] == '1',
		Formatted: isFormatted(path, ebr.PartStart),
	}
}

func (s *DiskService) readLogicalRecords(path string, extended disk.Partition) []logicalRecord {
	var records []logicalRecord

	extendedStart := extended.PartStart
	extendedEnd := extended.PartStart + extended.PartS
	position := extendedStart
	visited := make(map[int32]bool)

	for position >= extendedStart && position < extendedEnd && !visited[position] {
		visited[position] = true

		ebr, err := disk.ReadEBR(path, int64(position))
		if err != nil {
			break
		}

		if ebr.PartS > 0 {
			records = append(records, logicalRecord{
				position: position,
				ebr:      ebr,
			})
		}

		if ebr.PartNext == -1 {
			break
		}

		position = ebr.PartNext
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].position < records[j].position
	})

	return records
}
