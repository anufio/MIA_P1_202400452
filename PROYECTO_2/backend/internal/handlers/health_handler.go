package handlers

//health_handler.go: contiene el controlador para manejar las rutas de salud del sistema, incluyendo la verificación del estado del backend y la información del proyecto.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type HealthHandler struct {
	app *services.App
}

func NewHealthHandler(app *services.App) *HealthHandler {
	return &HealthHandler{
		app: app,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.MethodNotAllowed(w)
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Backend API funcionando correctamente",
		"student": "Ana Lucía Nufio Roblero",
		"carnet":  "202400452",
		"project": "MIA Proyecto Fase 2",
		"mode":    "api-rest",
	})
}
