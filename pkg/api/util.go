package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed writing json response", "err", err)
	}
}

func respondErr(w http.ResponseWriter, err error) {
	slog.Error("api request failed", "err", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func readJsonBody[T any](r *http.Request, v *T) error {
	return json.NewDecoder(r.Body).Decode(v)
}
