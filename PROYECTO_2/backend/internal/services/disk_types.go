package services

//disk_types.go: contiene las estructuras de datos utilizadas en el servicio de discos, incluyendo la información de discos, particiones y montajes, así como las entradas para crear, eliminar y redimensionar particiones.

import (
	"sync"
)

type DiskInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	Fit       string `json:"fit"`
	CreatedAt string `json:"createdAt"`
}

type PartitionInfo struct {
	ID          string `json:"id"`
	DiskID      string `json:"diskId"`
	DiskPath    string `json:"diskPath"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Fit         string `json:"fit"`
	Start       int32  `json:"start"`
	Size        int32  `json:"size"`
	Correlative int32  `json:"correlative"`
	Mounted     bool   `json:"mounted"`
	Formatted   bool   `json:"formatted"`
}

type MountedInfo struct {
	ID        string `json:"id"`
	DiskPath  string `json:"diskPath"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Start     int32  `json:"start"`
	Size      int32  `json:"size"`
	Index     int    `json:"index"`
	IsLogical bool   `json:"isLogical"`
	EBRStart  int64  `json:"ebrStart"`
}

type CreateDiskInput struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Unit string `json:"unit"`
	Fit  string `json:"fit"`
}

type CreatePartitionInput struct {
	DiskPath string `json:"diskPath"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Unit     string `json:"unit"`
	Type     string `json:"type"`
	Fit      string `json:"fit"`
}

type DeletePartitionInput struct {
	DiskPath   string `json:"diskPath"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	DeleteType string `json:"deleteType"`
	Delete     string `json:"delete"`
}

type ResizePartitionInput struct {
	DiskPath string `json:"diskPath"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Add      int64  `json:"add"`
	Unit     string `json:"unit"`
}

type MountInput struct {
	DiskPath string `json:"diskPath"`
	Path     string `json:"path"`
	Name     string `json:"name"`
}

type UnmountInput struct {
	ID string `json:"id"`
}

type DiskService struct {
	root        string
	mutex       sync.Mutex
	mounted     map[string]MountedInfo
	diskLetters map[string]string
	nextLetter  int
}

func NewDiskService(root string) *DiskService {
	return &DiskService{
		root:        root,
		mounted:     make(map[string]MountedInfo),
		diskLetters: make(map[string]string),
	}
}

func (s *DiskService) Root() string {
	return s.root
}
