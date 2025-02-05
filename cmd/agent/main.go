package main

import "github.com/MariMary/alertmetr/internal/metric"

func main() {
	metricCollector := metric.NewMetricCollector()
	metricCollector.Run()
}
