package db

import (
	"github.com/studio-b12/elk"
)

const (
	ErrDirectory = elk.ErrorCode("db:directory")
	ErrFile      = elk.ErrorCode("db:file")
	ErrDecode    = elk.ErrorCode("db:decode")
	ErrEncode    = elk.ErrorCode("db:encode")
	ErrDeepCopy  = elk.ErrorCode("db:deepcopy")
)
