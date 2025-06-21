package common

import (
	"reflect"

	"github.com/gubarz/gohtb/internal/deref"
	httpclient "github.com/gubarz/gohtb/internal/httpclient/v4"
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
