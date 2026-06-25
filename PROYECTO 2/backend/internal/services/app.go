package services

//app.go: contiene la estructura principal de la aplicación, que incluye los servicios de disco, autenticación, sistema de archivos y reportes.

import (
	"MIA_P2_202400452/internal/config"
)

type App struct {
	Config            config.Config
	DiskService       *DiskService
	AuthService       *AuthService
	FileSystemService *FileSystemService
	ReportService     *ReportService
}

func NewApp(cfg config.Config) *App {
	diskService := NewDiskService(cfg.DiskRoot)
	authService := NewAuthService(diskService)

	return &App{
		Config:            cfg,
		DiskService:       diskService,
		AuthService:       authService,
		FileSystemService: NewFileSystemService(diskService, authService),
		ReportService:     NewReportService(cfg.ReportRoot, diskService),
	}
}
