package server

import (
	"net/http"

	"github.com/bobgromozeka/metrics/internal/server/db"
	"github.com/bobgromozeka/metrics/internal/server/handlers"
	"github.com/bobgromozeka/metrics/internal/server/middlewares"
	"github.com/bobgromozeka/metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func new(s storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Group(
		func(r chi.Router) {
			r.Use(
				middlewares.WithLogging(
					[]string{
						"./http.log",
					},
				),
				middlewares.Gzippify,
			)
			r.Post("/update/{type}/{name}/{value}", handlers.Update(s))
			r.Get("/value/{type}/{name}", handlers.Get(s))
			r.Post("/update", handlers.UpdateJSON(s))
			r.Post("/value", handlers.GetJSON(s))
			r.Get("/", handlers.GetAll(s))
		},
	)
	r.Get("/ping", handlers.Ping)

	return r
}

func Start(startupConfig StartupConfig) error {
	db.Connect(startupConfig.DatabaseDsn)
	s := storage.New()
	s = storage.NewPersistenceStorage(
		s, storage.PersistenceSettings{
			Path:     startupConfig.FileStoragePath,
			Interval: startupConfig.StoreInterval,
			Restore:  startupConfig.Restore,
		},
	)

	server := new(s)

	return http.ListenAndServe(startupConfig.ServerAddr, server)
}
