package metric

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/MariMary/alertmetr/internal/client"
	"github.com/MariMary/alertmetr/internal/config"
)

type Metrics struct {
	MemStats    runtime.MemStats
	PollCount   int64
	RandomValue float64
}

type MetricCollector struct {
	Cfg        *config.ServerConfig
	Metric     Metrics
	HttpClient *client.HttpClient
}

func NewMetricCollector() *MetricCollector {
	cfg := config.NewAgtConfig()
	return &MetricCollector{
		Cfg:        cfg,
		HttpClient: client.NewHttpClient(cfg.Addr.StringHttp()),
		Metric:     Metrics{},
	}
}

func (mc *MetricCollector) ReadMetrics() {
	runtime.ReadMemStats(&mc.Metric.MemStats)
	mc.Metric.PollCount += 1
	mc.Metric.RandomValue = rand.Float64()
}

func (mc *MetricCollector) SendMetrics() {
	values := reflect.ValueOf(mc.Metric.MemStats)
	typs := reflect.TypeOf(mc.Metric.MemStats)
	for i := 0; i < values.NumField(); i++ {
		MetricValType := typs.Field(i).Type.Name()
		MetricName := typs.Field(i).Name
		MetricValue := ""
		MetricType := ""
		if MetricValType == "float64" {
			value := reflect.ValueOf(values.Field(i).Interface()).Float()
			MetricValue = fmt.Sprintf("%v", value)
			MetricType = "gauge"
		} else if MetricValType == "uint64" {
			value := reflect.ValueOf(values.Field(i).Interface()).Uint()
			MetricValue = fmt.Sprintf("%v", value)
			MetricType = "counter"
		} else if MetricValType == "int64" {
			value := reflect.ValueOf(values.Field(i).Interface()).Int()
			MetricValue = fmt.Sprintf("%v", value)
			MetricType = "counter"
		}
		mc.SendMetric(MetricType, MetricName, MetricValue)
	}
}

func (mc *MetricCollector) SendMetric(metricType string, metricName string, metricValue string) error {

	ApiName := "/update/" + metricType + "/" + metricName + "/" + metricValue
	return mc.HttpClient.CallApi(ApiName)
}

func (mc *MetricCollector) Run() {
	pollTick := time.NewTicker(mc.Cfg.PollInterval)
	reportTick := time.NewTicker(mc.Cfg.ReportInterval)
	for {
		select {
		case <-pollTick.C:
			mc.ReadMetrics()
		case <-reportTick.C:
			mc.SendMetrics()
		}

	}
}
