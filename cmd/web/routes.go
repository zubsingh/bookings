package main

import (
	"github.com/zubsingh/bookings/cmd/pkg/config"
	"github.com/zubsingh/bookings/cmd/pkg/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	//mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/home", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
