package commands

//cmd_cat.go: comando para mostrar el contenido de archivos

import (
	"MIA_P1_202400452/ext2"
	"fmt"
	"strings"
)

func CmdCAT(params map[string]string) string {
	if err := RequireSession(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	sess := GetSession()
	diskPath, partition, _, err := GetMountedPartitionInfo(sess.PartID)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	sb, err := ext2.ReadSuperBlock(diskPath, partition.PartStart)
	if err != nil {
		return fmt.Sprintf("Error al leer superbloque: %v", err)
	}
	if sb.SMagic != 0xEF53 {
		return "Error: la partición no tiene un sistema de archivos EXT2 formateado"
	}

	var contents []string
	for i := 1; ; i++ {
		key := fmt.Sprintf("file%d", i)
		filePath, ok := params[key]
		if !ok {
			break
		}

		inodeIdx, err := ext2.FindInodeByPath(diskPath, sb, filePath)
		if err != nil {
			return fmt.Sprintf("Error: el archivo '%s' no existe: %v", filePath, err)
		}

		inode, err := ext2.ReadInode(diskPath, sb, inodeIdx)
		if err != nil {
			return fmt.Sprintf("Error al leer inodo de '%s': %v", filePath, err)
		}

		if inode.IType[0] != '1' {
			return fmt.Sprintf("Error: '%s' no es un archivo", filePath)
		}

		if !ext2.HasReadPermission(inode, sess.UID, sess.GID) {
			return fmt.Sprintf("Error: no tiene permiso de lectura en '%s'", filePath)
		}

		content, err := ext2.GetFileContent(diskPath, sb, inodeIdx)
		if err != nil {
			return fmt.Sprintf("Error al leer '%s': %v", filePath, err)
		}
		contents = append(contents, content)
	}

	if len(contents) == 0 {
		return "Error: debe especificar al menos un archivo con -file1=<ruta>"
	}

	return strings.Join(contents, "\n")
}
