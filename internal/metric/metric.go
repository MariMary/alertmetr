package metric

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type Metrics struct {
	MemStats    runtime.MemStats
	PollCount   int64
	RandomValue float64
}

type MetricCollector struct {
	Metric Metrics
}

func NewMetricCollector() *MetricCollector {
	return &MetricCollector{
		Metric: Metrics{},
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
		SendMetric(MetricType, MetricName, MetricValue)
	}
}

func SendMetric(metricType string, metricName string, metricValue string) error {
	client := &http.Client{}
	url := "http://localhost:8080/update/" + metricType + "/" + metricName + "/" + metricValue
	var body []byte
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(io.Discard, response.Body)
	response.Body.Close()
	if err != nil {
		panic(err)
	}
	return nil
}

func (mc *MetricCollector) Run() {
	counter := 0
	for {
		if counter%5 == 0 {
			mc.SendMetrics()
		}
		counter++
		mc.ReadMetrics()
		time.Sleep(2 * time.Second)
	}
}
