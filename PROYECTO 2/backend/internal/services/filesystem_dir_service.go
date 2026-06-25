package services

//filesystem_dir_service.go: contiene la lógica para manejar la creación de directorios dentro del sistema de archivos, incluyendo la validación de rutas, la creación de directorios y la gestión de permisos.

import (
	"fmt"
	"strings"

	"MIA_P2_202400452/internal/ext2"
)

func (s *FileSystemService) ensureDirectory(diskPath string, sb *ext2.SuperBlock, targetPath string, session SessionInfo) error {
	targetPath = cleanFSPath(targetPath)

	if targetPath == "/" {
		return nil
	}

	parts := strings.Split(strings.Trim(targetPath, "/"), "")
	_ = parts

	current := ""

	for _, part := range strings.Split(strings.Trim(targetPath, "/"), "/") {
		if part == "" {
			continue
		}

		current = cleanFSPath(current + "/" + part)

		if _, err := ext2.FindInodeByPath(diskPath, *sb, current); err == nil {
			continue
		}

		if err := s.createDirectory(diskPath, sb, parentPath(current), baseName(current), session); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileSystemService) createDirectory(diskPath string, sb *ext2.SuperBlock, parent string, name string, session SessionInfo) error {
	parent = cleanFSPath(parent)
	name = strings.TrimSpace(name)

	if name == "" || name == "." || name == "/" {
		return fmt.Errorf("nombre de carpeta inválido")
	}

	if len(name) > 12 {
		return fmt.Errorf("el nombre no puede exceder 12 caracteres")
	}

	targetPath := joinFSPath(parent, name)

	if _, err := ext2.FindInodeByPath(diskPath, *sb, targetPath); err == nil {
		return fmt.Errorf("ya existe '%s'", targetPath)
	}

	parentIdx, err := ext2.FindInodeByPath(diskPath, *sb, parent)
	if err != nil {
		return err
	}

	parentInode, err := ext2.ReadInode(diskPath, *sb, parentIdx)
	if err != nil {
		return err
	}

	if parentInode.IType[0] != '0' {
		return fmt.Errorf("la ruta padre no es una carpeta")
	}

	inodeIdx, err := ext2.AllocateInode(diskPath, sb)
	if err != nil {
		return err
	}

	blockIdx, err := ext2.AllocateBlock(diskPath, sb)
	if err != nil {
		return err
	}

	inode := ext2.NewInode()
	inode.IType = [1]byte{'0'}
	inode.IPerm = [3]byte{'6', '6', '4'}
	inode.IUid = session.UID
	inode.IGid = session.GID
	inode.ICtime = ext2.GetCurrentTime()
	inode.IAtime = ext2.GetCurrentTime()
	inode.IMtime = ext2.GetCurrentTime()
	inode.IBlock[0] = blockIdx

	if inode.IUid == 0 {
		inode.IUid = 1
	}

	if inode.IGid == 0 {
		inode.IGid = 1
	}

	folder := ext2.NewFolderBlock()
	copy(folder.BContent[0].BName[:], ".")
	folder.BContent[0].BInode = inodeIdx
	copy(folder.BContent[1].BName[:], "..")
	folder.BContent[1].BInode = parentIdx

	if err := ext2.WriteFolderBlock(diskPath, *sb, blockIdx, folder); err != nil {
		return err
	}

	if err := ext2.WriteInode(diskPath, *sb, inodeIdx, inode); err != nil {
		return err
	}

	if err := ext2.AddEntryToFolder(diskPath, sb, parentIdx, name, inodeIdx); err != nil {
		return err
	}

	return nil
}
