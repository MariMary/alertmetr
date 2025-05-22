package storage

import (
	"errors"
	"fmt"
)

type Storage interface {
	RewriteGauge(name string, value float64)
	AppendCounter(name string, value int64)
	GetMetric(MetricType string, MetricName string) (Value string, err error)
	GetAllMetrics() (Metrics map[string]string)
}

type MemStorage struct {
	GaugeMap   map[string]float64
	CounterMap map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMap:   make(map[string]float64),
		CounterMap: make(map[string]int64),
	}
}

func (ms *MemStorage) GetMetric(MetricType string, MetricName string) (Value string, err error) {
	if MetricType == "gauge" {
		value, ok := ms.GaugeMap[MetricName]
		if ok {
			return fmt.Sprint(value), nil
		} else {
			return "", errors.New("unknown metric")
		}
	} else if MetricType == "counter" {
		value, ok := ms.CounterMap[MetricName]
		if ok {
			return fmt.Sprint(value), nil
		} else {
			return "", errors.New("unknown metric")
		}
	} else {
		return "", errors.New("unknown metric")
	}

}

func (ms *MemStorage) GetAllMetrics() (Metrics map[string]string) {
	Metrics = make(map[string]string)
	for name, value := range ms.GaugeMap {
		Metrics[name] = fmt.Sprint(value)
	}
	for name, value := range ms.CounterMap {
		Metrics[name] = fmt.Sprint(value)
	}
	return
}

func (ms *MemStorage) RewriteGauge(name string, value float64) {
	ms.GaugeMap[name] = value
}

func (ms *MemStorage) AppendCounter(name string, value int64) {
	_, ok := ms.CounterMap[name]
	if ok {
		ms.CounterMap[name] += value
	} else {
		ms.CounterMap[name] = value
	}

}
