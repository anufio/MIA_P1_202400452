package handlers

//filesystem_mutation_handler.go: contiene el controlador para manejar las rutas relacionadas con la mutación del sistema de archivos, incluyendo la eliminación, edición, renombrado, copia y movimiento de archivos y directorios.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

func (h *FileSystemHandler) Remove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.FSRemoveInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if input.Token == "" {
		input.Token = tokenFromRequest(r)
	}

	if err := h.app.FileSystemService.Remove(input); err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Entrada eliminada correctamente",
	})
}

func (h *FileSystemHandler) Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.FSEditInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if input.Token == "" {
		input.Token = tokenFromRequest(r)
	}

	file, err := h.app.FileSystemService.Edit(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Archivo editado correctamente",
		"file":    file,
		"data":    file,
	})
}

func (h *FileSystemHandler) Rename(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.FSRenameInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if input.Token == "" {
		input.Token = tokenFromRequest(r)
	}

	item, err := h.app.FileSystemService.Rename(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Entrada renombrada correctamente",
		"item":    item,
		"data":    item,
	})
}

func (h *FileSystemHandler) Copy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.FSCopyInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if input.Token == "" {
		input.Token = tokenFromRequest(r)
	}

	item, err := h.app.FileSystemService.Copy(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Entrada copiada correctamente",
		"item":    item,
		"data":    item,
	})
}

func (h *FileSystemHandler) Move(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.FSMoveInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if input.Token == "" {
		input.Token = tokenFromRequest(r)
	}

	item, err := h.app.FileSystemService.Move(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Entrada movida correctamente",
		"item":    item,
		"data":    item,
	})
}
