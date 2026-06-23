package commands

// cmd_files.go: comandos para crear archivos y directorios en un sistema de archivos EXT2

import (
	"MIA_P1_202400452/ext2"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CmdMKFILE(params map[string]string) string {
	if err := RequireSession(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	filePath, ok := params["path"]
	if !ok {
		return "Error: falta el parámetro obligatorio -path"
	}
	_, hasR := params["r"]
	if hasR && params["r"] != "" {
		return "Error: el parámetro -r no recibe ningún valor"
	}
	if strings.TrimSpace(filePath) == "" || filePath == "/" {
		return "Error: la ruta del archivo no es válida"
	}
	size := 0
	if sizeStr, ok := params["size"]; ok {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil || size < 0 {
			return "Error: -size debe ser un número no negativo"
		}
	}
	contPath := ""
	if cp, ok := params["cont"]; ok {
		contPath = cp
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

	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	if len(base) > 12 {
		return "Error: el nombre del archivo no puede exceder 12 caracteres"
	}

	parentInodeIdx, err := ext2.FindInodeByPath(diskPath, sb, dir)
	if err != nil {
		if !hasR {
			return fmt.Sprintf("Error: el directorio padre '%s' no existe. Use -r para crearlo", dir)
		}

		if err := createDirsRecursive(diskPath, &sb, dir, sess.UID, sess.GID, partition.PartStart); err != nil {
			return fmt.Sprintf("Error al crear directorios padre: %v", err)
		}
		parentInodeIdx, err = ext2.FindInodeByPath(diskPath, sb, dir)
		if err != nil {
			return fmt.Sprintf("Error: no se pudo crear el directorio padre: %v", err)
		}
	}

	parentInode, _ := ext2.ReadInode(diskPath, sb, parentInodeIdx)
	if !ext2.HasWritePermission(parentInode, sess.UID, sess.GID) {
		return fmt.Sprintf("Error: no tiene permiso de escritura en '%s'", dir)
	}

	existing, err := ext2.FindInodeByPath(diskPath, sb, filePath)
	if err == nil && existing != -1 {
		existingInode, err := ext2.ReadInode(diskPath, sb, existing)
		if err != nil {
			return fmt.Sprintf("Error al leer archivo existente: %v", err)
		}
		if existingInode.IType[0] != '1' {
			return fmt.Sprintf("Error: '%s' ya existe y no es un archivo", filePath)
		}
		fmt.Printf("El archivo '%s' ya existe. ¿Desea sobreescribirlo? (s/n): ", filePath)
		var resp string
		fmt.Scanln(&resp)
		if strings.ToLower(strings.TrimSpace(resp)) != "s" {
			return "Operación cancelada."
		}

		if err := deleteFile(diskPath, &sb, parentInodeIdx, base, existing); err != nil {
			return fmt.Sprintf("Error al eliminar archivo existente: %v", err)
		}
	}

	var content []byte
	if contPath != "" {
		data, err := os.ReadFile(contPath)
		if err != nil {
			return fmt.Sprintf("Error al leer archivo fuente '%s': %v", contPath, err)
		}
		content = data
	} else if size > 0 {
		content = make([]byte, size)
		for i := range content {
			content[i] = byte('0' + (i % 10))
		}
	}

	newInodeIdx, err := ext2.AllocateInode(diskPath, &sb)
	if err != nil {
		return fmt.Sprintf("Error al asignar inodo: %v", err)
	}
	newInode := ext2.NewInode()
	newInode.IUid = sess.UID
	newInode.IGid = sess.GID
	newInode.IAtime = ext2.GetCurrentTime()
	newInode.ICtime = ext2.GetCurrentTime()
	newInode.IMtime = ext2.GetCurrentTime()
	newInode.IType = [1]byte{'1'}
	copy(newInode.IPerm[:], "664")

	if len(content) > 0 {
		if err := ext2.WriteFileContent(diskPath, &sb, &newInode, content); err != nil {
			return fmt.Sprintf("Error al escribir contenido: %v", err)
		}
	}
	if err := ext2.WriteInode(diskPath, sb, newInodeIdx, newInode); err != nil {
		return fmt.Sprintf("Error al escribir inodo: %v", err)
	}

	if err := ext2.AddEntryToFolder(diskPath, &sb, parentInodeIdx, base, newInodeIdx); err != nil {
		return fmt.Sprintf("Error al agregar entrada al directorio: %v", err)
	}

	ext2.WriteSuperBlock(diskPath, partition.PartStart, sb)

	return fmt.Sprintf("Archivo '%s' creado exitosamente.", filePath)
}

func deleteFile(diskPath string, sb *ext2.SuperBlock, parentInodeIdx int32, name string, inodeIdx int32) error {

	inode, err := ext2.ReadInode(diskPath, *sb, inodeIdx)
	if err != nil {
		return err
	}
	if inode.IType[0] != '1' {
		return fmt.Errorf("la ruta existente no es un archivo")
	}
	if err := ext2.FreeInodeBlocks(diskPath, sb, &inode); err != nil {
		return err
	}
	if err := ext2.WriteInode(diskPath, *sb, inodeIdx, inode); err != nil {
		return err
	}
	if err := ext2.RemoveEntryFromFolder(diskPath, sb, parentInodeIdx, name); err != nil {
		return err
	}
	return ext2.FreeInode(diskPath, sb, inodeIdx)
}

func CmdMKDIR(params map[string]string) string {
	if err := RequireSession(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	dirPath, ok := params["path"]
	if !ok {
		return "Error: falta el parámetro obligatorio -path"
	}
	_, hasP := params["p"]
	if hasP && params["p"] != "" {
		return "Error: el parámetro -p no recibe ningún valor"
	}
	if strings.TrimSpace(dirPath) == "" || dirPath == "/" {
		return "Error: la ruta del directorio no es válida"
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

	if hasP {
		if err := createDirsRecursive(diskPath, &sb, dirPath, sess.UID, sess.GID, partition.PartStart); err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		ext2.WriteSuperBlock(diskPath, partition.PartStart, sb)
		return fmt.Sprintf("Directorio '%s' creado exitosamente.", dirPath)
	}

	dir := filepath.Dir(dirPath)
	base := filepath.Base(dirPath)
	if len(base) > 12 {
		return "Error: el nombre del directorio no puede exceder 12 caracteres"
	}
	parentInodeIdx, err := ext2.FindInodeByPath(diskPath, sb, dir)
	if err != nil {
		return fmt.Sprintf("Error: el directorio padre '%s' no existe. Use -p para crearlo", dir)
	}
	parentInode, _ := ext2.ReadInode(diskPath, sb, parentInodeIdx)
	if !ext2.HasWritePermission(parentInode, sess.UID, sess.GID) {
		return fmt.Sprintf("Error: no tiene permiso de escritura en '%s'", dir)
	}

	if _, err := ext2.FindInodeByPath(diskPath, sb, dirPath); err == nil {
		return fmt.Sprintf("Error: el directorio '%s' ya existe", dirPath)
	}

	if err := createSingleDir(diskPath, &sb, parentInodeIdx, base, sess.UID, sess.GID); err != nil {
		return fmt.Sprintf("Error al crear directorio: %v", err)
	}
	ext2.WriteSuperBlock(diskPath, partition.PartStart, sb)
	return fmt.Sprintf("Directorio '%s' creado exitosamente.", dirPath)
}

func createDirsRecursive(diskPath string, sb *ext2.SuperBlock, path string, uid, gid int32, partStart int32) error {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentPath := ""
	currentInodeIdx := int32(0)

	for _, part := range parts {
		if part == "" {
			continue
		}
		if len(part) > 12 {
			return fmt.Errorf("el nombre '%s' no puede exceder 12 caracteres", part)
		}
		currentPath += "/" + part

		idx, err := ext2.FindInodeByPath(diskPath, *sb, currentPath)
		if err == nil && idx != -1 {
			inode, err := ext2.ReadInode(diskPath, *sb, idx)
			if err != nil {
				return err
			}
			if inode.IType[0] != '0' {
				return fmt.Errorf("'%s' ya existe y no es un directorio", currentPath)
			}
			currentInodeIdx = idx
			continue
		}

		parentInode, err := ext2.ReadInode(diskPath, *sb, currentInodeIdx)
		if err != nil {
			return err
		}
		if !ext2.HasWritePermission(parentInode, uid, gid) {
			return fmt.Errorf("no tiene permiso de escritura en el directorio padre de '%s'", currentPath)
		}

		if err := createSingleDir(diskPath, sb, currentInodeIdx, part, uid, gid); err != nil {
			return fmt.Errorf("error creando '%s': %v", currentPath, err)
		}
		newIdx, err := ext2.FindInodeByPath(diskPath, *sb, currentPath)
		if err != nil {
			return fmt.Errorf("error verificando '%s': %v", currentPath, err)
		}
		currentInodeIdx = newIdx
	}
	return nil
}

func createSingleDir(diskPath string, sb *ext2.SuperBlock, parentInodeIdx int32, name string, uid, gid int32) error {

	newInodeIdx, err := ext2.AllocateInode(diskPath, sb)
	if err != nil {
		return err
	}

	newBlockIdx, err := ext2.AllocateBlock(diskPath, sb)
	if err != nil {
		return err
	}

	fb := ext2.NewFolderBlock()
	copy(fb.BContent[0].BName[:], ".")
	fb.BContent[0].BInode = newInodeIdx
	copy(fb.BContent[1].BName[:], "..")
	fb.BContent[1].BInode = parentInodeIdx

	if err := ext2.WriteFolderBlock(diskPath, *sb, newBlockIdx, fb); err != nil {
		return err
	}

	newInode := ext2.NewInode()
	newInode.IUid = uid
	newInode.IGid = gid
	newInode.IAtime = ext2.GetCurrentTime()
	newInode.ICtime = ext2.GetCurrentTime()
	newInode.IMtime = ext2.GetCurrentTime()
	newInode.IType = [1]byte{'0'}
	copy(newInode.IPerm[:], "664")
	newInode.IBlock[0] = newBlockIdx

	if err := ext2.WriteInode(diskPath, *sb, newInodeIdx, newInode); err != nil {
		return err
	}

	return ext2.AddEntryToFolder(diskPath, sb, parentInodeIdx, name, newInodeIdx)
}
