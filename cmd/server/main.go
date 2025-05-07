package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MariMary/alertmetr/internal/config"
	"github.com/MariMary/alertmetr/internal/handlers"
	"github.com/MariMary/alertmetr/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	log.Println("server started")
	cfg := config.NewSrvConfig()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	handlers.Sugar = *logger.Sugar()

	var srvHandler = handlers.MetricHandlers{
		Storage: storage.NewMemStorage(cfg.StoreInterval, cfg.StoragePath, cfg.Restore),
	}

	r := chi.NewMux()
	r.Use(handlers.GzipMiddleware, handlers.ZapLogging)
	r.Handle("/update/*", http.HandlerFunc(srvHandler.UpdateHandler))
	r.Handle("/update/", http.HandlerFunc(srvHandler.UpdateHandlerJSON))
	r.Handle("/value/*", http.HandlerFunc(srvHandler.GetSingleValueHandler))
	r.Handle("/value/", http.HandlerFunc(srvHandler.GetSingleValueHandlerJSON))
	r.Handle("/", http.HandlerFunc(srvHandler.GetAllValuesHandler))
	go func() {
		err = http.ListenAndServe(cfg.Addr.String(), r)
		if err != nil {
			panic(err)
		}
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopSignal

	srvHandler.Storage.SaveMetrics()

}
