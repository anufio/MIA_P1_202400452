package services

//report_types.go: contiene las definiciones de tipos de datos para manejar la generación de reportes dentro del sistema de archivos, incluyendo la estructura de entrada y salida para los reportes.

type ReportInput struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	FilePath string `json:"filePath"`
}

type ReportOutput struct {
	Type   string `json:"type"`
	Path   string `json:"path"`
	URL    string `json:"url"`
	Format string `json:"format"`
}
