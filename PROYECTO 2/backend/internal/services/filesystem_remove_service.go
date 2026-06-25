package services

//filesystem_remove_service.go: contiene la lógica para manejar la eliminación de archivos y carpetas dentro del sistema de archivos, incluyendo la validación de rutas, la eliminación recursiva y la gestión de permisos.

import (
	"fmt"

	"MIA_P2_202400452/internal/ext2"
)

func (s *FileSystemService) Remove(input FSRemoveInput) error {
	mounted, sb, session, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return err
	}

	if err := requireSession(session); err != nil {
		return err
	}

	targetPath := cleanFSPath(input.Path)

	if targetPath == "/" {
		return fmt.Errorf("no se puede eliminar la raíz")
	}

	inodeIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, targetPath)
	if err != nil {
		return err
	}

	parentIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, parentPath(targetPath))
	if err != nil {
		return err
	}

	parentInode, err := ext2.ReadInode(mounted.DiskPath, sb, parentIdx)
	if err != nil {
		return err
	}

	if err := s.requireWritePermission(parentInode, session, "modificar la carpeta padre"); err != nil {
		return err
	}

	if err := s.validateRemovePermissionsRecursive(mounted.DiskPath, sb, inodeIdx, session); err != nil {
		return err
	}

	if err := s.removeInodeRecursive(mounted.DiskPath, &sb, inodeIdx); err != nil {
		return err
	}

	if err := ext2.RemoveEntryFromFolder(mounted.DiskPath, &sb, parentIdx, baseName(targetPath)); err != nil {
		return err
	}

	return ext2.WriteSuperBlock(mounted.DiskPath, mounted.Start, sb)
}

func (s *FileSystemService) removeInodeRecursive(diskPath string, sb *ext2.SuperBlock, inodeIdx int32) error {
	inode, err := ext2.ReadInode(diskPath, *sb, inodeIdx)
	if err != nil {
		return err
	}

	if inode.IType[0] == '0' {
		for i := 0; i < 12; i++ {
			if inode.IBlock[i] == -1 {
				continue
			}

			block, err := ext2.ReadFolderBlock(diskPath, *sb, inode.IBlock[i])
			if err != nil {
				continue
			}

			for _, content := range block.BContent {
				name := trimBlockName(content.BName[:])

				if name == "" || name == "." || name == ".." || content.BInode == -1 {
					continue
				}

				if err := s.removeInodeRecursive(diskPath, sb, content.BInode); err != nil {
					return err
				}
			}
		}
	}

	if err := ext2.FreeInodeBlocks(diskPath, sb, &inode); err != nil {
		return err
	}

	if err := ext2.FreeInode(diskPath, sb, inodeIdx); err != nil {
		return err
	}

	return nil
}
