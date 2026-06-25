package services

//filesystem_copy_move_service.go: contiene la lógica para manejar la copia y el movimiento de archivos y carpetas dentro del sistema de archivos, incluyendo la validación de rutas, la creación de directorios y la gestión de permisos.

import (
	"fmt"

	"MIA_P2_202400452/internal/ext2"
)

func (s *FileSystemService) Copy(input FSCopyInput) (FSItem, error) {
	mounted, sb, session, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return FSItem{}, err
	}

	fromPath := cleanFSPath(input.From)
	toPath := cleanFSPath(input.To)

	if fromPath == "/" {
		return FSItem{}, fmt.Errorf("no se puede copiar la raíz")
	}

	if _, err := ext2.FindInodeByPath(mounted.DiskPath, sb, toPath); err == nil {
		return FSItem{}, fmt.Errorf("ya existe la ruta destino")
	}

	if err := s.copyRecursive(mounted.DiskPath, &sb, fromPath, toPath, session); err != nil {
		return FSItem{}, err
	}

	if err := ext2.WriteSuperBlock(mounted.DiskPath, mounted.Start, sb); err != nil {
		return FSItem{}, err
	}

	inodeIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, toPath)
	if err != nil {
		return FSItem{}, err
	}

	inode, err := ext2.ReadInode(mounted.DiskPath, sb, inodeIdx)
	if err != nil {
		return FSItem{}, err
	}

	return FSItem{
		Name:        baseName(toPath),
		Path:        toPath,
		Type:        inodeTypeLabel(inode),
		Inode:       inodeIdx,
		Size:        inode.IS,
		Permissions: inodePermissions(inode),
		UID:         inode.IUid,
		GID:         inode.IGid,
		Modified:    inodeTime(inode.IMtime),
	}, nil
}

func (s *FileSystemService) Move(input FSMoveInput) (FSItem, error) {
	item, err := s.Copy(FSCopyInput{
		Token: input.Token,
		ID:    input.ID,
		From:  input.From,
		To:    input.To,
	})
	if err != nil {
		return FSItem{}, err
	}

	if err := s.Remove(FSRemoveInput{
		Token: input.Token,
		ID:    input.ID,
		Path:  input.From,
	}); err != nil {
		return FSItem{}, err
	}

	return item, nil
}

func (s *FileSystemService) copyRecursive(diskPath string, sb *ext2.SuperBlock, fromPath string, toPath string, session SessionInfo) error {
	sourceIdx, err := ext2.FindInodeByPath(diskPath, *sb, fromPath)
	if err != nil {
		return err
	}

	sourceInode, err := ext2.ReadInode(diskPath, *sb, sourceIdx)
	if err != nil {
		return err
	}

	if sourceInode.IType[0] == '1' {
		content, err := ext2.GetFileContent(diskPath, *sb, sourceIdx)
		if err != nil {
			return err
		}

		item, err := s.Mkfile(FSMkfileInput{
			Token:   session.Token,
			ID:      session.PartID,
			Path:    toPath,
			Content: content,
			Parents: true,
		})
		_ = item

		return err
	}

	if err := s.ensureDirectory(diskPath, sb, toPath, session); err != nil {
		return err
	}

	for i := 0; i < 12; i++ {
		if sourceInode.IBlock[i] == -1 {
			continue
		}

		block, err := ext2.ReadFolderBlock(diskPath, *sb, sourceInode.IBlock[i])
		if err != nil {
			continue
		}

		for _, content := range block.BContent {
			name := trimBlockName(content.BName[:])

			if name == "" || name == "." || name == ".." || content.BInode == -1 {
				continue
			}

			if err := s.copyRecursive(diskPath, sb, joinFSPath(fromPath, name), joinFSPath(toPath, name), session); err != nil {
				return err
			}
		}
	}

	return nil
}
