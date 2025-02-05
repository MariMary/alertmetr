package metric

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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

type Metrics struct {
	MemStats    runtime.MemStats
	PollCount   int64
	RandomValue float64
}

type MetricCollector struct {
	Metric         Metrics
	Addr           NetAddress
	reportInterval time.Duration
	pollInterval   time.Duration
}

func NewMetricCollector() *MetricCollector {
	addr := NetAddress{
		Host: "http://localhost",
		Port: 8080,
	}
	addrEnv := os.Getenv("ADDRESS")
	if addr.Set(addrEnv) != nil {
		flag.Var(&addr, "a", "Net address host:port")
		flag.Parse()
	}
	pollEnv := os.Getenv("POLL_INTERVAL")
	poll, err := strconv.Atoi(pollEnv)
	if nil != err {
		pollFlag := flag.Int("p", 2, "pol interval")
		flag.Parse()
		poll = *pollFlag
	}
	reportEnv := os.Getenv("REPORT_INTERVAL")
	report, err := strconv.Atoi(reportEnv)
	if nil != err {
		reportFlag := flag.Int("r", 10, "report interval")
		flag.Parse()
		report = *reportFlag
	}

	return &MetricCollector{
		Addr:           addr,
		pollInterval:   time.Duration(poll),
		reportInterval: time.Duration(report),
		Metric:         Metrics{},
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
		SendMetric(mc.Addr.String(), MetricType, MetricName, MetricValue)
	}
}

func SendMetric(Addr string, metricType string, metricName string, metricValue string) error {
	client := &http.Client{}
	url := Addr + "/update/" + metricType + "/" + metricName + "/" + metricValue
	var body []byte
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	_, err = io.Copy(io.Discard, response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func (mc *MetricCollector) Run() {
	pollTick := time.NewTicker(mc.pollInterval)
	reportTick := time.NewTicker(mc.reportInterval)
	for {
		select {
		case <-pollTick.C:
			mc.ReadMetrics()
		case <-reportTick.C:
			mc.SendMetrics()
		}

	}
}
