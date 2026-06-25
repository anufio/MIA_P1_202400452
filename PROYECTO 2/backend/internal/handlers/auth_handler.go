package handlers

//auth_handler.go: contiene el controlador para manejar las rutas de autenticación, incluyendo login, logout y obtener información del usuario autenticado.

import (
	"net/http"
	"strings"

	"MIA_P2_202400452/internal/respond"
	"MIA_P2_202400452/internal/services"
)

type AuthHandler struct {
	app *services.App
}

func NewAuthHandler(app *services.App) *AuthHandler {
	return &AuthHandler{
		app: app,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.LoginInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	session, err := h.app.AuthService.Login(input)
	if err != nil {
		respond.Unauthorized(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login realizado correctamente",
		"session": session,
		"data":    session,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.MethodNotAllowed(w)
		return
	}

	var input services.LogoutInput

	if err := readJSON(r, &input); err != nil {
		respond.BadRequest(w, "JSON inválido")
		return
	}

	if err := h.app.AuthService.Logout(input); err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logout realizado correctamente",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.MethodNotAllowed(w)
		return
	}

	token := strings.TrimSpace(r.Header.Get("Authorization"))
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		respond.Unauthorized(w, "Token no enviado")
		return
	}

	session, ok := h.app.AuthService.GetSession(token)
	if !ok {
		respond.Unauthorized(w, "Sesión no encontrada")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sesión obtenida correctamente",
		"session": session,
		"data":    session,
	})
}
