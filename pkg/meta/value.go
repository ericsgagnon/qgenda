package meta

import "reflect"

func ValueOf(value any) reflect.Value {

	rv := reflect.ValueOf(value)
	rt := reflect.TypeOf(value)
	if !rv.IsValid() {
		rv = reflect.New(rt).Elem()
	}
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	
	return rv
}
