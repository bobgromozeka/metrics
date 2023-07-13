package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/bobgromozeka/metrics/internal/server/storage"
)

func Updates(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestMetrics []metrics.RequestPayload

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&requestMetrics); err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		metricsMap := metricsArrToMaps(requestMetrics)
		cErr := s.AddCounters(r.Context(), metricsMap.Counter)
		gErr := s.SetGauges(r.Context(), metricsMap.Gauge)

		if gErr != nil || cErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}
}

func metricsArrToMaps(arr []metrics.RequestPayload) storage.Metrics {
	m := storage.Metrics{
		Gauge:   storage.GaugeMetrics{},
		Counter: storage.CounterMetrics{},
	}

	for _, payload := range arr {
		if !metrics.IsValidType(payload.MType) {
			continue
		}

		if payload.MType == metrics.CounterType {
			var delta int64
			if payload.Delta == nil {
				delta = 0
			} else {
				delta = *payload.Delta
			}

			m.Counter[payload.ID] += delta
		} else if payload.MType == metrics.GaugeType {
			var value float64
			if payload.Value == nil {
				value = 0
			} else {
				value = *payload.Value
			}
			m.Gauge[payload.ID] = value
		}
	}

	return m
}
