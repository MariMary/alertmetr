package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MariMary/alertmetr/internal/storage"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/assert"
)

func TestHandlers_UpdateHandler(t *testing.T) {
	type want struct {
		method string
		code   int
		path   string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				path: "/update/counter/somemetric/800",
				code: 200,
			},
		},
		{
			name: "positive test #2",
			want: want{
				path: "/update/gauge/fmetric/4.27",
				code: 200,
			},
		},
		{
			name: "negative test #1",
			want: want{
				path: "/update/counter/34",
				code: 404,
			},
		},
		{
			name: "negative test #2",
			want: want{
				path: "/update/counter/mnn/1rrr00",
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srvHandler := MetricHandlers{
				Storage: storage.NewMemStorage(0, "metrics.txt", false),
			}
			r := httptest.NewRequest(http.MethodPost, test.want.path, nil)
			w := httptest.NewRecorder()
			srvHandler.UpdateHandler(w, r)
			result := w.Result()
			assert.Equal(t, result.StatusCode, test.want.code)
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			result.Body.Close()
		})
	}
}

func TestHandlers_GetSingleValueHandler(t *testing.T) {
	type want struct {
		method string
		code   int
		path   string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				path: "/value/counter/val208",
				code: 404,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srvHandler := MetricHandlers{
				Storage: storage.NewMemStorage(0, "metrics.txt", false),
			}
			r := httptest.NewRequest(http.MethodGet, test.want.path, nil)
			w := httptest.NewRecorder()
			srvHandler.GetSingleValueHandler(w, r)
			result := w.Result()
			assert.Equal(t, result.StatusCode, test.want.code)
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			result.Body.Close()
		})
	}
}
