package main

import (
	"log"
	"net/http"

	"github.com/MariMary/alertmetr/internal/config"
	"github.com/MariMary/alertmetr/internal/handlers"
	"github.com/MariMary/alertmetr/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var srvHandler = handlers.MetricHandlers{
	Storage: storage.NewMemStorage(),
}

func main() {
	cfg := config.NewSrvConfig()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	handlers.Sugar = *logger.Sugar()

	r := chi.NewMux()
	r.Use(handlers.ZapLogging)
	r.Handle("/update/*", http.HandlerFunc(srvHandler.UpdateHandler))
	r.Handle("/value/*", http.HandlerFunc(srvHandler.GetSingleValueHandler))
	r.Handle("/", http.HandlerFunc(srvHandler.GetAllValuesHandler))
	log.Println(cfg.Addr.String())

	err = http.ListenAndServe(cfg.Addr.String(), r)
	if err != nil {
		panic(err)
	}

}
