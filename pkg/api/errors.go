package api

import (
	"github.com/studio-b12/elk"
)

const (
	ErrParseJsonBody = elk.ErrorCode("api:parse-json-body")
	ErrValidation    = elk.ErrorCode("api:validation")
)

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   any    `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	elk.ErrorResponseModel
	ValidationErrors []*ValidationError `json:"validation_errors"`
}
