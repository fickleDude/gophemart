package main

import (
	"net/http"

	"github.com/fickleDude/gophemart/internal/config"
	"github.com/fickleDude/gophemart/internal/config/db"
	"github.com/fickleDude/gophemart/internal/handler"
	"github.com/fickleDude/gophemart/internal/repository"
	"github.com/fickleDude/gophemart/internal/service"
	"github.com/go-chi/chi"
)

func main() {
	//repository
	cfg := config.GetConfig()
	storage := db.GetDBConnection(cfg.DatabaseURI())
	defer db.CloseDBConnection()
	interApiRepository := repository.NewInternalApiRepository(storage)
	//services
	internalApiService := service.NewInternalApiService(interApiRepository)

	//handlers
	internalApiHandler := handler.NewInternalApiHandler(internalApiService)

	r := chi.NewRouter()
	r.Get("/api/orders/{number}", internalApiHandler.GetData)

	//start server
	err := http.ListenAndServe(cfg.AccrualSystenAddress(), r)
	if err != nil {
		panic(err)
	}
}
