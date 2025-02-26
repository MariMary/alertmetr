package main

import (
	"log"

	"github.com/MariMary/alertmetr/internal/metric"
)

func main() {
	metricCollector := metric.NewMetricCollector()
	log.Println(metricCollector.Cfg.Addr.String())
	metricCollector.Run()
}
