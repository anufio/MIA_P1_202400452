package commands

// cmd_users.go: comandos para gestionar usuarios y grupos en el sistema de archivos EXT2, utilizando el archivo users.txt como base de datos de usuarios y grupos

import (
	"MIA_P1_202400452/ext2"
	"fmt"
	"strings"
)

type groupRecord struct {
	ID   int32
	Name string
	Raw  string
}

type userRecord struct {
	ID       int32
	Group    string
	Username string
	Password string
	Raw      string
}

func readUsersFile() (string, ext2.SuperBlock, string, error) {
	sess := GetSession()
	diskPath, partition, _, err := GetMountedPartitionInfo(sess.PartID)
	if err != nil {
		return "", ext2.SuperBlock{}, "", err
	}
	sb, err := ext2.ReadSuperBlock(diskPath, partition.PartStart)
	if err != nil {
		return "", sb, diskPath, err
	}
	if sb.SMagic != 0xEF53 {
		return "", sb, diskPath, fmt.Errorf("la partición no tiene un sistema de archivos EXT2 formateado")
	}
	inode, err := ext2.ReadInode(diskPath, sb, 1)
	if err != nil {
		return "", sb, diskPath, err
	}
	if inode.IType[0] != '1' {
		return "", sb, diskPath, fmt.Errorf("users.txt no es un archivo válido")
	}
	content, err := ext2.GetFileContent(diskPath, sb, 1)
	return content, sb, diskPath, err
}

func saveUsers(diskPath string, sb ext2.SuperBlock, content string, partStart int32) error {
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	inode, err := ext2.ReadInode(diskPath, sb, 1)
	if err != nil {
		return err
	}
	if inode.IType[0] != '1' {
		return fmt.Errorf("users.txt no es un archivo válido")
	}
	if err := ext2.WriteFileContent(diskPath, &sb, &inode, []byte(content)); err != nil {
		return err
	}
	inode.IMtime = ext2.GetCurrentTime()
	if err := ext2.WriteInode(diskPath, sb, 1, inode); err != nil {
		return err
	}
	return ext2.WriteSuperBlock(diskPath, partStart, sb)
}

func parseGroups(content string) []groupRecord {
	var groups []groupRecord
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		parts := splitCSV(trimmed)
		if len(parts) >= 3 && strings.TrimSpace(parts[1]) == "G" {
			groups = append(groups, groupRecord{
				ID:   parseID(parts[0]),
				Name: strings.TrimSpace(parts[2]),
				Raw:  trimmed,
			})
		}
	}
	return groups
}

func parseUsers(content string) []userRecord {
	var users []userRecord
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		parts := splitCSV(trimmed)
		if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" {
			users = append(users, userRecord{
				ID:       parseID(parts[0]),
				Group:    strings.TrimSpace(parts[2]),
				Username: strings.TrimSpace(parts[3]),
				Password: strings.TrimSpace(parts[4]),
				Raw:      trimmed,
			})
		}
	}
	return users
}

func activeGroupExists(content, name string) bool {
	for _, group := range parseGroups(content) {
		if group.ID != 0 && group.Name == name {
			return true
		}
	}
	return false
}

func activeUserExists(content, username string) bool {
	for _, user := range parseUsers(content) {
		if user.ID != 0 && user.Username == username {
			return true
		}
	}
	return false
}

func maxGroupID(content string) int32 {
	max := int32(0)
	for _, group := range parseGroups(content) {
		if group.ID > max {
			max = group.ID
		}
	}
	return max
}

func maxUserID(content string) int32 {
	max := int32(0)
	for _, user := range parseUsers(content) {
		if user.ID > max {
			max = user.ID
		}
	}
	return max
}

func CmdMKGRP(params map[string]string) string {
	if err := RequireRoot(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	name, ok := params["name"]
	if !ok {
		return "Error: falta el parámetro obligatorio -name"
	}
	if len(name) > 10 {
		return "Error: el nombre del grupo no puede exceder 10 caracteres"
	}

	content, sb, diskPath, err := readUsersFile()
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	if activeGroupExists(content, name) {
		return fmt.Sprintf("Error: el grupo '%s' ya existe", name)
	}
	newGID := maxGroupID(content) + 1
	newLine := fmt.Sprintf("%d,G,%s\n", newGID, name)

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += newLine

	sess := GetSession()
	_, partition, _, _ := GetMountedPartitionInfo(sess.PartID)
	if err := saveUsers(diskPath, sb, content, partition.PartStart); err != nil {
		return fmt.Sprintf("Error al guardar users.txt: %v", err)
	}
	return fmt.Sprintf("Grupo '%s' creado exitosamente (GID=%d)", name, newGID)
}

func CmdRMGRP(params map[string]string) string {
	if err := RequireRoot(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	name, ok := params["name"]
	if !ok {
		return "Error: falta el parámetro obligatorio -name"
	}

	content, sb, diskPath, err := readUsersFile()
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	found := false
	var newLines []string
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		parts := splitCSV(trimmed)
		if len(parts) >= 3 && strings.TrimSpace(parts[1]) == "G" {
			grpName := strings.TrimSpace(parts[2])
			gid := parseID(strings.TrimSpace(parts[0]))
			if grpName == name && gid != 0 {
				found = true
				newLines = append(newLines, fmt.Sprintf("0,G,%s", name))
				continue
			}
		}
		newLines = append(newLines, trimmed)
	}
	if !found {
		return fmt.Sprintf("Error: el grupo '%s' no existe o ya fue eliminado", name)
	}
	newContent := strings.Join(newLines, "\n") + "\n"

	sess := GetSession()
	_, partition, _, _ := GetMountedPartitionInfo(sess.PartID)
	if err := saveUsers(diskPath, sb, newContent, partition.PartStart); err != nil {
		return fmt.Sprintf("Error al guardar users.txt: %v", err)
	}
	return fmt.Sprintf("Grupo '%s' eliminado exitosamente.", name)
}

func CmdMKUSR(params map[string]string) string {
	if err := RequireRoot(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	user, ok := params["user"]
	if !ok {
		return "Error: falta el parámetro obligatorio -user"
	}
	if len(user) > 10 {
		return "Error: el nombre de usuario no puede exceder 10 caracteres"
	}
	pass, ok := params["pass"]
	if !ok {
		return "Error: falta el parámetro obligatorio -pass"
	}
	if len(pass) > 10 {
		return "Error: la contraseña no puede exceder 10 caracteres"
	}
	grp, ok := params["grp"]
	if !ok {
		return "Error: falta el parámetro obligatorio -grp"
	}
	if len(grp) > 10 {
		return "Error: el nombre del grupo no puede exceder 10 caracteres"
	}

	content, sb, diskPath, err := readUsersFile()
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	if !activeGroupExists(content, grp) {
		return fmt.Sprintf("Error: el grupo '%s' no existe o está eliminado", grp)
	}
	if activeUserExists(content, user) {
		return fmt.Sprintf("Error: el usuario '%s' ya existe", user)
	}
	newUID := maxUserID(content) + 1
	newLine := fmt.Sprintf("%d,U,%s,%s,%s\n", newUID, grp, user, pass)

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += newLine

	sess := GetSession()
	_, partition, _, _ := GetMountedPartitionInfo(sess.PartID)
	if err := saveUsers(diskPath, sb, content, partition.PartStart); err != nil {
		return fmt.Sprintf("Error al guardar users.txt: %v", err)
	}
	return fmt.Sprintf("Usuario '%s' creado exitosamente (UID=%d, Grupo=%s)", user, newUID, grp)
}

func CmdRMUSR(params map[string]string) string {
	if err := RequireRoot(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	user, ok := params["user"]
	if !ok {
		return "Error: falta el parámetro obligatorio -user"
	}

	content, sb, diskPath, err := readUsersFile()
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	found := false
	var newLines []string
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		parts := splitCSV(trimmed)
		if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" {
			uid := parseID(strings.TrimSpace(parts[0]))
			uname := strings.TrimSpace(parts[3])
			if uname == user && uid != 0 {
				found = true
				newLines = append(newLines, fmt.Sprintf("0,U,%s,%s,%s", strings.TrimSpace(parts[2]), uname, strings.TrimSpace(parts[4])))
				continue
			}
		}
		newLines = append(newLines, trimmed)
	}
	if !found {
		return fmt.Sprintf("Error: el usuario '%s' no existe o ya fue eliminado", user)
	}
	newContent := strings.Join(newLines, "\n") + "\n"

	sess := GetSession()
	_, partition, _, _ := GetMountedPartitionInfo(sess.PartID)
	if err := saveUsers(diskPath, sb, newContent, partition.PartStart); err != nil {
		return fmt.Sprintf("Error al guardar users.txt: %v", err)
	}
	return fmt.Sprintf("Usuario '%s' eliminado exitosamente.", user)
}

func CmdCHGRP(params map[string]string) string {
	if err := RequireRoot(); err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	user, ok := params["user"]
	if !ok {
		return "Error: falta el parámetro obligatorio -user"
	}
	grp, ok := params["grp"]
	if !ok {
		return "Error: falta el parámetro obligatorio -grp"
	}

	content, sb, diskPath, err := readUsersFile()
	if err != nil {
		return fmt.Sprintf("Error al leer users.txt: %v", err)
	}

	if !activeGroupExists(content, grp) {
		return fmt.Sprintf("Error: el grupo '%s' no existe o está eliminado", grp)
	}

	found := false
	var newLines []string
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		parts := splitCSV(trimmed)
		if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" {
			uid := parseID(strings.TrimSpace(parts[0]))
			uname := strings.TrimSpace(parts[3])
			if uname == user && uid != 0 {
				found = true
				newLines = append(newLines, fmt.Sprintf("%d,U,%s,%s,%s", uid, grp, uname, strings.TrimSpace(parts[4])))
				continue
			}
		}
		newLines = append(newLines, trimmed)
	}
	if !found {
		return fmt.Sprintf("Error: el usuario '%s' no existe o está eliminado", user)
	}
	newContent := strings.Join(newLines, "\n") + "\n"

	sess := GetSession()
	_, partition, _, _ := GetMountedPartitionInfo(sess.PartID)
	if err := saveUsers(diskPath, sb, newContent, partition.PartStart); err != nil {
		return fmt.Sprintf("Error al guardar users.txt: %v", err)
	}
	return fmt.Sprintf("Grupo del usuario '%s' cambiado a '%s' exitosamente.", user, grp)
}
