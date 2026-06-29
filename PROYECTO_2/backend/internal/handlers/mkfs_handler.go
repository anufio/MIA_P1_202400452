package handlers

//mkfs_handler.go: contiene el controlador para manejar las rutas relacionadas con la creación de sistemas de archivos, incluyendo la formateación de particiones.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type MkfsHandler struct {
	app *services.App
}

func NewMkfsHandler(app *services.App) *MkfsHandler {
	return &MkfsHandler{
		app: app,
	}
}

func (h *MkfsHandler) Mkfs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.MkfsInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	partition, err := h.app.DiskService.FormatPartition(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"message":   "Partición formateada correctamente",
		"partition": partition,
		"data":      partition,
	})
}
