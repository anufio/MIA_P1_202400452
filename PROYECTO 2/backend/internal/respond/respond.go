package respond

//respond.go: contiene funciones auxiliares para manejar las respuestas HTTP, incluyendo la serialización de JSON y el manejo de errores.

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func OK(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, Response{
		Success: false,
		Error:   message,
	})
}

func Unauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, Response{
		Success: false,
		Error:   message,
	})
}

func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, Response{
		Success: false,
		Error:   message,
	})
}

func MethodNotAllowed(w http.ResponseWriter) {
	JSON(w, http.StatusMethodNotAllowed, Response{
		Success: false,
		Error:   "Método no permitido",
	})
}

func InternalError(w http.ResponseWriter, message string) {
	JSON(w, http.StatusInternalServerError, Response{
		Success: false,
		Error:   message,
	})
}
