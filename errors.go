package gohtb

import (
	"errors"

	"github.com/gubarz/gohtb/internal/errutil"
)

type APIError = errutil.APIError

var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")
var ErrRateLimited = errors.New("rate limited")

func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	ok := errors.As(err, &apiErr)
	return apiErr, ok
}
