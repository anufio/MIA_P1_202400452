package services

//filesystem_types.go: contiene las estructuras de datos utilizadas en el servicio del sistema de archivos, incluyendo la información de archivos y carpetas, así como las entradas para listar, leer, crear directorios y crear archivos.

type FileSystemService struct {
	diskService *DiskService
	authService *AuthService
}

type FSListInput struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Path  string `json:"path"`
}

type FSReadInput struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Path  string `json:"path"`
}

type FSMkdirInput struct {
	Token   string `json:"token"`
	ID      string `json:"id"`
	Path    string `json:"path"`
	Parents bool   `json:"parents"`
}

type FSMkfileInput struct {
	Token   string `json:"token"`
	ID      string `json:"id"`
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
	Parents bool   `json:"parents"`
}

type FSItem struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Inode       int32  `json:"inode"`
	Size        int32  `json:"size"`
	Permissions string `json:"permissions"`
	UID         int32  `json:"uid"`
	GID         int32  `json:"gid"`
	Modified    string `json:"modified"`
}

type FSReadResult struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Inode   int32  `json:"inode"`
	Size    int32  `json:"size"`
	Content string `json:"content"`
}

func NewFileSystemService(diskService *DiskService, authService *AuthService) *FileSystemService {
	return &FileSystemService{
		diskService: diskService,
		authService: authService,
	}
}
