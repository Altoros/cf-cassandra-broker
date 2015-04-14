package api

import "fmt"

type apiError struct {
	Description string `json:"description"`
}

func ApiError(format string, args ...interface{}) apiError {
	return apiError{Description: fmt.Sprintf(format, args...)}
}
