package services

//reset_service.go: contiene la lógica para manejar el reinicio del sistema de archivos, incluyendo la eliminación de todos los discos y reportes, y la reconfiguración de los servicios del sistema.

import "os"

func (app *App) ResetAll() error {
	if err := os.RemoveAll(app.Config.DiskRoot); err != nil {
		return err
	}

	if err := os.RemoveAll(app.Config.ReportRoot); err != nil {
		return err
	}

	if err := os.MkdirAll(app.Config.DiskRoot, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(app.Config.ReportRoot, 0755); err != nil {
		return err
	}

	fresh := NewApp(app.Config)

	app.DiskService = fresh.DiskService
	app.AuthService = fresh.AuthService
	app.FileSystemService = fresh.FileSystemService
	app.ReportService = fresh.ReportService

	return nil
}
