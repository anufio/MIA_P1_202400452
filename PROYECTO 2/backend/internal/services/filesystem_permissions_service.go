package services

//filesystem_permissions_service.go: contiene las funciones para validar los permisos de lectura y escritura en el sistema de archivos, incluyendo la verificación de permisos para eliminar, copiar y mover archivos y carpetas, así como la validación de permisos recursivos en directorios.

import (
	"fmt"

	"MIA_P2_202400452/internal/ext2"
)

func requireSession(session SessionInfo) error {
	if !session.Active || session.Token == "" {
		return fmt.Errorf("debe iniciar sesión para ejecutar este comando")
	}

	return nil
}

func (s *FileSystemService) requireWritePermission(inode ext2.Inode, session SessionInfo, action string) error {
	if err := requireSession(session); err != nil {
		return err
	}

	if !ext2.HasWritePermission(inode, session.UID, session.GID) {
		return fmt.Errorf("no tiene permiso de escritura para %s", action)
	}

	return nil
}

func (s *FileSystemService) requireReadPermission(inode ext2.Inode, session SessionInfo, action string) error {
	if err := requireSession(session); err != nil {
		return err
	}

	if !ext2.HasReadPermission(inode, session.UID, session.GID) {
		return fmt.Errorf("no tiene permiso de lectura para %s", action)
	}

	return nil
}

func (s *FileSystemService) requireReadWritePermission(inode ext2.Inode, session SessionInfo, action string) error {
	if err := s.requireReadPermission(inode, session, action); err != nil {
		return err
	}

	if err := s.requireWritePermission(inode, session, action); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemService) validateRemovePermissionsRecursive(diskPath string, sb ext2.SuperBlock, inodeIdx int32, session SessionInfo) error {
	inode, err := ext2.ReadInode(diskPath, sb, inodeIdx)
	if err != nil {
		return err
	}

	if err := s.requireWritePermission(inode, session, "eliminar esta entrada"); err != nil {
		return err
	}

	if inode.IType[0] != '0' {
		return nil
	}

	for i := 0; i < 12; i++ {
		if inode.IBlock[i] == -1 {
			continue
		}

		block, err := ext2.ReadFolderBlock(diskPath, sb, inode.IBlock[i])
		if err != nil {
			continue
		}

		for _, content := range block.BContent {
			name := trimBlockName(content.BName[:])

			if name == "" || name == "." || name == ".." || content.BInode == -1 {
				continue
			}

			if err := s.validateRemovePermissionsRecursive(diskPath, sb, content.BInode, session); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *FileSystemService) validateCopyPermissionsRecursive(diskPath string, sb ext2.SuperBlock, inodeIdx int32, session SessionInfo) error {
	inode, err := ext2.ReadInode(diskPath, sb, inodeIdx)
	if err != nil {
		return err
	}

	if err := s.requireReadPermission(inode, session, "copiar esta entrada"); err != nil {
		return err
	}

	if inode.IType[0] != '0' {
		return nil
	}

	for i := 0; i < 12; i++ {
		if inode.IBlock[i] == -1 {
			continue
		}

		block, err := ext2.ReadFolderBlock(diskPath, sb, inode.IBlock[i])
		if err != nil {
			continue
		}

		for _, content := range block.BContent {
			name := trimBlockName(content.BName[:])

			if name == "" || name == "." || name == ".." || content.BInode == -1 {
				continue
			}

			if err := s.validateCopyPermissionsRecursive(diskPath, sb, content.BInode, session); err != nil {
				return err
			}
		}
	}

	return nil
}
