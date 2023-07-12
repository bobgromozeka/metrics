package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bobgromozeka/metrics/internal/server/storage"
)

func GetAll(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gaugeMetrics, gErr := s.GetAllGaugeMetrics(r.Context())
		counterMetrics, cErr := s.GetAllCounterMetrics(r.Context())
		if gErr != nil || cErr != nil {
			log.Println(gErr, cErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := ""

		for k, v := range gaugeMetrics {
			response += fmt.Sprintf("%s:   %f\r\n", k, v)
		}

		for k, v := range counterMetrics {
			response += fmt.Sprintf("%s:   %d\n", k, v)
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}
