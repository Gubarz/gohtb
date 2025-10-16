package common

import (
	"reflect"

	"github.com/microcosm-cc/bluemonday"
)

func SafeStatus(resp any) int {
	switch r := resp.(type) {
	case interface{ StatusCode() int }:
		// Check if underlying value is nil
		if reflect.ValueOf(r).IsNil() {
			return -1
		}
		return r.StatusCode()
	default:
		return -1
	}
}

var (
	strictPolicy   = bluemonday.StrictPolicy()
	sanitizePolicy = bluemonday.UGCPolicy()
)

func StrictHTML(input string) string {
	if input == "" {
		return ""
	}
	return strictPolicy.Sanitize(input)
}

func SanitizeHTML(input string) string {
	if input == "" {
		return ""
	}
	return sanitizePolicy.Sanitize(input)
}
