package main

import "net/http"

func (h *HandlerRepo) EnableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
		return
	})
}

func (h *HandlerRepo) LoadAndSaveSession(handler http.Handler) http.Handler {
	return h.Config.SCSManager.LoadAndSave(handler)
}