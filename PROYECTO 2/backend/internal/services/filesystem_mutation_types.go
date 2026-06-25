package services

//filesystem_mutation_types.go: contiene las estructuras de datos utilizadas en el servicio de mutación del sistema de archivos, incluyendo las entradas para eliminar, editar, renombrar, copiar y mover archivos y carpetas.

type FSRemoveInput struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Path  string `json:"path"`
}

type FSEditInput struct {
	Token   string `json:"token"`
	ID      string `json:"id"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

type FSRenameInput struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Path  string `json:"path"`
	Name  string `json:"name"`
}

type FSCopyInput struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
}

type FSMoveInput struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
}
