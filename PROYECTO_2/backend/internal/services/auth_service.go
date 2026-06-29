package services

//auth_service.go: contiene la lógica para manejar la autenticación de usuarios, incluyendo el inicio y cierre de sesión, así como la validación de credenciales y la gestión de sesiones.

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"MIA_P2_202400452/internal/ext2"
)

type AuthService struct {
	diskService *DiskService
	mutex       sync.Mutex
	sessions    map[string]SessionInfo
}

func NewAuthService(diskService *DiskService) *AuthService {
	return &AuthService{
		diskService: diskService,
		sessions:    make(map[string]SessionInfo),
	}
}

func (s *AuthService) Login(input LoginInput) (SessionInfo, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	partID := strings.ToUpper(strings.TrimSpace(input.ID))
	username := strings.TrimSpace(firstNonEmpty(input.User, input.Username))
	password := strings.TrimSpace(input.Password)

	if partID == "" {
		return SessionInfo{}, fmt.Errorf("debe indicar el id de la partición")
	}

	if username == "" {
		return SessionInfo{}, fmt.Errorf("debe indicar el usuario")
	}

	if password == "" {
		return SessionInfo{}, fmt.Errorf("debe indicar la contraseña")
	}

	mounted, ok := s.diskService.GetMounted(partID)
	if !ok {
		return SessionInfo{}, fmt.Errorf("la partición '%s' no está montada", partID)
	}

	sb, err := ext2.ReadSuperBlock(mounted.DiskPath, mounted.Start)
	if err != nil {
		return SessionInfo{}, fmt.Errorf("error leyendo superbloque: %v", err)
	}

	if sb.SMagic != 0xEF53 {
		return SessionInfo{}, fmt.Errorf("la partición no tiene EXT2 formateado")
	}

	content, err := ext2.GetFileContent(mounted.DiskPath, sb, 1)
	if err != nil {
		return SessionInfo{}, fmt.Errorf("error leyendo users.txt: %v", err)
	}

	uid, gid, group, ok := validateUser(content, username, password)
	if !ok {
		return SessionInfo{}, fmt.Errorf("usuario o contraseña incorrectos")
	}

	token := fmt.Sprintf("%s-%d", username, time.Now().UnixNano())

	session := SessionInfo{
		Token:    token,
		Active:   true,
		User:     username,
		UID:      uid,
		GID:      gid,
		Group:    group,
		PartID:   partID,
		DiskPath: mounted.DiskPath,
	}

	s.sessions[token] = session

	return session, nil
}

func (s *AuthService) Logout(input LogoutInput) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	token := strings.TrimSpace(input.Token)

	if token == "" {
		return fmt.Errorf("debe indicar el token")
	}

	if _, ok := s.sessions[token]; !ok {
		return fmt.Errorf("la sesión no existe")
	}

	delete(s.sessions, token)

	return nil
}

func (s *AuthService) GetSession(token string) (SessionInfo, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	token = strings.TrimSpace(token)
	session, ok := s.sessions[token]

	return session, ok
}

func validateUser(content string, username string, password string) (int32, int32, string, bool) {
	groups := map[string]int32{}

	for _, line := range strings.Split(content, "\n") {
		parts := splitCSVLine(line)

		if len(parts) >= 3 && strings.TrimSpace(parts[1]) == "G" {
			id := int32(parseInt64(parts[0]))
			name := strings.TrimSpace(parts[2])

			if id != 0 {
				groups[name] = id
			}
		}
	}

	for _, line := range strings.Split(content, "\n") {
		parts := splitCSVLine(line)

		if len(parts) >= 5 && strings.TrimSpace(parts[1]) == "U" {
			uid := int32(parseInt64(parts[0]))
			group := strings.TrimSpace(parts[2])
			user := strings.TrimSpace(parts[3])
			pass := strings.TrimSpace(parts[4])

			if uid != 0 && user == username && pass == password {
				return uid, groups[group], group, true
			}
		}
	}

	return 0, 0, "", false
}

func splitCSVLine(line string) []string {
	raw := strings.Split(strings.TrimSpace(line), ",")
	result := make([]string, 0, len(raw))

	for _, part := range raw {
		result = append(result, strings.TrimSpace(part))
	}

	return result
}
