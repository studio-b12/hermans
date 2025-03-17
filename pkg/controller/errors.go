package controller

import (
	"strings"

	"github.com/studio-b12/elk"
)

const (
	ErrInvalidStoreItem = elk.ErrorCode("controller:invalid-store-item")
)

type ListError []string

func (t ListError) Error() string {
	return strings.Join(t, ", ")
}
