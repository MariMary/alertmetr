package metric

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/MariMary/alertmetr/internal/client"
	"github.com/MariMary/alertmetr/internal/config"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Metrics struct {
	MemStats    runtime.MemStats
	PollCount   int64
	RandomValue float64
}

type MetricCollector struct {
	Cfg        *config.ServerConfig
	Metric     Metrics
	HTTPClient *client.HTTPClient
}

func NewMetricCollector() *MetricCollector {
	cfg := config.NewAgtConfig()
	return &MetricCollector{
		Cfg:        cfg,
		HTTPClient: client.NewHTTPClient(cfg.Addr.StringHTTP()),
		Metric:     Metrics{},
	}
}

func (mc *MetricCollector) ReadMetrics() {
	runtime.ReadMemStats(&mc.Metric.MemStats)
	mc.Metric.PollCount += 1
	mc.Metric.RandomValue = rand.Float64()
}

func (mc *MetricCollector) SendMetrics() {
	mc.SendMetricJSON("counter", "PollCount", nil, &mc.Metric.PollCount)
	mc.SendMetricJSON("gauge", "RandomValue", &mc.Metric.RandomValue, nil)

	values := reflect.ValueOf(mc.Metric.MemStats)
	typs := reflect.TypeOf(mc.Metric.MemStats)
	for i := 0; i < values.NumField(); i++ {

		MetricValType := typs.Field(i).Type.Name()
		MetricName := typs.Field(i).Name
		if MetricValType == "float64" {
			value := reflect.ValueOf(values.Field(i).Interface()).Float()
			mc.SendMetricJSON("gauge", MetricName, &value, nil)
		} else if MetricValType == "uint64" {
			value := reflect.ValueOf(values.Field(i).Interface()).Uint()
			v64 := float64(value)
			mc.SendMetricJSON("gauge", MetricName, &v64, nil)
		} else if MetricValType == "int64" {
			value := reflect.ValueOf(values.Field(i).Interface()).Int()
			v64 := float64(value)
			mc.SendMetricJSON("gauge", MetricName, &v64, nil)
		} else if MetricValType == "uint32" {
			value := reflect.ValueOf(values.Field(i).Interface()).Uint()
			v64 := float64(value)
			mc.SendMetricJSON("gauge", MetricName, &v64, nil)
		}
	}
}

func (mc *MetricCollector) SendMetric(metricType string, metricName string, metricValue string) error {

	APIName := "/update/" + metricType + "/" + metricName + "/" + metricValue
	return mc.HTTPClient.CallAPI(APIName, nil, "text/plain")
	//return mc.HTTPClient.CallAPI(APIName)
}

func (mc *MetricCollector) SendMetricJSON(metricType string, metricName string, Value *float64, Delta *int64) error {

	metric := Metric{
		ID:    metricName,
		MType: metricType,
		Value: Value,
		Delta: Delta,
	}
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	gz.Write(body)
	gz.Close()
	return mc.HTTPClient.CallAPIBuf("/update/", buf, "application/json")
}

func (mc *MetricCollector) Run() {
	mc.ReadMetrics()
	mc.SendMetrics()
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
