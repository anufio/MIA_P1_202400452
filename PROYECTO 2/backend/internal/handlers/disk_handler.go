package handlers

//disk_handler.go: contiene el controlador para manejar las rutas relacionadas con los discos, incluyendo la obtención de la lista de discos, la creación y eliminación de discos, y la obtención de particiones por ruta de disco.

import (
	"encoding/json"
	"net/http"
	"strings"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type DiskHandler struct {
	app *services.App
}

type DeleteDiskRequest struct {
	Path string `json:"path"`
}

func NewDiskHandler(app *services.App) *DiskHandler {
	return &DiskHandler{
		app: app,
	}
}

func (h *DiskHandler) Disks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		disks, err := h.app.DiskService.ListDisks()
		if err != nil {
			respond.InternalError(w, err.Error())
			return
		}

		respond.JSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Discos obtenidos correctamente",
			"disks":   disks,
			"data":    disks,
		})

	case http.MethodPost:
		var input services.CreateDiskInput

		if err := readJSON(r, &input); err != nil {
			respond.BadRequest(w, "JSON inválido")
			return
		}

		diskInfo, err := h.app.DiskService.CreateDisk(input)
		if err != nil {
			respond.BadRequest(w, err.Error())
			return
		}

		respond.JSON(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Disco creado correctamente",
			"disk":    diskInfo,
			"data":    diskInfo,
		})

	case http.MethodDelete:
		var input DeleteDiskRequest

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
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

	default:
		respond.MethodNotAllowed(w)
	}
}

func (h *DiskHandler) DiskByPath(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/disks/")
	path = strings.Trim(path, "/")

	if path == "" {
		respond.NotFound(w, "Ruta no encontrada")
		return
	}

	if strings.HasSuffix(path, "/partitions") {
		diskID := strings.TrimSuffix(path, "/partitions")
		diskID = strings.Trim(diskID, "/")

		diskPath, err := services.DecodeDiskID(diskID)
		if err != nil {
			respond.BadRequest(w, "ID de disco inválido")
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
		return
	}

	respond.NotFound(w, "Ruta no encontrada")
}
