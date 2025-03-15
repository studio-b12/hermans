package cache

import (
	"github.com/studio-b12/elk"
)

const (
	ErrDirectory = elk.ErrorCode("cache:directory")
	ErrFile      = elk.ErrorCode("cache:file")
	ErrDecode    = elk.ErrorCode("cache:decode")
	ErrEncode    = elk.ErrorCode("cache:encode")
	ErrDeepCopy  = elk.ErrorCode("cache:deepcopy")
)
