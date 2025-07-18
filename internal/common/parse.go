package common

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
)

func Parse[T any](
	resp *http.Response,
	parse func(*http.Response) (*T, error),
) (parsed *T, meta ResponseMeta, err error) {
	raw := extract.Raw(resp)

	var cfRay string
	var headers http.Header
	if resp != nil && resp.Header != nil {
		cfRay = resp.Header.Get("CF-Ray")
		headers = resp.Header
	}
	meta = ResponseMeta{
		Raw:        raw,
		StatusCode: SafeStatus(resp),
		Headers:    headers,
		CFRay:      cfRay,
	}

	if resp == nil {
		parsed, err = errutil.UnwrapFailure(errors.New("nil HTTP response"), raw, meta.StatusCode, func([]byte) *T { return nil })
		return parsed, meta, err
	}

	parsed, err = parse(resp)
	if err != nil {
		parsed, err = errutil.UnwrapFailure(err, raw, meta.StatusCode, func([]byte) *T { return nil })
		return parsed, meta, err
	}

	if parsed == nil {
		parsed, err = errutil.UnwrapFailure(errors.New("parsed response is nil"), raw, meta.StatusCode, func([]byte) *T { return nil })
		return parsed, meta, err
	}

	val := reflect.ValueOf(parsed)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		var jsonFound bool
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			if strings.HasPrefix(field.Name, "JSON") {
				jsonField := val.Field(i)
				if jsonField.IsValid() && jsonField.Kind() == reflect.Ptr && !jsonField.IsNil() {
					jsonFound = true

					if field.Name != "JSON200" {
						parsed, err = errutil.UnwrapFailure(
							fmt.Errorf("%+v", jsonField.Interface()),
							raw,
							meta.StatusCode,
							func([]byte) *T { return nil },
						)
						return parsed, meta, err
					}
					break
				}
			}
		}

		if !jsonFound {
			parsed, err = errutil.UnwrapFailure(errors.New("no populated JSON* field"), raw, meta.StatusCode, func([]byte) *T { return nil })
			return parsed, meta, err
		}
	}

	return parsed, meta, nil
}
