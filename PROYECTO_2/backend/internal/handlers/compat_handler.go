package handlers

//compat_handler.go: contiene el controlador para manejar las rutas de compatibilidad, incluyendo la eliminación de discos y la creación/listado de particiones.

import (
	"net/http"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type CompatHandler struct {
	app *services.App
}

func NewCompatHandler(app *services.App) *CompatHandler {
	return &CompatHandler{
		app: app,
	}
}

func (h *CompatHandler) DiskDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		respond.MethodNotAllowed(w)
		return
	}

	var input DeleteDiskRequest

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if err := h.app.DiskService.DeleteDisk(input.Path); err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Disco eliminado correctamente",
	})
}

func (h *CompatHandler) PartitionCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.CreatePartitionInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	partition, err := h.app.DiskService.CreatePartition(input)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusCreated, map[string]interface{}{
		"success":   true,
		"message":   "Partición creada correctamente",
		"partition": partition,
		"data":      partition,
	})
}

func (h *CompatHandler) PartitionList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.MethodNotAllowed(w)
		return
	}

	diskPath := r.URL.Query().Get("diskPath")
	if diskPath == "" {
		diskPath = r.URL.Query().Get("path")
	}

	if diskPath == "" {
		respond.BadRequest(w, "Debe indicar diskPath o path")
		return
	}

	partitions, err := h.app.DiskService.ListPartitions(diskPath)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"message":    "Particiones obtenidas correctamente",
		"partitions": partitions,
		"data":       partitions,
	})
}
