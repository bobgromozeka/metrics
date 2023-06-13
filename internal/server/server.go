package server

import (
	"github.com/bobgromozeka/metrics/internal/server/handlers"
	"github.com/bobgromozeka/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func new(s storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Post("/update/{type}/{name}/{value}", handlers.Update(s))
	r.Get("/value/{type}/{name}", handlers.Get(s))
	r.Get("/", handlers.GetAll(s))

	return r
}

func Start(serverAddr string) error {

	s := storage.New()
	server := new(s)

	return http.ListenAndServe(serverAddr, server)
}