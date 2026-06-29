package services

//filesystem_utils.go: contiene funciones auxiliares para manejar rutas y nombres de archivos y carpetas dentro del sistema de archivos, incluyendo la limpieza de rutas, la obtención de nombres base y padres, y la resolución de accesos a particiones montadas.

import (
	"fmt"
	"path"
	"strings"

	"MIA_P2_202400452/internal/ext2"
)

func cleanFSPath(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return "/"
	}

	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}

	return path.Clean(value)
}

func baseName(value string) string {
	value = cleanFSPath(value)
	return path.Base(value)
}

func parentPath(value string) string {
	value = cleanFSPath(value)
	parent := path.Dir(value)

	if parent == "." {
		return "/"
	}

	return parent
}

func joinFSPath(parent string, name string) string {
	parent = cleanFSPath(parent)

	if parent == "/" {
		return "/" + name
	}

	return parent + "/" + name
}

func trimBlockName(raw []byte) string {
	return strings.TrimRight(string(raw), "\x00")
}

func inodeTypeLabel(inode ext2.Inode) string {
	if inode.IType[0] == '0' {
		return "folder"
	}

	return "file"
}

func inodePermissions(inode ext2.Inode) string {
	return string(inode.IPerm[:])
}

func inodeTime(raw [19]byte) string {
	return strings.TrimRight(string(raw[:]), "\x00")
}

func (s *FileSystemService) resolveAccess(token string, id string) (MountedInfo, ext2.SuperBlock, SessionInfo, error) {
	token = strings.TrimSpace(token)
	id = strings.ToUpper(strings.TrimSpace(id))

	var session SessionInfo

	if token != "" {
		foundSession, ok := s.authService.GetSession(token)
		if !ok {
			return MountedInfo{}, ext2.SuperBlock{}, SessionInfo{}, fmt.Errorf("sesión no válida")
		}

		session = foundSession

		if id == "" {
			id = session.PartID
		}
	}

	if id == "" {
		return MountedInfo{}, ext2.SuperBlock{}, SessionInfo{}, fmt.Errorf("debe indicar id de partición o token")
	}

	mounted, ok := s.diskService.GetMounted(id)
	if !ok {
		return MountedInfo{}, ext2.SuperBlock{}, SessionInfo{}, fmt.Errorf("la partición '%s' no está montada", id)
	}

	sb, err := ext2.ReadSuperBlock(mounted.DiskPath, mounted.Start)
	if err != nil {
		return MountedInfo{}, ext2.SuperBlock{}, SessionInfo{}, fmt.Errorf("error leyendo superbloque: %v", err)
	}

	if sb.SMagic != 0xEF53 {
		return MountedInfo{}, ext2.SuperBlock{}, SessionInfo{}, fmt.Errorf("la partición no está formateada como EXT2")
	}

	return mounted, sb, session, nil
}
