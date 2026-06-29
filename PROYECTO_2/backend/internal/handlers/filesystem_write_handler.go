package handlers

//filesystem_write_handler.go: contiene el controlador para manejar las rutas relacionadas con la escritura en el sistema de archivos, incluyendo la creación de directorios y archivos.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

func (h *FileSystemHandler) Mkdir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.FSMkdirInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if input.Token == "" {
		input.Token = tokenFromRequest(r)
	}

	item, err := h.app.FileSystemService.Mkdir(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Carpeta creada correctamente",
		"item":    item,
		"data":    item,
	})
}

func (h *FileSystemHandler) Mkfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.FSMkfileInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if input.Token == "" {
		input.Token = tokenFromRequest(r)
	}

	item, err := h.app.FileSystemService.Mkfile(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Archivo creado correctamente",
		"item":    item,
		"data":    item,
	})
}
