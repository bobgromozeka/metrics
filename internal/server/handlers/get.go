package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/bobgromozeka/metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
)

func Get(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsType := chi.URLParam(r, "type")
		metricsName := chi.URLParam(r, "name")

		m, err := s.GetMetricsByType(r.Context(), metricsType, metricsName)

		w.Header().Set("Content-Type", "text/html")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", m)))
	}
}

func GetJSON(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestMetrics metrics.RequestPayload

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&requestMetrics); err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		}

		if !metrics.IsValidType(requestMetrics.MType) {
			log.Println("Got wrong metrics type in request: ", requestMetrics.MType)
			http.Error(w, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if requestMetrics.MType == metrics.CounterType {
			val, err := s.GetCounterMetrics(r.Context(), requestMetrics.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			requestMetrics.Delta = &val
		} else {
			val, err := s.GetGaugeMetrics(r.Context(), requestMetrics.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			requestMetrics.Value = &val
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(requestMetrics)
	}
}
