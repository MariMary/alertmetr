package handlers

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/MariMary/alertmetr/internal/storage"
)

type MetricHandlers struct {
	Storage storage.Storage
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
