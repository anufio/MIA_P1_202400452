package handlers

//partition_handler.go: contiene el controlador para manejar las rutas relacionadas con las particiones, incluyendo la creación, eliminación, montaje, desmontaje y redimensionamiento de particiones.

import (
	"net/http"
	"strings"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type PartitionHandler struct {
	app *services.App
}

func NewPartitionHandler(app *services.App) *PartitionHandler {
	return &PartitionHandler{
		app: app,
	}
}

func (h *PartitionHandler) Partitions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
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

	case http.MethodGet:
		mounted := h.app.DiskService.ListMounted()

		respond.JSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Particiones montadas obtenidas correctamente",
			"mounted": mounted,
			"data":    mounted,
		})

	default:
		respond.MethodNotAllowed(w)
	}
}

func (h *PartitionHandler) PartitionAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	action := strings.TrimPrefix(r.URL.Path, "/api/partitions/")
	action = strings.Trim(action, "/")

	switch action {
	case "mount":
		var input services.MountInput

		if err := readJSON(r, &input); err != nil {
			respond.BadRequest(w, "JSON inválido")
			return
		}

		mounted, err := h.app.DiskService.Mount(input)
		if err != nil {
			respond.BadRequest(w, err.Error())
			return
		}

		respond.JSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Partición montada correctamente",
			"mounted": mounted,
			"data":    mounted,
		})

	case "unmount":
		var input services.UnmountInput

		if err := readJSON(r, &input); err != nil {
			respond.BadRequest(w, "JSON inválido")
			return
		}

		if err := h.app.DiskService.Unmount(input); err != nil {
			respond.BadRequest(w, err.Error())
			return
		}

		respond.JSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Partición desmontada correctamente",
		})

	case "delete":
		var input services.DeletePartitionInput

		if err := readJSON(r, &input); err != nil {
			respond.BadRequest(w, "JSON inválido")
			return
		}

		if err := h.app.DiskService.DeletePartition(input); err != nil {
			respond.BadRequest(w, err.Error())
			return
		}

		respond.JSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Partición eliminada correctamente",
		})

	case "resize":
		var input services.ResizePartitionInput

		if err := readJSON(r, &input); err != nil {
			respond.BadRequest(w, "JSON inválido")
			return
		}

		partition, err := h.app.DiskService.ResizePartition(input)
		if err != nil {
			respond.BadRequest(w, err.Error())
			return
		}

		respond.JSON(w, http.StatusOK, map[string]interface{}{
			"success":   true,
			"message":   "Partición redimensionada correctamente",
			"partition": partition,
			"data":      partition,
		})

	default:
		respond.NotFound(w, "Acción de partición no encontrada")
	}
}
