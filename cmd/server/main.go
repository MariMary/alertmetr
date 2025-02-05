package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"flag"

	"github.com/MariMary/alertmetr/internal/handlers"
	"github.com/MariMary/alertmetr/internal/storage"
	"github.com/go-chi/chi/v5"
)

var srvHandler = handlers.MetricHandlers{
	Storage: storage.NewMemStorage(),
}

type NetAddress struct {
	Host string
	Port int
}

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

func main() {
	addr := NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	addrEnv := os.Getenv("ADDRESS")
	if addr.Set(addrEnv) != nil {
		flag.Var(&addr, "a", "Net address host:port")
		flag.Parse()
	}
	log.Println(addr.String())
	r := chi.NewMux()
	r.Handle("/update/*", http.HandlerFunc(srvHandler.UpdateHandler))
	r.Handle("/value/*", http.HandlerFunc(srvHandler.GetSingleValueHandler))
	r.Handle("/", http.HandlerFunc(srvHandler.GetAllValuesHandler))

	err := http.ListenAndServe(addr.String(), r)
	if err != nil {
		panic(err)
	}

}
