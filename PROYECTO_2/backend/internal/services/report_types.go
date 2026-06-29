package services

// report_types.go: entradas y salidas para generar reportes desde la API.

type ReportInput struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	DiskPath string `json:"diskPath"`
	FilePath string `json:"filePath"`
}

type ReportOutput struct {
	Type   string `json:"type"`
	Path   string `json:"path"`
	URL    string `json:"url"`
	Format string `json:"format"`
}
