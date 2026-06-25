package handlers

//filesystem_handler.go: contiene el controlador para manejar las rutas relacionadas con el sistema de archivos, incluyendo la obtención de la lista de archivos y directorios, la lectura de archivos, la creación de directorios y archivos, la eliminación, edición, renombrado, copia y movimiento de archivos y directorios.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type FileSystemHandler struct {
	app *services.App
}

func NewFileSystemHandler(app *services.App) *FileSystemHandler {
	return &FileSystemHandler{
		app: app,
	}
}

func (h *FileSystemHandler) List(w http.ResponseWriter, r *http.Request) {
	var input services.FSListInput

	if r.Method == http.MethodGet {
		input.Token = tokenFromRequest(r)
		input.ID = r.URL.Query().Get("id")
		input.Path = r.URL.Query().Get("path")
	} else if r.Method == http.MethodPost {
		if err := readJSON(r, &input); err != nil {
			respond.BadRequest(w, "JSON inválido")
			return
		}
		if input.Token == "" {
			input.Token = tokenFromRequest(r)
		}
	} else {
		respond.MethodNotAllowed(w)
		return
	}

	items, err := h.app.FileSystemService.List(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Carpeta leída correctamente",
		"items":   items,
		"data":    items,
	})
}

func (h *FileSystemHandler) Read(w http.ResponseWriter, r *http.Request) {
	var input services.FSReadInput

	if r.Method == http.MethodGet {
		input.Token = tokenFromRequest(r)
		input.ID = r.URL.Query().Get("id")
		input.Path = r.URL.Query().Get("path")
	} else if r.Method == http.MethodPost {
		if err := readJSON(r, &input); err != nil {
			respond.BadRequest(w, "JSON inválido")
			return
		}
		if input.Token == "" {
			input.Token = tokenFromRequest(r)
		}
	} else {
		respond.MethodNotAllowed(w)
		return
	}

	file, err := h.app.FileSystemService.ReadFile(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Archivo leído correctamente",
		"file":    file,
		"data":    file,
	})
}
