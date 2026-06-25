package services

//filesystem_edit_rename_service.go: contiene la lógica para manejar la edición y el renombrado de archivos dentro del sistema de archivos, incluyendo la validación de rutas, la escritura de contenido y la gestión de permisos.

import (
	"fmt"
	"os"
	"strings"

	"MIA_P2_202400452/internal/ext2"
)

func (s *FileSystemService) Edit(input FSEditInput) (FSReadResult, error) {
	mounted, sb, session, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return FSReadResult{}, err
	}

	if err := requireSession(session); err != nil {
		return FSReadResult{}, err
	}

	targetPath := cleanFSPath(input.Path)

	inodeIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, targetPath)
	if err != nil {
		return FSReadResult{}, err
	}

	inode, err := ext2.ReadInode(mounted.DiskPath, sb, inodeIdx)
	if err != nil {
		return FSReadResult{}, err
	}

	if inode.IType[0] != '1' {
		return FSReadResult{}, fmt.Errorf("la ruta no es un archivo")
	}

	if err := s.requireReadWritePermission(inode, session, "editar el archivo"); err != nil {
		return FSReadResult{}, err
	}

	content := firstNonEmpty(input.Content, input.Cont)

	if strings.TrimSpace(input.Contenido) != "" {
		data, err := os.ReadFile(strings.TrimSpace(input.Contenido))
		if err != nil {
			return FSReadResult{}, fmt.Errorf("no se pudo leer el archivo indicado en -contenido: %v", err)
		}

		content = string(data)
	}

	if strings.TrimSpace(content) == "" {
		return FSReadResult{}, fmt.Errorf("debe indicar contenido con -contenido o -cont")
	}

	if err := ext2.WriteFileContent(mounted.DiskPath, &sb, &inode, []byte(content)); err != nil {
		return FSReadResult{}, err
	}

	inode.IMtime = ext2.GetCurrentTime()

	if err := ext2.WriteInode(mounted.DiskPath, sb, inodeIdx, inode); err != nil {
		return FSReadResult{}, err
	}

	if err := ext2.WriteSuperBlock(mounted.DiskPath, mounted.Start, sb); err != nil {
		return FSReadResult{}, err
	}

	return FSReadResult{
		Path:    targetPath,
		Name:    baseName(targetPath),
		Type:    "file",
		Inode:   inodeIdx,
		Size:    inode.IS,
		Content: content,
	}, nil
}

func (s *FileSystemService) Rename(input FSRenameInput) (FSItem, error) {
	mounted, sb, session, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return FSItem{}, err
	}

	if err := requireSession(session); err != nil {
		return FSItem{}, err
	}

	targetPath := cleanFSPath(input.Path)
	newName := strings.TrimSpace(input.Name)

	if targetPath == "/" {
		return FSItem{}, fmt.Errorf("no se puede renombrar la raíz")
	}

	if newName == "" || newName == "." || newName == ".." || strings.Contains(newName, "/") {
		return FSItem{}, fmt.Errorf("nombre nuevo inválido")
	}

	if len(newName) > 12 {
		return FSItem{}, fmt.Errorf("el nombre no puede exceder 12 caracteres")
	}

	inodeIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, targetPath)
	if err != nil {
		return FSItem{}, err
	}

	inode, err := ext2.ReadInode(mounted.DiskPath, sb, inodeIdx)
	if err != nil {
		return FSItem{}, err
	}

	if err := s.requireWritePermission(inode, session, "renombrar esta entrada"); err != nil {
		return FSItem{}, err
	}

	parent := parentPath(targetPath)
	parentIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, parent)
	if err != nil {
		return FSItem{}, err
	}

	parentInode, err := ext2.ReadInode(mounted.DiskPath, sb, parentIdx)
	if err != nil {
		return FSItem{}, err
	}

	if err := s.requireWritePermission(parentInode, session, "modificar la carpeta padre"); err != nil {
		return FSItem{}, err
	}

	newPath := joinFSPath(parent, newName)

	if _, err := ext2.FindInodeByPath(mounted.DiskPath, sb, newPath); err == nil {
		return FSItem{}, fmt.Errorf("ya existe una entrada con ese nombre")
	}

	if err := ext2.RemoveEntryFromFolder(mounted.DiskPath, &sb, parentIdx, baseName(targetPath)); err != nil {
		return FSItem{}, err
	}

	if err := ext2.AddEntryToFolder(mounted.DiskPath, &sb, parentIdx, newName, inodeIdx); err != nil {
		return FSItem{}, err
	}

	inode.IMtime = ext2.GetCurrentTime()

	if err := ext2.WriteInode(mounted.DiskPath, sb, inodeIdx, inode); err != nil {
		return FSItem{}, err
	}

	if err := ext2.WriteSuperBlock(mounted.DiskPath, mounted.Start, sb); err != nil {
		return FSItem{}, err
	}

	return FSItem{
		Name:        newName,
		Path:        newPath,
		Type:        inodeTypeLabel(inode),
		Inode:       inodeIdx,
		Size:        inode.IS,
		Permissions: inodePermissions(inode),
		UID:         inode.IUid,
		GID:         inode.IGid,
		Modified:    inodeTime(inode.IMtime),
	}, nil
}
