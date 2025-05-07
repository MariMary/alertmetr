package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/MariMary/alertmetr/internal/metric"
	"github.com/MariMary/alertmetr/internal/storage"
)

type MetricHandlers struct {
	Storage storage.Storage
}

func (ms *MetricHandlers) GetSingleValueHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var metric metric.Metric
	if err := json.Unmarshal(body, &metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rMetric, err := ms.Storage.GetMetricJSON(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resp, err := json.Marshal(rMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ms *MetricHandlers) GetSingleValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusNotFound)
	}
	pth := r.URL.Path
	params := strings.Split(pth, "/")
	if len(params) < 4 {
		http.Error(w, "No such metric", http.StatusNotFound)
		return
	}

	mType := params[2]
	mName := params[3]
	metric, err := ms.Storage.GetMetric(mType, mName)
	if nil != err {
		http.Error(w, "No such metric", http.StatusNotFound)
		return
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(metric))
		if nil != err {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (ms *MetricHandlers) GetAllValuesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
	metrics := ms.Storage.GetAllMetrics()
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.New("data").Parse(`<!doctype html>
	<html>
	<head>
		<title> List of metrics </title>
	</head>
	<body>
		<table>
			<thead>
				<tr>
					<th>Name</th><th>Value</th>
				</tr>
			</thead>
			<tbody>
			{{- range $key, $value := . -}}
				<tr><td>{{- $key -}}</td><td>{{- $value -}}</td></tr>
			{{end}}
			</tbody>
		</table>
	</body>
	</html>`))
	tmpl.Execute(w, metrics)

}

func (ms *MetricHandlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}

	pth := r.URL.Path
	re := regexp.MustCompile(`^/update/(.+)$`)
	if !re.MatchString(pth) {
		http.Error(w, "Url not found", http.StatusNotFound)
		return
	}

	params := strings.Split(pth, "/")
	if len(params) < 5 {
		http.Error(w, "No metric name", http.StatusNotFound)
		return
	}
	mType := params[2]
	mName := params[3]
	mValue := params[4]
	if strings.Contains(mType, "gauge") {
		gaugeValue, er := strconv.ParseFloat(mValue, 64)
		if nil != er {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
		ms.Storage.RewriteGauge(mName, gaugeValue)

	} else if strings.Contains(mType, "counter") {
		counterValue, er := strconv.ParseInt(mValue, 10, 64)
		if nil != er {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}
		ms.Storage.AppendCounter(mName, int64(counterValue))
	} else {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (ms *MetricHandlers) UpdateHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var metric metric.Metric
	if err := json.Unmarshal(body, &metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch metric.MType {
	case "gauge":
		ms.Storage.RewriteGauge(metric.ID, *metric.Value)
	case "counter":
		ms.Storage.AppendCounter(metric.ID, *metric.Delta)
	default:
		http.Error(w, "invalid metric type", http.StatusInternalServerError)
	}

	rMetric, err := ms.Storage.GetMetricJSON(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(rMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
