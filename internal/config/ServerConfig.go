package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ServerConfig struct {
	Addr           NetAddress
	Port           int
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewSrvConfig() *ServerConfig {
	addr := NetAddress{
		Host: "localhost",
		Port: 8080,
	}

	addrEnv := os.Getenv("ADDRESS")
	address := flag.String("a", "localhost:8080", "Net address host:port")

	flag.Parse()
	if addr.Set(addrEnv) != nil {
		addr.Set(*address)
	}
	return &ServerConfig{
		Addr: addr,
	}
}

func NewAgtConfig() *ServerConfig {
	addr := NetAddress{
		Host: "localhost",
		Port: 8080,
	}

	addrEnv := os.Getenv("ADDRESS")
	address := flag.String("a", "localhost:8080", "Net address host:port")
	pollEnv := os.Getenv("POLL_INTERVAL")
	poll, err := strconv.Atoi(pollEnv)
	if nil != err {
		flag.IntVar(&poll, "p", 2, "poll interval")
	}
	reportEnv := os.Getenv("REPORT_INTERVAL")
	report, err := strconv.Atoi(reportEnv)
	if nil != err {
		flag.IntVar(&report, "r", 10, "report interval")
	}
	flag.Parse()
	if addr.Set(addrEnv) != nil {
		addr.Set(*address)
	}
	return &ServerConfig{
		Addr:           addr,
		PollInterval:   time.Duration(poll),
		ReportInterval: time.Duration(report),
	}
}

type NetAddress struct {
	Host string
	Port int
}

func (a *NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a NetAddress) StringHTTP() string {
	return "http://" + a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddress) Set(s string) error {
	log.Println("set", s)
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
