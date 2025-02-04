package storage

type Storage interface {
	RewriteGauge(name string, value float64)
	AppendCounter(name string, value int64)
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
