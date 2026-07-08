package utils

import "errors"

var (
	ErrNotFound            = errors.New("no content found")
	ErrInternalServerError = errors.New("internal server error")
)
