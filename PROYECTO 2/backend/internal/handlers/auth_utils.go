package handlers

//auth_utils.go: contiene funciones auxiliares para manejar la autenticación, incluyendo la extracción del token de autenticación de la solicitud HTTP.

import (
	"net/http"
	"strings"
)

func tokenFromRequest(r *http.Request) string {
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if token != "" {
		return token
	}

	token = strings.TrimSpace(r.URL.Query().Get("token"))

	return token
}
