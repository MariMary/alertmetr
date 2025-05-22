package metric

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetric_MetricCollector_SendMetric(t *testing.T) {
	type want struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negativee test #1",
			want: want{
				metricType:  "counter",
				metricName:  "testmetric",
				metricValue: "=====123",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metricCollector := NewMetricCollector()
			err := metricCollector.SendMetric(test.want.metricType, test.want.metricName, test.want.metricValue)
			require.Error(t, err)

		})
	}
}
