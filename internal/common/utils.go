package common

import (
	"reflect"

	httpclient "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/deref"
	"github.com/microcosm-cc/bluemonday"
)

func FromAPITodoItem(data httpclient.Item) TodoItem {
	return TodoItem{
		Id: deref.Int(data.Id),
	}
}

func FromAPIFlag(data httpclient.Flag) Flag {
	return Flag{
		Id:     deref.Int(data.Id),
		Owned:  deref.Bool(data.Owned),
		Points: deref.Int(data.Points),
		Title:  deref.String(data.Title),
	}
}

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
