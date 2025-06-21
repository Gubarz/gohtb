package errutil

import (
	"errors"
	"fmt"
)

type APIError struct {
	StatusCode int
	Message    string
	Raw        []byte
	Err        error
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("status %d: %v", e.StatusCode, e.Err)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func UnwrapFailure[T any](err error, raw []byte, status int, constructor func([]byte) T) (T, *APIError) {
	if err != nil {
		return constructor(raw), &APIError{
			StatusCode: status,
			Message:    "Request failed",
			Raw:        raw,
			Err:        err,
		}
	}
	switch status {
	case 401:
		return constructor(raw), &APIError{
			StatusCode: status,
			Message:    "Unauthorized",
			Raw:        raw,
			Err:        errors.New("unauthorized"),
		}
	case 403:
		return constructor(raw), &APIError{
			StatusCode: status,
			Message:    "Forbidden",
			Raw:        raw,
			Err:        errors.New("forbidden"),
		}
	case 429:
		return constructor(raw), &APIError{
			StatusCode: status,
			Message:    "Rate limit exceeded",
			Raw:        raw,
			Err:        errors.New("rate limit exceeded"),
		}
	}
	return constructor(raw), &APIError{
		StatusCode: status,
		Message:    "Unknown error",
		Raw:        raw,
		Err:        fmt.Errorf("unknown error: %d", status),
	}
}
