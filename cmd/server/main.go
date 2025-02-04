package main

import (
	"net/http"

	"github.com/MariMary/alertmetr/internal/handlers"
	"github.com/MariMary/alertmetr/internal/storage"
)

var srvHandler = handlers.MetricHandlers{
	Storage: storage.NewMemStorage(),
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(srvHandler.UpdateHandler))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
