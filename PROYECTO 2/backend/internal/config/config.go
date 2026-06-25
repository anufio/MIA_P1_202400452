package config

//config.go: contiene la configuración del servidor, incluyendo la dirección del servidor, la ruta raíz de los discos y la ruta raíz de los reportes.

import (
	"os"
	"path/filepath"
)

type Config struct {
	ServerAddress string
	DiskRoot      string
	ReportRoot    string
}

func Load() Config {
	diskRoot := getEnv("DISK_ROOT", "data/disks")
	reportRoot := getEnv("REPORT_ROOT", "data/reports")

	diskRoot = toAbsolute(diskRoot)
	reportRoot = toAbsolute(reportRoot)

	_ = os.MkdirAll(diskRoot, 0755)
	_ = os.MkdirAll(reportRoot, 0755)

	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DiskRoot:      diskRoot,
		ReportRoot:    reportRoot,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func toAbsolute(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	return abs
}
