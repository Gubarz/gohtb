package extract

import "reflect"

func Raw(resp interface{}) []byte {
	if resp == nil {
		return nil
	}

	v := reflect.ValueOf(resp)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if !v.IsValid() || v.Kind() != reflect.Struct {
		return nil
	}

	field := v.FieldByName("Body")
	if !field.IsValid() || !field.CanInterface() {
		return nil
	}

	if field.Kind() != reflect.Slice || field.Type().Elem().Kind() != reflect.Uint8 {
		return nil
	}

	data, ok := field.Interface().([]byte)
	if !ok {
		return nil
	}

	return data
}
