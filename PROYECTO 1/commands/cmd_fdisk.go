package commands

// cmd_fdisk.go

import (
	"MIA_P1_202400452/disk"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func CmdFDISK(params map[string]string) string {
	path, ok := params["path"]
	if !ok {
		return "Error: falta el parámetro obligatorio -path"
	}

	if !disk.DiskExists(path) {
		return fmt.Sprintf("Error: el disco '%s' no existe", path)
	}

	name, ok := params["name"]
	if !ok {
		return "Error: falta el parámetro obligatorio -name"
	}
	if len(name) > 16 {
		return "Error: el nombre de la partición no puede exceder 16 caracteres"
	}

	partType := "P"
	if t, ok := params["type"]; ok {
		partType = strings.ToUpper(t)
		if partType != "P" && partType != "E" && partType != "L" {
			return fmt.Sprintf("Error: tipo de partición inválido '%s'. Use P, E o L", params["type"])
		}
	}

	fit := "WF"
	if f, ok := params["fit"]; ok {
		fit = strings.ToUpper(f)
		if fit != "BF" && fit != "FF" && fit != "WF" {
			return fmt.Sprintf("Error: ajuste inválido '%s'. Use BF, FF o WF", params["fit"])
		}
	}

	unit := "K"
	if u, ok := params["unit"]; ok {
		unit = strings.ToUpper(u)
		if unit != "B" && unit != "K" && unit != "M" {
			return fmt.Sprintf("Error: unidad inválida '%s'. Use B, K o M", params["unit"])
		}
	}

	sizeStr, hasSz := params["size"]
	if !hasSz {
		return "Error: falta el parámetro obligatorio -size para crear una partición"
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		return "Error: -size debe ser un número positivo mayor que cero"
	}

	var sizeBytes int32
	switch unit {
	case "B":
		sizeBytes = int32(size)
	case "K":
		sizeBytes = int32(size) * 1024
	case "M":
		sizeBytes = int32(size) * 1024 * 1024
	}

	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return fmt.Sprintf("Error al leer MBR: %v", err)
	}

	fitChar, _ := disk.GetFitChar(fit)

	switch partType {
	case "P":
		return createPrimary(path, mbr, name, sizeBytes, fitChar)
	case "E":
		return createExtended(path, mbr, name, sizeBytes, fitChar)
	case "L":
		return createLogical(path, mbr, name, sizeBytes, fitChar)
	}
	return "Error desconocido en fdisk"
}

func createPrimary(path string, mbr disk.MBR, name string, sizeBytes int32, fitChar byte) string {

	if partitionNameExists(path, mbr, name) {
		return fmt.Sprintf("Error: ya existe una partición con el nombre '%s'", name)
	}

	if disk.CountPrimaryExtended(mbr) >= 4 {
		return "Error: ya existen 4 particiones (primarias + extendidas), no se puede crear más"
	}

	start := disk.FindFreeSpace(mbr, sizeBytes, fitChar)
	if start == -1 {
		return "Error: no hay espacio suficiente en el disco para la partición"
	}

	slot := disk.FindFreePartitionSlot(mbr)
	if slot == -1 {
		return "Error: no hay slots disponibles en el MBR"
	}

	correlative := disk.GetNextCorrelative(mbr)

	var nameArr [16]byte
	copy(nameArr[:], name)

	mbr.MbrPartitions[slot] = disk.Partition{
		PartStatus:      [1]byte{'0'},
		PartType:        [1]byte{'P'},
		PartFit:         [1]byte{fitChar},
		PartStart:       start,
		PartS:           sizeBytes,
		PartName:        nameArr,
		PartCorrelative: correlative,
		PartId:          [4]byte{},
	}

	if err := disk.WriteMBR(path, mbr); err != nil {
		return fmt.Sprintf("Error al escribir MBR: %v", err)
	}

	return fmt.Sprintf("Partición primaria '%s' creada exitosamente (inicio: %d, tamaño: %d bytes)", name, start, sizeBytes)
}

func createExtended(path string, mbr disk.MBR, name string, sizeBytes int32, fitChar byte) string {
	if disk.HasExtended(mbr) {
		return "Error: ya existe una partición extendida en este disco"
	}
	if partitionNameExists(path, mbr, name) {
		return fmt.Sprintf("Error: ya existe una partición con el nombre '%s'", name)
	}
	if disk.CountPrimaryExtended(mbr) >= 4 {
		return "Error: ya existen 4 particiones, no se puede crear más"
	}
	start := disk.FindFreeSpace(mbr, sizeBytes, fitChar)
	if start == -1 {
		return "Error: no hay espacio suficiente en el disco"
	}
	slot := disk.FindFreePartitionSlot(mbr)
	if slot == -1 {
		return "Error: no hay slots disponibles en el MBR"
	}
	correlative := disk.GetNextCorrelative(mbr)

	var nameArr [16]byte
	copy(nameArr[:], name)

	mbr.MbrPartitions[slot] = disk.Partition{
		PartStatus:      [1]byte{'0'},
		PartType:        [1]byte{'E'},
		PartFit:         [1]byte{fitChar},
		PartStart:       start,
		PartS:           sizeBytes,
		PartName:        nameArr,
		PartCorrelative: correlative,
		PartId:          [4]byte{},
	}

	if err := disk.WriteMBR(path, mbr); err != nil {
		return fmt.Sprintf("Error al escribir MBR: %v", err)
	}

	firstEBR := disk.EBR{
		PartMount: [1]byte{'0'},
		PartFit:   [1]byte{fitChar},
		PartStart: start,
		PartS:     0,
		PartNext:  -1,
	}
	if err := disk.WriteEBR(path, int64(start), firstEBR); err != nil {
		return fmt.Sprintf("Error al escribir EBR inicial: %v", err)
	}

	return fmt.Sprintf("Partición extendida '%s' creada exitosamente (inicio: %d, tamaño: %d bytes)", name, start, sizeBytes)
}

func createLogical(path string, mbr disk.MBR, name string, sizeBytes int32, fitChar byte) string {
	ext, _, found := disk.GetExtended(mbr)
	if !found {
		return "Error: no existe una partición extendida. Cree una partición extendida primero"
	}

	if partitionNameExists(path, mbr, name) {
		return fmt.Sprintf("Error: ya existe una partición con el nombre '%s'", name)
	}

	ebrSize := int32(disk.SizeEBR())
	if sizeBytes+ebrSize > ext.PartS {
		return "Error: no hay espacio suficiente dentro de la partición extendida"
	}

	records := readLogicalRecords(path, ext)
	freeSpaces := logicalFreeSpaces(ext, records)
	selected := selectLogicalSpace(freeSpaces, sizeBytes+ebrSize, fitChar)
	if selected.start == -1 {
		return "Error: no hay espacio suficiente dentro de la partición extendida"
	}

	newEBRPos := selected.start
	logStart := newEBRPos + ebrSize

	var nameArr [16]byte
	copy(nameArr[:], name)

	newEBR := disk.EBR{
		PartMount: [1]byte{'0'},
		PartFit:   [1]byte{fitChar},
		PartStart: logStart,
		PartS:     sizeBytes,
		PartNext:  -1,
		PartName:  nameArr,
	}

	allRecords := append(records, logicalRecord{pos: newEBRPos, ebr: newEBR})
	sort.Slice(allRecords, func(i, j int) bool {
		return allRecords[i].pos < allRecords[j].pos
	})

	for i := range allRecords {
		if i == len(allRecords)-1 {
			allRecords[i].ebr.PartNext = -1
		} else {
			allRecords[i].ebr.PartNext = allRecords[i+1].pos
		}
		if err := disk.WriteEBR(path, int64(allRecords[i].pos), allRecords[i].ebr); err != nil {
			return fmt.Sprintf("Error al escribir EBR: %v", err)
		}
	}

	return fmt.Sprintf("Partición lógica '%s' creada exitosamente (inicio: %d, tamaño: %d bytes)", name, logStart, sizeBytes)
}

type logicalRecord struct {
	pos int32
	ebr disk.EBR
}

type logicalSpace struct {
	start int32
	size  int32
}

func partitionNameExists(path string, mbr disk.MBR, name string) bool {
	if _, _, found := disk.FindPartitionByName(mbr, name); found {
		return true
	}
	if ext, _, found := disk.GetExtended(mbr); found {
		_, _, logicalFound := disk.FindLogicalByName(path, ext, name)
		return logicalFound
	}
	return false
}

func readLogicalRecords(path string, ext disk.Partition) []logicalRecord {
	var records []logicalRecord
	extStart := ext.PartStart
	extEnd := ext.PartStart + ext.PartS
	pos := extStart

	visited := make(map[int32]bool)
	for pos >= extStart && pos < extEnd && !visited[pos] {
		visited[pos] = true
		ebr, err := disk.ReadEBR(path, int64(pos))
		if err != nil {
			break
		}
		if ebr.PartS > 0 {
			records = append(records, logicalRecord{pos: pos, ebr: ebr})
		}
		if ebr.PartNext == -1 {
			break
		}
		pos = ebr.PartNext
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].pos < records[j].pos
	})
	return records
}

func logicalFreeSpaces(ext disk.Partition, records []logicalRecord) []logicalSpace {
	var spaces []logicalSpace
	lastEnd := ext.PartStart
	extEnd := ext.PartStart + ext.PartS

	for _, rec := range records {
		if rec.pos > lastEnd {
			spaces = append(spaces, logicalSpace{start: lastEnd, size: rec.pos - lastEnd})
		}
		recEnd := rec.ebr.PartStart + rec.ebr.PartS
		if recEnd > lastEnd {
			lastEnd = recEnd
		}
	}

	if lastEnd < extEnd {
		spaces = append(spaces, logicalSpace{start: lastEnd, size: extEnd - lastEnd})
	}
	return spaces
}

func selectLogicalSpace(spaces []logicalSpace, required int32, fitChar byte) logicalSpace {
	selected := logicalSpace{start: -1}
	for _, space := range spaces {
		if space.size < required {
			continue
		}
		switch fitChar {
		case 'F':
			return space
		case 'B':
			if selected.start == -1 || space.size < selected.size {
				selected = space
			}
		case 'W':
			if selected.start == -1 || space.size > selected.size {
				selected = space
			}
		default:
			if selected.start == -1 {
				selected = space
			}
		}
	}
	return selected
}
