package handlers

//report_handler.go: contiene el controlador para manejar las rutas relacionadas con los reportes del sistema, incluyendo la generación de reportes de disco, bitmap, árbol de directorios y journaling.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type ReportHandler struct {
	app *services.App
}

func NewReportHandler(app *services.App) *ReportHandler {
	return &ReportHandler{
		app: app,
	}
}

func (h *ReportHandler) Reports(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)

	case http.MethodPost:
		h.Generate(w, r)

	default:
		respond.MethodNotAllowed(w)
	}
}

func (h *ReportHandler) Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.ReportInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	output, err := h.app.ReportService.Generate(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Reporte generado correctamente",
		"report":  output,
		"data":    output,
	})
}

func (h *ReportHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.MethodNotAllowed(w)
		return
	}

	reports, err := h.app.ReportService.List()
	if err != nil {
		respond.InternalError(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Reportes obtenidos correctamente",
		"reports": reports,
		"data":    reports,
	})
}
