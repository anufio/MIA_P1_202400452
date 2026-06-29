package services

//partition_create_service.go: contiene la lógica para manejar la creación de particiones dentro del sistema de archivos, incluyendo la validación de parámetros, la escritura de MBR y EBR, y la gestión de particiones primarias, extendidas y lógicas.

import (
	"fmt"
	"sort"
	"strings"

	"MIA_P2_202400452/internal/disk"
)

func (s *DiskService) CreatePartition(input CreatePartitionInput) (PartitionInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	path := s.absolute(firstNonEmpty(input.DiskPath, input.Path))
	name := strings.TrimSpace(input.Name)

	if path == "" {
		return PartitionInfo{}, fmt.Errorf("debe indicar la ruta del disco")
	}

	if name == "" {
		return PartitionInfo{}, fmt.Errorf("debe indicar el nombre de la partición")
	}

	if len(name) > 16 {
		return PartitionInfo{}, fmt.Errorf("el nombre de la partición no puede exceder 16 caracteres")
	}

	if !disk.DiskExists(path) {
		return PartitionInfo{}, fmt.Errorf("el disco no existe: %s", path)
	}

	partType := strings.ToUpper(strings.TrimSpace(input.Type))
	if partType == "" {
		partType = "P"
	}

	if partType != "P" && partType != "E" && partType != "L" {
		return PartitionInfo{}, fmt.Errorf("tipo de partición inválido '%s'. Use P, E o L", partType)
	}

	unit := strings.ToUpper(strings.TrimSpace(input.Unit))
	if unit == "" {
		unit = "K"
	}

	sizeBytes, err := positiveSizeToBytes(input.Size, unit)
	if err != nil {
		return PartitionInfo{}, err
	}

	fit := strings.ToUpper(strings.TrimSpace(input.Fit))
	if fit == "" {
		fit = "WF"
	}

	fitByte, err := disk.GetFitChar(fit)
	if err != nil {
		return PartitionInfo{}, err
	}

	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return PartitionInfo{}, fmt.Errorf("error leyendo MBR: %v", err)
	}

	if s.partitionNameExists(path, mbr, name) {
		return PartitionInfo{}, fmt.Errorf("ya existe una partición con el nombre '%s'", name)
	}

	switch partType {
	case "P":
		return s.createPrimary(path, mbr, name, sizeBytes, fitByte)

	case "E":
		return s.createExtended(path, mbr, name, sizeBytes, fitByte)

	case "L":
		return s.createLogical(path, mbr, name, sizeBytes, fitByte)
	}

	return PartitionInfo{}, fmt.Errorf("tipo de partición inválido")
}

func (s *DiskService) createPrimary(path string, mbr disk.MBR, name string, sizeBytes int32, fitByte byte) (PartitionInfo, error) {
	if disk.CountPrimaryExtended(mbr) >= 4 {
		return PartitionInfo{}, fmt.Errorf("ya existen 4 particiones primarias/extendidas")
	}

	start := disk.FindFreeSpace(mbr, sizeBytes, fitByte)
	if start == -1 {
		return PartitionInfo{}, fmt.Errorf("no hay espacio suficiente para la partición")
	}

	slot := disk.FindFreePartitionSlot(mbr)
	if slot == -1 {
		return PartitionInfo{}, fmt.Errorf("no hay slots disponibles en el MBR")
	}

	correlative := disk.GetNextCorrelative(mbr)

	partition := disk.Partition{
		PartStatus:      [1]byte{'0'},
		PartType:        [1]byte{'P'},
		PartFit:         [1]byte{fitByte},
		PartStart:       start,
		PartS:           sizeBytes,
		PartName:        copyName16(name),
		PartCorrelative: correlative,
	}

	mbr.MbrPartitions[slot] = partition

	if err := disk.WriteMBR(path, mbr); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo MBR: %v", err)
	}

	return s.primaryPartitionInfo(path, EncodeDiskID(path), partition), nil
}

func (s *DiskService) createExtended(path string, mbr disk.MBR, name string, sizeBytes int32, fitByte byte) (PartitionInfo, error) {
	if disk.HasExtended(mbr) {
		return PartitionInfo{}, fmt.Errorf("ya existe una partición extendida")
	}

	if disk.CountPrimaryExtended(mbr) >= 4 {
		return PartitionInfo{}, fmt.Errorf("ya existen 4 particiones primarias/extendidas")
	}

	start := disk.FindFreeSpace(mbr, sizeBytes, fitByte)
	if start == -1 {
		return PartitionInfo{}, fmt.Errorf("no hay espacio suficiente para la partición extendida")
	}

	slot := disk.FindFreePartitionSlot(mbr)
	if slot == -1 {
		return PartitionInfo{}, fmt.Errorf("no hay slots disponibles en el MBR")
	}

	correlative := disk.GetNextCorrelative(mbr)

	partition := disk.Partition{
		PartStatus:      [1]byte{'0'},
		PartType:        [1]byte{'E'},
		PartFit:         [1]byte{fitByte},
		PartStart:       start,
		PartS:           sizeBytes,
		PartName:        copyName16(name),
		PartCorrelative: correlative,
	}

	mbr.MbrPartitions[slot] = partition

	if err := disk.WriteMBR(path, mbr); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo MBR: %v", err)
	}

	firstEBR := disk.EBR{
		PartMount: [1]byte{'0'},
		PartFit:   [1]byte{fitByte},
		PartStart: start,
		PartS:     0,
		PartNext:  -1,
	}

	if err := disk.WriteEBR(path, int64(start), firstEBR); err != nil {
		return PartitionInfo{}, fmt.Errorf("error escribiendo EBR inicial: %v", err)
	}

	return s.primaryPartitionInfo(path, EncodeDiskID(path), partition), nil
}

func (s *DiskService) createLogical(path string, mbr disk.MBR, name string, sizeBytes int32, fitByte byte) (PartitionInfo, error) {
	extended, _, found := disk.GetExtended(mbr)
	if !found {
		return PartitionInfo{}, fmt.Errorf("no existe una partición extendida")
	}

	ebrSize := int32(disk.SizeEBR())
	required := sizeBytes + ebrSize

	records := s.readLogicalRecords(path, extended)
	spaces := s.logicalFreeSpaces(extended, records)
	selected := s.selectLogicalSpace(spaces, required, fitByte)

	if selected.start == -1 {
		return PartitionInfo{}, fmt.Errorf("no hay espacio suficiente dentro de la partición extendida")
	}

	newEBRPosition := selected.start
	logicalStart := newEBRPosition + ebrSize

	ebr := disk.EBR{
		PartMount: [1]byte{'0'},
		PartFit:   [1]byte{fitByte},
		PartStart: logicalStart,
		PartS:     sizeBytes,
		PartNext:  -1,
		PartName:  copyName16(name),
	}

	allRecords := append(records, logicalRecord{
		position: newEBRPosition,
		ebr:      ebr,
	})

	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].position < allRecords[j].position
	})

	for i := range allRecords {
		if i == len(allRecords)-1 {
			allRecords[i].ebr.PartNext = -1
		} else {
			allRecords[i].ebr.PartNext = allRecords[i+1].position
		}

		if err := disk.WriteEBR(path, int64(allRecords[i].position), allRecords[i].ebr); err != nil {
			return PartitionInfo{}, fmt.Errorf("error escribiendo EBR: %v", err)
		}
	}

	return s.logicalPartitionInfo(path, EncodeDiskID(path), newEBRPosition, ebr), nil
}

func (s *DiskService) partitionNameExists(path string, mbr disk.MBR, name string) bool {
	for _, partition := range mbr.MbrPartitions {
		if partition.PartStart <= 0 || partition.PartS <= 0 {
			continue
		}

		if equalName(cleanName(partition.PartName[:]), name) {
			return true
		}
	}

	if extended, _, found := disk.GetExtended(mbr); found {
		_, _, logicalFound := disk.FindLogicalByName(path, extended, name)
		return logicalFound
	}

	return false
}
