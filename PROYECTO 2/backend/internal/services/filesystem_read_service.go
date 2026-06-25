package services

//filesystem_read_service.go: contiene la lógica para manejar la lectura de archivos y carpetas dentro del sistema de archivos, incluyendo la validación de rutas, la obtención de contenido y la gestión de permisos.

import (
	"fmt"
	"sort"

	"MIA_P2_202400452/internal/ext2"
)

func (s *FileSystemService) List(input FSListInput) ([]FSItem, error) {
	mounted, sb, _, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
		return nil, err
	}

	targetPath := cleanFSPath(input.Path)

	inodeIdx, err := ext2.FindInodeByPath(mounted.DiskPath, sb, targetPath)
	if err != nil {
		return nil, err
	}

	inode, err := ext2.ReadInode(mounted.DiskPath, sb, inodeIdx)
	if err != nil {
		return nil, err
	}

	if inode.IType[0] != '0' {
		return nil, fmt.Errorf("la ruta no es una carpeta")
	}

	items := make([]FSItem, 0)

	for i := 0; i < 12; i++ {
		if inode.IBlock[i] == -1 {
			continue
		}

		block, err := ext2.ReadFolderBlock(mounted.DiskPath, sb, inode.IBlock[i])
		if err != nil {
			continue
		}

		for _, content := range block.BContent {
			name := trimBlockName(content.BName[:])

			if name == "" || name == "." || name == ".." || content.BInode == -1 {
				continue
			}

			childInode, err := ext2.ReadInode(mounted.DiskPath, sb, content.BInode)
			if err != nil {
				continue
			}

			items = append(items, FSItem{
				Name:        name,
				Path:        joinFSPath(targetPath, name),
				Type:        inodeTypeLabel(childInode),
				Inode:       content.BInode,
				Size:        childInode.IS,
				Permissions: inodePermissions(childInode),
				UID:         childInode.IUid,
				GID:         childInode.IGid,
				Modified:    inodeTime(childInode.IMtime),
			})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Type == items[j].Type {
			return items[i].Name < items[j].Name
		}

		return items[i].Type == "folder"
	})

	return items, nil
}

func (s *FileSystemService) ReadFile(input FSReadInput) (FSReadResult, error) {
	mounted, sb, _, err := s.resolveAccess(input.Token, input.ID)
	if err != nil {
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

	content, err := ext2.GetFileContent(mounted.DiskPath, sb, inodeIdx)
	if err != nil {
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
