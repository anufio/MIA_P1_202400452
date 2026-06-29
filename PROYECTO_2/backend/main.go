package main

//main.go: contiene la función principal para iniciar el servidor HTTP del backend, incluyendo la carga de la configuración, la creación del enrutador y la escucha de solicitudes entrantes.

import (
	"fmt"
	"log"
	"net/http"

	"MIA_P2_202400452/internal/api"
	"MIA_P2_202400452/internal/config"
)

func main() {
	cfg := config.Load()
	router := api.NewRouter(cfg)

	fmt.Println("========================================")
	fmt.Println(" MIA Proyecto 2 - Backend API")
	fmt.Println(" Ana Lucía Nufio Roblero - 202400452")
	fmt.Println(" Servidor: http://localhost" + cfg.ServerAddress)
	fmt.Println("========================================")

	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatal("Error iniciando servidor: ", err)
	}
}
