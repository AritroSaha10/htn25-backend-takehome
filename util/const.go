package util

import "errors"

const (
	ISO8601 = "2006-01-02T15:04:05.99999"
)

var (
	ErrNotFound   = errors.New("resource not found")
	ErrBadRequest = errors.New("bad request")
)
