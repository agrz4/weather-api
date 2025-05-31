package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.addr,
		"version": "1.1.0",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.logger.Errorw("Failed to encode health check response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
