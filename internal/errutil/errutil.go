package errutil

import (
	"errors"
	"fmt"
	"strings"
)

type APIError struct {
	StatusCode int
	Message    string
	Raw        []byte
	Err        error
}

const (
	StatusUnmarshalError = 1001
)

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
		if isUnmarshalError(err) {
			return constructor(raw), &APIError{
				StatusCode: StatusUnmarshalError,
				Message:    "Failed to parse response JSON",
				Raw:        raw,
				Err:        err,
			}
		}

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

	case 500, 502, 503, 504:
		return constructor(raw), &APIError{
			StatusCode: status,
			Message:    "Server error",
			Raw:        raw,
			Err:        fmt.Errorf("server error: %d", status),
		}
	}

	return constructor(raw), &APIError{
		StatusCode: status,
		Message:    "Unknown error",
		Raw:        raw,
		Err:        fmt.Errorf("unknown error: %d", status),
	}
}

func isUnmarshalError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "json:") ||
		strings.Contains(errStr, "unmarshal")
}
