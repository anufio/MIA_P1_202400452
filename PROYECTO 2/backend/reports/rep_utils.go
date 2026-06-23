package reports

// rep_utils.go

import (
	"fmt"
	"strings"

	"MIA_P1_202400452/disk"
	"MIA_P1_202400452/ext2"
)

func readEXT2SuperBlock(diskPath string, partition disk.Partition) (ext2.SuperBlock, error) {
	sb, err := ext2.ReadSuperBlock(diskPath, partition.PartStart)
	if err != nil {
		return sb, fmt.Errorf("error al leer superbloque: %v", err)
	}
	if sb.SMagic != 0xEF53 {
		return sb, fmt.Errorf("la partición no tiene un sistema de archivos EXT2 formateado")
	}
	return sb, nil
}

func readUserGroupNames(diskPath string, sb ext2.SuperBlock) (map[int32]string, map[int32]string) {
	users := make(map[int32]string)
	groups := make(map[int32]string)

	content, err := ext2.GetFileContent(diskPath, sb, 1)
	if err != nil {
		return users, groups
	}

	for _, line := range strings.Split(content, "\n") {
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}
		id := parseReportID(parts[0])
		if id == 0 {
			continue
		}
		recordType := strings.TrimSpace(parts[1])
		switch recordType {
		case "G":
			groups[id] = strings.TrimSpace(parts[2])
		case "U":
			if len(parts) >= 5 {
				users[id] = strings.TrimSpace(parts[3])
			}
		}
	}
	return users, groups
}

func parseReportID(value string) int32 {
	var id int32
	fmt.Sscanf(strings.TrimSpace(value), "%d", &id)
	return id
}

func nameOrID(names map[int32]string, id int32) string {
	if name, ok := names[id]; ok {
		return name
	}
	return fmt.Sprintf("%d", id)
}
