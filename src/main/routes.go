package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Use(handlerRepo.EnableCORS)
	mux.Use(handlerRepo.LoadAndSaveSession)

	mux.Get("/status", handlerRepo.Status)
	mux.Post("/v1/add-image", handlerRepo.ImageUploadingHandler)
	mux.Post("/v1/add-user", handlerRepo.UserInfoHandler)

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
