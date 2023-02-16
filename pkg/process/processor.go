package process

import (
	"reflect"

	"github.com/ericsgagnon/qgenda/pkg/meta"
)

var IDTag string = "id"

type ValueMap map[string]any

func ToValueMap(value any) ValueMap {

	
	vm := ValueMap{}
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil
	}
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		fields := meta.ToFields(value)
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			var v any
			if field.IsValid() {
				v = field.Interface()
			}
			vm[fields[i].Name] = v
		}
		return vm
	}
	return nil
}

func IDHashFields(value any) meta.Fields {
	return nil
}

func IDHashValues(value any)

func IDHash(value any) []byte {
	return nil
}
