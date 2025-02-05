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
	r := http.NewServeMux()
	r.Handle("/update/", http.HandlerFunc(srvHandler.UpdateHandler))
	r.Handle("/value/", http.HandlerFunc(srvHandler.GetSingleValueHandler))
	r.Handle("/", http.HandlerFunc(srvHandler.GetAllValuesHandler))

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}

}
