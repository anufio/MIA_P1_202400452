package commands

// cmd_mount.go

import (
	"MIA_P1_202400452/disk"
	"fmt"
	"sort"
	"strings"
)

type MountedPartition struct {
	ID             string
	DiskPath       string
	Name           string
	PartitionIndex int
	Partition      disk.Partition
}

var mountTable = make(map[string]MountedPartition)

var diskLetterMap = make(map[string]string)

var nextDiskLetterIndex = 0

const CarnetPrefix = "52"

func CmdMOUNT(params map[string]string) string {
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

	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return fmt.Sprintf("Error al leer MBR: %v", err)
	}

	partIdx := -1
	var selectedPartition disk.Partition

	for i, p := range mbr.MbrPartitions {
		pName := strings.TrimRight(string(p.PartName[:]), "\x00")

		if pName == name && p.PartStart > 0 {
			if p.PartType[0] != 'P' {
				return fmt.Sprintf("Error: la partición '%s' no es primaria. El manejo de archivos EXT2 se realiza sobre particiones primarias", name)
			}

			partIdx = i
			selectedPartition = p
			break
		}
	}

	if partIdx == -1 {
		return fmt.Sprintf("Error: no existe una partición primaria con el nombre '%s'", name)
	}

	for _, mp := range mountTable {
		if mp.DiskPath == path && mp.Name == name {
			return fmt.Sprintf("Partición '%s' ya está montada con ID: '%s'", name, mp.ID)
		}
	}

	partitionNumber := getPartitionNumber(selectedPartition, partIdx)
	diskLetter := getDiskLetter(path)

	id := fmt.Sprintf("%s%d%s", CarnetPrefix, partitionNumber, diskLetter)
	id = strings.ToUpper(id)

	if _, exists := mountTable[id]; exists {
		return fmt.Sprintf("Error: ya existe una partición montada con el ID '%s'", id)
	}

	selectedPartition.PartStatus = [1]byte{'1'}

	var idArr [4]byte
	copy(idArr[:], id)
	selectedPartition.PartId = idArr

	mountTable[id] = MountedPartition{
		ID:             id,
		DiskPath:       path,
		Name:           name,
		PartitionIndex: partIdx,
		Partition:      selectedPartition,
	}

	return fmt.Sprintf("Partición '%s' montada exitosamente con ID: '%s'", name, id)
}

func CmdUNMOUNT(params map[string]string) string {
	id, ok := params["id"]
	if !ok {
		return "Error: falta el parámetro obligatorio -id"
	}

	id = strings.ToUpper(strings.TrimSpace(id))

	mp, ok := mountTable[id]
	if !ok {
		return fmt.Sprintf("Error: el ID '%s' no está montado", id)
	}

	delete(mountTable, id)

	if currentSession.Active && strings.ToUpper(currentSession.PartID) == id {
		currentSession = Session{}
		return fmt.Sprintf("Partición '%s' desmontada exitosamente. La sesión activa fue cerrada porque pertenecía a esa partición.", mp.Name)
	}

	return fmt.Sprintf("Partición '%s' desmontada exitosamente.", mp.Name)
}

func GetMountedPartition(id string) (MountedPartition, bool) {
	id = strings.ToUpper(strings.TrimSpace(id))
	mp, ok := mountTable[id]
	return mp, ok
}

func GetMountedPartitionInfo(id string) (string, disk.Partition, int, error) {
	id = strings.ToUpper(strings.TrimSpace(id))

	mp, ok := mountTable[id]
	if !ok {
		return "", disk.Partition{}, -1, fmt.Errorf("el ID '%s' no está montado", id)
	}

	mbr, err := disk.ReadMBR(mp.DiskPath)
	if err != nil {
		return "", disk.Partition{}, -1, err
	}

	if mp.PartitionIndex < 0 || mp.PartitionIndex >= len(mbr.MbrPartitions) {
		return "", disk.Partition{}, -1, fmt.Errorf("partición montada con ID '%s' no encontrada en el disco", id)
	}

	p := mbr.MbrPartitions[mp.PartitionIndex]
	pName := strings.TrimRight(string(p.PartName[:]), "\x00")

	if p.PartStart <= 0 || p.PartType[0] != 'P' || pName != mp.Name {
		return "", disk.Partition{}, -1, fmt.Errorf("partición montada con ID '%s' no encontrada en el disco", id)
	}

	p.PartStatus = mp.Partition.PartStatus
	p.PartId = mp.Partition.PartId
	p.PartCorrelative = mp.Partition.PartCorrelative

	return mp.DiskPath, p, mp.PartitionIndex, nil
}

func CmdMOUNTED() string {
	if len(mountTable) == 0 {
		return "No hay particiones montadas."
	}

	var ids []string
	for id := range mountTable {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	var sb strings.Builder
	sb.WriteString("Particiones montadas:\n")

	for _, id := range ids {
		mp := mountTable[id]
		sb.WriteString(fmt.Sprintf(" ID: %s | Disco: %s | Partición: %s\n", id, mp.DiskPath, mp.Name))
	}

	return sb.String()
}

func getPartitionNumber(partition disk.Partition, index int) int32 {
	if partition.PartCorrelative > 0 {
		return partition.PartCorrelative
	}

	return int32(index + 1)
}

func getDiskLetter(path string) string {
	if letter, ok := diskLetterMap[path]; ok {
		return letter
	}

	letter := string(rune('A' + nextDiskLetterIndex))
	diskLetterMap[path] = letter
	nextDiskLetterIndex++

	return letter
}
