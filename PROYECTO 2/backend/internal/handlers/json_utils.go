package handlers

//json_utils.go: contiene funciones auxiliares para manejar la lectura de JSON desde la solicitud HTTP.

import (
	"encoding/json"
	"net/http"
)

func readJSON(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
