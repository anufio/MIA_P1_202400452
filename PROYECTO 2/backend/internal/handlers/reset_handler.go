package handlers

//reset_handler.go: contiene el controlador para manejar las rutas relacionadas con el reinicio del sistema, incluyendo la eliminación de todas las particiones y archivos.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type ResetHandler struct {
	app *services.App
}

func NewResetHandler(app *services.App) *ResetHandler {
	return &ResetHandler{
		app: app,
	}
}

func (h *ResetHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		respond.MethodNotAllowed(w)
		return
	}

	if err := h.app.ResetAll(); err != nil {
		respond.InternalError(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sistema reiniciado correctamente",
	})
}
