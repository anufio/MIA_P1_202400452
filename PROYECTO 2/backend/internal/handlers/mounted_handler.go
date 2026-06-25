package handlers

//mounted_handler.go: contiene el controlador para manejar las rutas relacionadas con las particiones montadas, incluyendo la obtención de la lista de particiones montadas.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type MountedHandler struct {
	app *services.App
}

func NewMountedHandler(app *services.App) *MountedHandler {
	return &MountedHandler{
		app: app,
	}
}

func (h *MountedHandler) Mounted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.MethodNotAllowed(w)
		return
	}

	mounted := h.app.DiskService.ListMounted()

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Particiones montadas obtenidas correctamente",
		"mounted": mounted,
		"data":    mounted,
	})
}
