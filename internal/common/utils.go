package common

import (
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
