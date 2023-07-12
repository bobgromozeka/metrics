package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/bobgromozeka/metrics/internal/server/db"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if db.Connection() == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := db.Connection().Ping(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
