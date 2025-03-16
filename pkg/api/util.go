package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/studio-b12/elk"
)

func respondJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error("failed to encode response json", err)
		return
	}
}

func respondErr(w http.ResponseWriter, err error) {
	if vErr, ok := elk.As[validator.ValidationErrors](err); ok {
		// FIXME: better wrapping
		respondJson(w, http.StatusBadRequest, vErr)
		return
	}

	m := elk.Cast(err).ToResponseModel(http.StatusInternalServerError)
	respondJson(w, http.StatusInternalServerError, m)
}
