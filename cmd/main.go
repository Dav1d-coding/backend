// @title My App API
// @version 1.0
// @description Это API моего сервиса
// @host localhost:8080
// @BasePath /
package main

import (
	"app/backendv1/internal/config"
	"app/backendv1/internal/delivery/http_handler"
	"app/backendv1/internal/repository/postgres"
	"app/backendv1/internal/usecase"
	"database/sql"
	"log"
	"net/http"

	_ "app/backendv1/docs"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config.LoadEnv()
	dsn := config.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB CONNECTION FAILED", err)
	}
	defer db.Close()

	if err := ensureTables(db); err != nil {
		log.Fatalf("could not create tables: %v", err)
	}
	//namespace setup
	namespaceRepo := postgres.NewNamespaceRepo(db)
	namespaceUC := usecase.NewNamespaceService(namespaceRepo)
	namespaceHandler := http_handler.NewHandler(namespaceUC)

	//app setup
	appRepo := postgres.NewAppRepo(db)
	appUC := usecase.NewAppUsecase(appRepo)
	appHandler := http_handler.NewAppHandler(appUC)

	//appData setup
	appDataRepo := postgres.NewAppDataRepo(db)
	appDataUC := usecase.NewAppDataUsecase(appDataRepo)
	appDataHandler := http_handler.NewAppDataHandler(appDataUC)

	r := mux.NewRouter()
	namespaceHandler.RegisterRoutes(r)
	appHandler.RegisterRoutes(r)
	appDataHandler.RegisterRoutes(r)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// init tables
func ensureTables(db *sql.DB) error {
	createAppsTable := `
	CREATE TABLE IF NOT EXISTS namespaces (
		code TEXT PRIMARY KEY,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS apps (
		code TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		namespace_code TEXT NOT NULL,
		icon TEXT,
		FOREIGN KEY (namespace_code) REFERENCES namespaces(code) ON DELETE CASCADE
	);`

	_, err := db.Exec(createAppsTable)
	return err
}
