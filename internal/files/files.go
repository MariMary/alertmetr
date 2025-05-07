package files

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/MariMary/alertmetr/internal/metric"
)

func SaveToFile(Path string, Metrics []*metric.Metric) error {
	data, err := json.MarshalIndent(Metrics, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path, data, 0666)
}

func LoadFromFile(Path string) (Metrics []*metric.Metric, err error) {
	data, err := os.ReadFile(Path)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	if err := json.NewDecoder(reader).Decode(&Metrics); err != nil {
		return nil, err
	}
	return
}
