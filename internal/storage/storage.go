package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/MariMary/alertmetr/internal/files"
	"github.com/MariMary/alertmetr/internal/metric"
)

type Storage interface {
	RewriteGauge(name string, value float64)
	AppendCounter(name string, value int64)
	GetMetric(MetricType string, MetricName string) (Value string, err error)
	GetAllMetrics() (Metrics map[string]string)
	GetMetricJSON(Metric *metric.Metric) (MetricResult *metric.Metric, err error)
	RunFileStorageUpdatesAsync()
	SaveMetrics()
}

type MemStorage struct {
	GaugeMap      map[string]float64
	CounterMap    map[string]int64
	GaugeMutex    sync.RWMutex
	CounterMutex  sync.RWMutex
	StoreInterval time.Duration
	FilePath      string
}

func NewMemStorage(StoreInterval time.Duration, FilePath string, Restore bool) *MemStorage {
	ms := MemStorage{
		GaugeMap:      make(map[string]float64),
		CounterMap:    make(map[string]int64),
		GaugeMutex:    sync.RWMutex{},
		CounterMutex:  sync.RWMutex{},
		StoreInterval: StoreInterval,
		FilePath:      FilePath,
	}
	if Restore {
		ms.LoadMetrics()
	}
	go ms.RunFileStorageUpdatesAsync()
	return &ms
}

func (ms *MemStorage) SaveMetrics() {
	metrics := ms.GetAllMetricsJSON()
	files.SaveToFile(ms.FilePath, metrics)
}

func (ms *MemStorage) LoadMetrics() {
	metrics, err := files.LoadFromFile(ms.FilePath)
	if nil == err {
		for _, m := range metrics {
			switch m.MType {
			case "gauge":
				ms.RewriteGauge(m.ID, *m.Value)

			case "counter":
				ms.AppendCounter(m.ID, *m.Delta)
			}
		}
	}
}

func (ms *MemStorage) RunFileStorageUpdatesAsync() {
	if ms.StoreInterval != 0 {
		saveTick := time.NewTicker(ms.StoreInterval)
		for range saveTick.C {
			ms.SaveMetrics()
		}
	}
}

func (ms *MemStorage) GetMetric(MetricType string, MetricName string) (Value string, err error) {
	if MetricType == "gauge" {
		ms.GaugeMutex.RLock()
		value, ok := ms.GaugeMap[MetricName]
		ms.GaugeMutex.RUnlock()
		if ok {
			return fmt.Sprint(value), nil
		} else {
			return "", errors.New("unknown metric")
		}
	} else if MetricType == "counter" {
		ms.CounterMutex.RLock()
		value, ok := ms.CounterMap[MetricName]
		ms.CounterMutex.RUnlock()
		if ok {
			return fmt.Sprint(value), nil
		} else {
			return "", errors.New("unknown metric")
		}
	} else {
		return "", errors.New("unknown metric")
	}

}

func (ms *MemStorage) GetMetricJSON(Metric *metric.Metric) (MetricResult *metric.Metric, err error) {
	switch Metric.MType {
	case "gauge":
		ms.GaugeMutex.RLock()
		value, ok := ms.GaugeMap[Metric.ID]
		ms.GaugeMutex.RUnlock()
		if ok {
			Metric.Value = &value
			return Metric, nil
		} else {
			return Metric, errors.New("unknown metric")
		}

	case "counter":
		ms.CounterMutex.RLock()
		value, ok := ms.CounterMap[Metric.ID]
		ms.CounterMutex.RUnlock()
		if ok {
			Metric.Delta = &value
			return Metric, nil
		} else {
			return Metric, errors.New("unknown metric")
		}
	default:
		return Metric, errors.New("unknown metric")
	}

}

func (ms *MemStorage) GetAllMetrics() (Metrics map[string]string) {
	Metrics = make(map[string]string)
	ms.GaugeMutex.RLock()
	for name, value := range ms.GaugeMap {
		Metrics[name] = fmt.Sprint(value)
	}
	ms.GaugeMutex.RUnlock()
	ms.CounterMutex.RLock()
	for name, value := range ms.CounterMap {
		Metrics[name] = fmt.Sprint(value)
	}
	ms.CounterMutex.RUnlock()
	return
}

func (ms *MemStorage) GetAllMetricsJSON() (Metrics []*metric.Metric) {
	ms.GaugeMutex.RLock()
	for name, value := range ms.GaugeMap {
		Metrics = append(Metrics, &metric.Metric{ID: name, MType: "gauge", Value: &value})
	}
	ms.GaugeMutex.RUnlock()
	ms.CounterMutex.RLock()
	for name, value := range ms.CounterMap {
		Metrics = append(Metrics, &metric.Metric{ID: name, MType: "counter", Delta: &value})
	}
	ms.CounterMutex.RUnlock()
	return
}

func (ms *MemStorage) RewriteGauge(name string, value float64) {
	ms.GaugeMutex.Lock()
	ms.GaugeMap[name] = value
	ms.GaugeMutex.Unlock()
	if ms.StoreInterval == 0 {
		ms.SaveMetrics()
	}
}

func (ms *MemStorage) AppendCounter(name string, value int64) {
	ms.CounterMutex.Lock()
	_, ok := ms.CounterMap[name]
	if ok {
		ms.CounterMap[name] += value
	} else {
		ms.CounterMap[name] = value
	}
	ms.CounterMutex.Unlock()
	if ms.StoreInterval == 0 {
		ms.SaveMetrics()
	}

}
