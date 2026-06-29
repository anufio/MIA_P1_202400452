package api

//router.go: define las rutas de la API y asocia cada ruta con su correspondiente controlador.

import (
	"net/http"

	"MIA_P2_202400452/internal/config"
	"MIA_P2_202400452/internal/handlers"
	"MIA_P2_202400452/internal/services"
)

func NewRouter(cfg config.Config) http.Handler {
	app := services.NewApp(cfg)

	mux := http.NewServeMux()

	healthHandler := handlers.NewHealthHandler(app)
	diskHandler := handlers.NewDiskHandler(app)
	partitionHandler := handlers.NewPartitionHandler(app)
	mkfsHandler := handlers.NewMkfsHandler(app)
	authHandler := handlers.NewAuthHandler(app)
	fsHandler := handlers.NewFileSystemHandler(app)
	reportHandler := handlers.NewReportHandler(app)
	mountedHandler := handlers.NewMountedHandler(app)
	compatHandler := handlers.NewCompatHandler(app)
	resetHandler := handlers.NewResetHandler(app)

	mux.HandleFunc("/api/health", healthHandler.Health)
	mux.HandleFunc("/api/reset", resetHandler.Reset)

	mux.HandleFunc("/api/disks", diskHandler.Disks)
	mux.HandleFunc("/api/disks/delete", compatHandler.DiskDelete)
	mux.HandleFunc("/api/disks/", diskHandler.DiskByPath)

	mux.HandleFunc("/api/partitions", partitionHandler.Partitions)
	mux.HandleFunc("/api/partitions/list", compatHandler.PartitionList)
	mux.HandleFunc("/api/partitions/add", compatHandler.PartitionCreate)
	mux.HandleFunc("/api/partitions/create", compatHandler.PartitionCreate)
	mux.HandleFunc("/api/partitions/", partitionHandler.PartitionAction)

	mux.HandleFunc("/api/mounted", mountedHandler.Mounted)
	mux.HandleFunc("/api/partitions/mounted", mountedHandler.Mounted)

	mux.HandleFunc("/api/mkfs", mkfsHandler.Mkfs)

	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/auth/me", authHandler.Me)

	mux.HandleFunc("/api/fs/list", fsHandler.List)
	mux.HandleFunc("/api/fs/read", fsHandler.Read)
	mux.HandleFunc("/api/fs/mkdir", fsHandler.Mkdir)
	mux.HandleFunc("/api/fs/mkfile", fsHandler.Mkfile)
	mux.HandleFunc("/api/fs/remove", fsHandler.Remove)
	mux.HandleFunc("/api/fs/edit", fsHandler.Edit)
	mux.HandleFunc("/api/fs/rename", fsHandler.Rename)
	mux.HandleFunc("/api/fs/copy", fsHandler.Copy)
	mux.HandleFunc("/api/fs/move", fsHandler.Move)
	mux.HandleFunc("/api/fs/", fsHandler.FSAction)

	mux.HandleFunc("/api/explorer/list", fsHandler.List)
	mux.HandleFunc("/api/explorer/read", fsHandler.Read)
	mux.HandleFunc("/api/explorer/mkdir", fsHandler.Mkdir)
	mux.HandleFunc("/api/explorer/mkfile", fsHandler.Mkfile)
	mux.HandleFunc("/api/explorer/remove", fsHandler.Remove)
	mux.HandleFunc("/api/explorer/edit", fsHandler.Edit)
	mux.HandleFunc("/api/explorer/rename", fsHandler.Rename)
	mux.HandleFunc("/api/explorer/copy", fsHandler.Copy)
	mux.HandleFunc("/api/explorer/move", fsHandler.Move)
	mux.HandleFunc("/api/explorer/", fsHandler.Explorer)

	mux.HandleFunc("/api/reports", reportHandler.Reports)
	mux.HandleFunc("/api/reports/list", reportHandler.List)
	mux.HandleFunc("/api/reports/generate", reportHandler.Generate)

	mux.Handle("/reports/", http.StripPrefix("/reports/", http.FileServer(http.Dir(cfg.ReportRoot))))

	return WithCORS(mux)
}
