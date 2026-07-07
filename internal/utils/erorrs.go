package utils

import "fmt"

var (
	ErrNoContent           error = fmt.Errorf("no content found")
	ErrInternalServerError error = fmt.Errorf("internal server error")
)
