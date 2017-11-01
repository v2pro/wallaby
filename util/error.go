package util

import "errors"

var (
	EOF      = errors.New("EOF")
	NOTFOUND = errors.New("NOTFOUND")
)
