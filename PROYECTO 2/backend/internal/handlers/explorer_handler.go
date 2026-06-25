package handlers

//explorer_handler.go: contiene el controlador para manejar las rutas relacionadas con el explorador de archivos, incluyendo la obtención de la lista de archivos y directorios, la lectura de archivos, la creación de directorios y archivos, la eliminación, edición, renombrado, copia y movimiento de archivos y directorios.

import (
	"net/http"
	"strings"

	"MIA_P2_202400452/internal/respond"
)

func (h *FileSystemHandler) Explorer(w http.ResponseWriter, r *http.Request) {
	action := strings.TrimPrefix(r.URL.Path, "/api/explorer/")
	action = strings.Trim(action, "/")

	switch action {
	case "list":
		h.List(w, r)

	case "read":
		h.Read(w, r)

	case "mkdir":
		h.Mkdir(w, r)

	case "mkfile":
		h.Mkfile(w, r)

	case "remove":
		h.Remove(w, r)

	case "edit":
		h.Edit(w, r)

	case "rename":
		h.Rename(w, r)

	case "copy":
		h.Copy(w, r)

	case "move":
		h.Move(w, r)

	default:
		respond.NotFound(w, "Acción de explorador no encontrada")
	}
}
