package controller

import (
	"strings"

	"github.com/studio-b12/elk"
)

const (
	ErrInvalidStoreItem = elk.ErrorCode("controller:invalid-store-item")
	ErrInvalidVariants  = elk.ErrorCode("controller:invalid-variants")
	ErrInvalidDips      = elk.ErrorCode("controller:invalid-dips")
	ErrInvalidEditKey   = elk.ErrorCode("controller:invalid-edit-key")
)

type ListError []string

func (t ListError) Error() string {
	return strings.Join(t, ", ")
}

func (t ListError) Details() any {
	return t
}
