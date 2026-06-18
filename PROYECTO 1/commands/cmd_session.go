package commands

// cmd_session.go: comandos relacionados con la gestión de sesiones de usuario (login/logout)

import (
	"MIA_P1_202400452/ext2"
	"fmt"
	"strings"
)

type Session struct {
	Active   bool
	Username string
	UID      int32
	GID      int32
	PartID   string
}

var currentSession = Session{}

func CmdLOGIN(params map[string]string) string {
	if currentSession.Active {
		return fmt.Sprintf("Error: ya existe una sesión activa del usuario '%s'. Ejecute logout primero.", currentSession.Username)
	}

	user, ok := params["user"]
	if !ok {
		return "Error: falta el parámetro obligatorio -user"
	}
	pass, ok := params["pass"]
	if !ok {
		return "Error: falta el parámetro obligatorio -pass"
	}
	id, ok := params["id"]
	if !ok {
		return "Error: falta el parámetro obligatorio -id"
	}

	diskPath, partition, _, err := GetMountedPartitionInfo(id)
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

	usersInodeIdx, err := ext2.FindInodeByPath(diskPath, sb, "/users.txt")
	if err != nil {
		return fmt.Sprintf("Error: no se encontró users.txt: %v", err)
	}
	usersInode, err := ext2.ReadInode(diskPath, sb, usersInodeIdx)
	if err != nil {
		return fmt.Sprintf("Error al leer inodo de users.txt: %v", err)
	}
	if usersInode.IType[0] != '1' {
		return "Error: users.txt no es un archivo válido"
	}

	content, err := ext2.GetFileContent(diskPath, sb, usersInodeIdx)
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	uid, gid, found := findUser(content, user, pass)
	if !found {
		return "Error: usuario o contraseña incorrectos"
	}
	if uid == 0 {
		return fmt.Sprintf("Error: el usuario '%s' está eliminado", user)
	}

	currentSession = Session{
		Active:   true,
		Username: user,
		UID:      uid,
		GID:      gid,
		PartID:   strings.ToUpper(id),
	}

	return fmt.Sprintf("Login exitoso. Bienvenido, %s (UID=%d, GID=%d)", user, uid, gid)
}

func CmdLOGOUT(params map[string]string) string {
	if !currentSession.Active {
		return "Error: no hay una sesión activa"
	}
	user := currentSession.Username
	currentSession = Session{}
	return fmt.Sprintf("Sesión de '%s' cerrada exitosamente.", user)
}

func GetSession() Session {
	return currentSession
}

func RequireSession() error {
	if !currentSession.Active {
		return fmt.Errorf("debe iniciar sesión para ejecutar este comando")
	}
	return nil
}

func RequireRoot() error {
	if err := RequireSession(); err != nil {
		return err
	}
	if currentSession.Username != "root" {
		return fmt.Errorf("solo el usuario root puede ejecutar este comando")
	}
	return nil
}

func findUser(content, username, password string) (int32, int32, bool) {
	lines := strings.Split(content, "\n")

	groupMap := make(map[string]int32)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := splitCSV(line)
		if len(parts) >= 3 && strings.TrimSpace(parts[1]) == "G" {
			gidStr := strings.TrimSpace(parts[0])
			gid := parseID(gidStr)
			grpName := strings.TrimSpace(parts[2])
			if gid != 0 {
				groupMap[grpName] = gid
			}
		}
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := splitCSV(line)
		if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" {
			uidStr := strings.TrimSpace(parts[0])
			uid := parseID(uidStr)
			grpName := strings.TrimSpace(parts[2])
			uname := strings.TrimSpace(parts[3])
			upass := strings.TrimSpace(parts[4])
			if uname == username && upass == password {
				gid, groupFound := groupMap[grpName]
				if !groupFound {
					return 0, 0, false
				}
				return uid, gid, true
			}
		}
	}
	return 0, 0, false
}

func splitCSV(line string) []string {
	return strings.Split(line, ",")
}

func parseID(s string) int32 {
	s = strings.TrimSpace(s)
	var n int32
	fmt.Sscanf(s, "%d", &n)
	return n
}
