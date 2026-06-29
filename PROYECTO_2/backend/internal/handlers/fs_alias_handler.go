package handlers

//fs_alias_handler.go: contiene el controlador para manejar las rutas relacionadas con los alias del sistema de archivos, incluyendo la creación, eliminación y listado de alias.

import (
	"net/http"
	"strings"

	"MIA_P2_202400452/internal/respond"
)

func (h *FileSystemHandler) FSAction(w http.ResponseWriter, r *http.Request) {
	action := strings.TrimPrefix(r.URL.Path, "/api/fs/")
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

	default:
		respond.NotFound(w, "Acción de sistema de archivos no encontrada")
	}
}
