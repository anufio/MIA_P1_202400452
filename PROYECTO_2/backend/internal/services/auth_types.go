package services

//auth_types.go: contiene las estructuras de datos utilizadas en el servicio de autenticación, incluyendo las entradas para login y logout, así como la información de sesión del usuario autenticado.

type MkfsInput struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type LoginInput struct {
	ID       string `json:"id"`
	User     string `json:"user"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutInput struct {
	Token string `json:"token"`
}

type SessionInfo struct {
	Token    string `json:"token"`
	Active   bool   `json:"active"`
	User     string `json:"user"`
	UID      int32  `json:"uid"`
	GID      int32  `json:"gid"`
	Group    string `json:"group"`
	PartID   string `json:"partId"`
	DiskPath string `json:"diskPath"`
}
