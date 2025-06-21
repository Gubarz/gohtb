package deref

import (
	"time"

	"github.com/oapi-codegen/runtime/types"
)

func String(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func Int(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func Float32(f *float32) float32 {
	if f != nil {
		return *f
	}
	return 0
}

func Bool(p *bool) bool {
	if p != nil {
		return *p
	}
	return false
}

func Time(p *time.Time) time.Time {
	if p != nil {
		return *p
	}
	return time.Time{}
}

func Slice[T any](s *[]T) []T {
	if s == nil {
		return nil
	}
	return *s
}

func TimeFromDate(d *types.Date) time.Time {
	if d == nil {
		return time.Time{}
	}
	return d.Time
}
