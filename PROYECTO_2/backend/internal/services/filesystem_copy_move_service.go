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

	if err := requireSession(session); err != nil {
		return FSItem{}, err
	}

	fromPath := cleanFSPath(input.From)
	toPath := cleanFSPath(input.To)

	if fromPath == "/" {
		return FSItem{}, fmt.Errorf("no se puede copiar la raíz")
	}

	sourceIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, fromPath)
	if err != nil {
		return FSItem{}, err
	}

	if err := s.validateCopyPermissionsRecursive(mounted.DiskPath, sb, sourceIdx, session); err != nil {
		return FSItem{}, err
	}

	destinationParentPath := parentPath(toPath)

	destinationParentIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, destinationParentPath)
	if err != nil {
		return FSItem{}, fmt.Errorf("la carpeta destino no existe")
	}

	destinationParentInode, err := ext2.ReadInode(mounted.DiskPath, sb, destinationParentIdx)
	if err != nil {
		return FSItem{}, err
	}

	if destinationParentInode.IType[0] != '0' {
		return FSItem{}, fmt.Errorf("el destino debe estar dentro de una carpeta")
	}

	if err := s.requireWritePermission(destinationParentInode, session, "escribir en la carpeta destino"); err != nil {
		return FSItem{}, err
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
	mounted, sb, session, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return FSItem{}, err
	}

	if err := requireSession(session); err != nil {
		return FSItem{}, err
	}

	fromPath := cleanFSPath(input.From)
	toPath := cleanFSPath(input.To)

	if fromPath == "/" {
		return FSItem{}, fmt.Errorf("no se puede mover la raíz")
	}

	if fromPath == toPath {
		return FSItem{}, fmt.Errorf("la ruta origen y destino son iguales")
	}

	if len(baseName(toPath)) > 12 || baseName(toPath) == "" || baseName(toPath) == "." || baseName(toPath) == ".." {
		return FSItem{}, fmt.Errorf("nombre destino inválido")
	}

	if len(toPath) > len(fromPath) && toPath[:len(fromPath)] == fromPath && toPath[len(fromPath)] == '/' {
		return FSItem{}, fmt.Errorf("no se puede mover una carpeta dentro de sí misma")
	}

	sourceIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, fromPath)
	if err != nil {
		return FSItem{}, err
	}

	sourceInode, err := ext2.ReadInode(mounted.DiskPath, sb, sourceIdx)
	if err != nil {
		return FSItem{}, err
	}

	if err := s.requireWritePermission(sourceInode, session, "mover esta entrada"); err != nil {
		return FSItem{}, err
	}

	sourceParentIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, parentPath(fromPath))
	if err != nil {
		return FSItem{}, err
	}

	sourceParentInode, err := ext2.ReadInode(mounted.DiskPath, sb, sourceParentIdx)
	if err != nil {
		return FSItem{}, err
	}

	if err := s.requireWritePermission(sourceParentInode, session, "modificar la carpeta origen"); err != nil {
		return FSItem{}, err
	}

	destinationParentIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, parentPath(toPath))
	if err != nil {
		return FSItem{}, fmt.Errorf("la carpeta destino no existe")
	}

	destinationParentInode, err := ext2.ReadInode(mounted.DiskPath, sb, destinationParentIdx)
	if err != nil {
		return FSItem{}, err
	}

	if destinationParentInode.IType[0] != '0' {
		return FSItem{}, fmt.Errorf("el destino debe estar dentro de una carpeta")
	}

	if err := s.requireWritePermission(destinationParentInode, session, "modificar la carpeta destino"); err != nil {
		return FSItem{}, err
	}

	if _, err := ext2.FindInodeByPath(mounted.DiskPath, sb, toPath); err == nil {
		return FSItem{}, fmt.Errorf("ya existe la ruta destino")
	}

	if err := ext2.RemoveEntryFromFolder(mounted.DiskPath, &sb, sourceParentIdx, baseName(fromPath)); err != nil {
		return FSItem{}, err
	}

	if err := ext2.AddEntryToFolder(mounted.DiskPath, &sb, destinationParentIdx, baseName(toPath), sourceIdx); err != nil {
		// Intento de reversión para no perder la referencia si falla el destino.
		_ = ext2.AddEntryToFolder(mounted.DiskPath, &sb, sourceParentIdx, baseName(fromPath), sourceIdx)
		return FSItem{}, err
	}

	sourceInode.IMtime = ext2.GetCurrentTime()
	if err := ext2.WriteInode(mounted.DiskPath, sb, sourceIdx, sourceInode); err != nil {
		return FSItem{}, err
	}

	if err := ext2.WriteSuperBlock(mounted.DiskPath, mounted.Start, sb); err != nil {
		return FSItem{}, err
	}

	return FSItem{
		Name:        baseName(toPath),
		Path:        toPath,
		Type:        inodeTypeLabel(sourceInode),
		Inode:       sourceIdx,
		Size:        sourceInode.IS,
		Permissions: inodePermissions(sourceInode),
		UID:         sourceInode.IUid,
		GID:         sourceInode.IGid,
		Modified:    inodeTime(sourceInode.IMtime),
	}, nil
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

	if err := s.requireReadPermission(sourceInode, session, "copiar esta entrada"); err != nil {
		return err
	}

	if sourceInode.IType[0] == '1' {
		content, err := ext2.GetFileContent(diskPath, *sb, sourceIdx)
		if err != nil {
			return err
		}

		_, err = s.Mkfile(FSMkfileInput{
			Token:   session.Token,
			ID:      session.PartID,
			Path:    toPath,
			Content: content,
			Parents: true,
		})

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
