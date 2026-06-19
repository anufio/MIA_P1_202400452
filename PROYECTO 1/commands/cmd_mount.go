package commands

// cmd_mount.go: comandos para montar, desmontar y listar particiones montadas

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
	IsLogical      bool
	EBRStart       int64
}

var mountTable = make(map[string]MountedPartition)

var diskLetterMap = make(map[string]string)
var nextDiskLetterIndex = 0

var logicalNumberMap = make(map[string]int32)
var nextLogicalNumber int32 = 1

const CarnetPrefix = "52"

func CmdMOUNT(params map[string]string) string {
	path, ok := params["path"]
	if !ok {
		return "Error: falta el parámetro obligatorio -path"
	}

	path = strings.TrimSpace(path)

	if path == "" {
		return "Error: -path no puede estar vacío"
	}

	if !disk.DiskExists(path) {
		return fmt.Sprintf("Error: el disco '%s' no existe", path)
	}

	name, ok := params["name"]
	if !ok {
		return "Error: falta el parámetro obligatorio -name"
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return "Error: -name no puede estar vacío"
	}

	for _, mp := range mountTable {
		if mp.DiskPath == path && strings.EqualFold(mp.Name, name) {
			return fmt.Sprintf(
				"Partición '%s' ya está montada con ID: '%s'\n\n%s",
				name,
				mp.ID,
				CmdMOUNTED(),
			)
		}
	}

	selectedPartition, partIdx, isLogical, ebrStart, found, err := findMountablePartition(path, name)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	if !found {
		return fmt.Sprintf("Error: no existe una partición primaria o lógica con el nombre '%s'", name)
	}

	partitionNumber := getPartitionNumberForID(path, name, selectedPartition, partIdx, isLogical)
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

	mp := MountedPartition{
		ID:             id,
		DiskPath:       path,
		Name:           name,
		PartitionIndex: partIdx,
		Partition:      selectedPartition,
		IsLogical:      isLogical,
		EBRStart:       ebrStart,
	}

	if err := writeMountStateToDisk(mp, true); err != nil {
		return fmt.Sprintf("Error al actualizar estado de montaje: %v", err)
	}

	mountTable[id] = mp

	return fmt.Sprintf(
		"Partición '%s' montada exitosamente con ID: '%s'\n\n%s",
		name,
		id,
		CmdMOUNTED(),
	)
}

func CmdUNMOUNT(params map[string]string) string {
	id, ok := params["id"]
	if !ok {
		return "Error: falta el parámetro obligatorio -id"
	}

	id = strings.ToUpper(strings.TrimSpace(id))

	if id == "" {
		return "Error: -id no puede estar vacío"
	}

	mp, ok := mountTable[id]
	if !ok {
		return fmt.Sprintf("Error: el ID '%s' no está montado", id)
	}

	if err := writeMountStateToDisk(mp, false); err != nil {
		return fmt.Sprintf("Error al actualizar estado de desmontaje: %v", err)
	}

	delete(mountTable, id)

	if currentSession.Active && strings.ToUpper(currentSession.PartID) == id {
		currentSession = Session{}
		return fmt.Sprintf(
			"Partición '%s' desmontada exitosamente. La sesión activa fue cerrada porque pertenecía a esa partición.",
			mp.Name,
		)
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

	if mp.IsLogical {
		ebr, err := disk.ReadEBR(mp.DiskPath, mp.EBRStart)
		if err != nil {
			return "", disk.Partition{}, -1, err
		}

		partition := logicalEBRToPartition(ebr, mp.Name)
		partition.PartStatus = [1]byte{'1'}
		partition.PartId = mp.Partition.PartId

		return mp.DiskPath, partition, -1, nil
	}

	mbr, err := disk.ReadMBR(mp.DiskPath)
	if err != nil {
		return "", disk.Partition{}, -1, err
	}

	if mp.PartitionIndex < 0 || mp.PartitionIndex >= len(mbr.MbrPartitions) {
		return "", disk.Partition{}, -1, fmt.Errorf("partición montada con ID '%s' no encontrada en el disco", id)
	}

	p := mbr.MbrPartitions[mp.PartitionIndex]
	pName := cleanPartitionName(p.PartName[:])

	if p.PartStart <= 0 || p.PartType[0] != 'P' || !strings.EqualFold(pName, mp.Name) {
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
		sb.WriteString(fmt.Sprintf(
			" ID: %s | Disco: %s | Partición: %s\n",
			id,
			mp.DiskPath,
			mp.Name,
		))
	}

	return sb.String()
}

func findMountablePartition(path string, name string) (disk.Partition, int, bool, int64, bool, error) {
	mbr, err := disk.ReadMBR(path)
	if err != nil {
		return disk.Partition{}, -1, false, -1, false, fmt.Errorf("error al leer MBR: %v", err)
	}

	for i, p := range mbr.MbrPartitions {
		if p.PartStart <= 0 {
			continue
		}

		pName := cleanPartitionName(p.PartName[:])

		if strings.EqualFold(pName, name) {
			switch p.PartType[0] {
			case 'P':
				return p, i, false, -1, true, nil

			case 'E':
				return disk.Partition{}, -1, false, -1, false,
					fmt.Errorf("la partición '%s' es extendida. Las extendidas no se montan como EXT2; solo contienen particiones lógicas", name)

			default:
				return disk.Partition{}, -1, false, -1, false,
					fmt.Errorf("la partición '%s' tiene un tipo inválido", name)
			}
		}
	}

	if ext, _, foundExt := disk.GetExtended(mbr); foundExt {
		ebr, ebrPos, foundLogical := disk.FindLogicalByName(path, ext, name)
		if foundLogical {
			partition := logicalEBRToPartition(ebr, name)
			return partition, -1, true, ebrPos, true, nil
		}
	}

	return disk.Partition{}, -1, false, -1, false, nil
}

func logicalEBRToPartition(ebr disk.EBR, name string) disk.Partition {
	var nameArr [16]byte
	copy(nameArr[:], name)

	return disk.Partition{
		PartStatus:      ebr.PartMount,
		PartType:        [1]byte{'L'},
		PartFit:         ebr.PartFit,
		PartStart:       ebr.PartStart,
		PartS:           ebr.PartS,
		PartName:        nameArr,
		PartCorrelative: 0,
		PartId:          [4]byte{},
	}
}

func writeMountStateToDisk(mp MountedPartition, mounted bool) error {
	if mp.IsLogical {
		ebr, err := disk.ReadEBR(mp.DiskPath, mp.EBRStart)
		if err != nil {
			return err
		}

		if mounted {
			ebr.PartMount = [1]byte{'1'}
		} else {
			ebr.PartMount = [1]byte{'0'}
		}

		return disk.WriteEBR(mp.DiskPath, mp.EBRStart, ebr)
	}

	mbr, err := disk.ReadMBR(mp.DiskPath)
	if err != nil {
		return err
	}

	if mp.PartitionIndex < 0 || mp.PartitionIndex >= len(mbr.MbrPartitions) {
		return fmt.Errorf("índice de partición inválido")
	}

	if mounted {
		mbr.MbrPartitions[mp.PartitionIndex].PartStatus = [1]byte{'1'}
		mbr.MbrPartitions[mp.PartitionIndex].PartId = mp.Partition.PartId
	} else {
		mbr.MbrPartitions[mp.PartitionIndex].PartStatus = [1]byte{'0'}
		mbr.MbrPartitions[mp.PartitionIndex].PartId = [4]byte{}
	}

	return disk.WriteMBR(mp.DiskPath, mbr)
}

func getPartitionNumberForID(path string, name string, partition disk.Partition, index int, isLogical bool) int32 {
	if !isLogical {
		if partition.PartCorrelative > 0 {
			return partition.PartCorrelative
		}

		return int32(index + 1)
	}

	key := strings.ToLower(strings.TrimSpace(path + "::" + name))

	if n, ok := logicalNumberMap[key]; ok {
		return n
	}

	n := nextLogicalNumber
	logicalNumberMap[key] = n
	nextLogicalNumber++

	return n
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

func cleanPartitionName(raw []byte) string {
	return strings.TrimSpace(strings.TrimRight(string(raw), "\x00"))
}
