package services

//filesystem_write_service.go: contiene la lógica para manejar la escritura de archivos y carpetas dentro del sistema de archivos, incluyendo la validación de rutas, la escritura de contenido y la gestión de permisos.

import (
	"fmt"
	"strings"

	"MIA_P2_202400452/internal/ext2"
)

func (s *FileSystemService) Mkdir(input FSMkdirInput) (FSItem, error) {
	mounted, sb, session, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return FSItem{}, err
	}

	targetPath := cleanFSPath(input.Path)

	if targetPath == "/" {
		return FSItem{}, fmt.Errorf("no se puede crear la raíz")
	}

	if input.Parents {
		if err := s.ensureDirectory(mounted.DiskPath, &sb, targetPath, session); err != nil {
			return FSItem{}, err
		}
	} else {
		if err := s.createDirectory(mounted.DiskPath, &sb, parentPath(targetPath), baseName(targetPath), session); err != nil {
			return FSItem{}, err
		}
	}

	_ = ext2.WriteSuperBlock(mounted.DiskPath, mounted.Start, sb)

	inodeIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, targetPath)
	if err != nil {
		return FSItem{}, err
	}

	inode, _ := ext2.ReadInode(mounted.DiskPath, sb, inodeIdx)

	return FSItem{
		Name:        baseName(targetPath),
		Path:        targetPath,
		Type:        "folder",
		Inode:       inodeIdx,
		Size:        inode.IS,
		Permissions: inodePermissions(inode),
		UID:         inode.IUid,
		GID:         inode.IGid,
		Modified:    inodeTime(inode.IMtime),
	}, nil
}

func (s *FileSystemService) Mkfile(input FSMkfileInput) (FSItem, error) {
	mounted, sb, session, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return FSItem{}, err
	}

	targetPath := cleanFSPath(input.Path)
	fileName := baseName(targetPath)
	parent := parentPath(targetPath)

	if fileName == "." || fileName == "/" || fileName == "" {
		return FSItem{}, fmt.Errorf("nombre de archivo inválido")
	}

	if len(fileName) > 12 {
		return FSItem{}, fmt.Errorf("el nombre no puede exceder 12 caracteres")
	}

	if input.Parents {
		if err := s.ensureDirectory(mounted.DiskPath, &sb, parent, session); err != nil {
			return FSItem{}, err
		}
	}

	if _, err := ext2.FindInodeByPath(mounted.DiskPath, sb, targetPath); err == nil {
		return FSItem{}, fmt.Errorf("ya existe un archivo o carpeta en esa ruta")
	}

	parentIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, parent)
	if err != nil {
		return FSItem{}, err
	}

	content := input.Content

	if input.Size > int64(len(content)) {
		content += strings.Repeat("0", int(input.Size)-len(content))
	}

	inodeIdx, err := ext2.AllocateInode(mounted.DiskPath, &sb)
	if err != nil {
		return FSItem{}, err
	}

	inode := ext2.NewInode()
	inode.IType = [1]byte{'1'}
	inode.IPerm = [3]byte{'6', '6', '4'}
	inode.IUid = session.UID
	inode.IGid = session.GID
	inode.ICtime = ext2.GetCurrentTime()
	inode.IAtime = ext2.GetCurrentTime()
	inode.IMtime = ext2.GetCurrentTime()

	if inode.IUid == 0 {
		inode.IUid = 1
	}

	if inode.IGid == 0 {
		inode.IGid = 1
	}

	if err := ext2.WriteFileContent(mounted.DiskPath, &sb, &inode, []byte(content)); err != nil {
		return FSItem{}, err
	}

	if err := ext2.WriteInode(mounted.DiskPath, sb, inodeIdx, inode); err != nil {
		return FSItem{}, err
	}

	if err := ext2.AddEntryToFolder(mounted.DiskPath, &sb, parentIdx, fileName, inodeIdx); err != nil {
		return FSItem{}, err
	}

	_ = ext2.WriteSuperBlock(mounted.DiskPath, mounted.Start, sb)

	return FSItem{
		Name:        fileName,
		Path:        targetPath,
		Type:        "file",
		Inode:       inodeIdx,
		Size:        inode.IS,
		Permissions: inodePermissions(inode),
		UID:         inode.IUid,
		GID:         inode.IGid,
		Modified:    inodeTime(inode.IMtime),
	}, nil
}
