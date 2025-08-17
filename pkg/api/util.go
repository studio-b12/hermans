package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/studio-b12/elk"
	"github.com/zekrotja/hermans/pkg/controller"
	"github.com/zekrotja/hermans/pkg/database"
)

func multiHandler(handler ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, h := range handler {
			h(w, r)
		}
	}
}

func readJsonBody[T any](r *http.Request) (v T, err error) {
	limitReader := io.LimitReader(r.Body, 1*1024*1024)
	err = json.NewDecoder(limitReader).Decode(&v)
	if err != nil {
		return v, elk.Wrap(ErrParseJsonBody, err, "failed parsing json body")
	}
	return v, err
}

func respondJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error("failed to encode response json", "err", err)
		return
	}
}

func respondErr(w http.ResponseWriter, err error) {
	if vErrs, ok := elk.As[validator.ValidationErrors](err); ok {
		resp := ValidationErrors{
			ErrorResponseModel: elk.NewError(ErrValidation, "invalid request object").
				ToResponseModel(http.StatusBadRequest),
		}
		for _, vErr := range vErrs {
			resp.ValidationErrors = append(resp.ValidationErrors, &ValidationError{
				Field:   vErr.Namespace(),
				Tag:     vErr.ActualTag(),
				Value:   vErr.Value(),
				Message: vErr.Error(),
			})
		}

		respondJson(w, http.StatusBadRequest, resp)
		return
	}

	eErr := elk.Cast(err)

	switch eErr.Code() {
	case database.ErrNotFound:
		respondJson(w, http.StatusNotFound,
			eErr.ToResponseModel(http.StatusNotFound))
		return
	case ErrParseJsonBody,
		controller.ErrInvalidDips,
		controller.ErrInvalidVariants:
		respondJson(w, http.StatusBadRequest,
			eErr.ToResponseModel(http.StatusBadRequest))
		return
	}

	callFrame, _ := eErr.CallStack().First()
	slog.Error("request failed", "err", fmt.Sprintf("%v", eErr), "callFrame", callFrame)

	respondJson(w, http.StatusInternalServerError,
		eErr.ToResponseModel(http.StatusInternalServerError))
}
