package services

//path_utils.go: contiene funciones auxiliares para manejar rutas y nombres de archivos y carpetas dentro del sistema de archivos, incluyendo la limpieza de rutas, la obtención de nombres base y padres, y la resolución de accesos a particiones montadas.

import (
	"path/filepath"
)

func absolutePath(value string) string {
	abs, err := filepath.Abs(value)
	if err != nil {
		return value
	}

	return abs
}
