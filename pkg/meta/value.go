package meta

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"reflect"
)

// ValueMap's keys should match a struct's members' names (or tags) and values
// are wrapped in any()
type ValueMap map[string]any

func ToValueMap(value any, tagKey string) ValueMap {
	var fields Fields
	switch v := value.(type) {
	case Fields:
		fields = v
	case Struct:
		fields = v.Fields
	default:
		fields = ToFields(v)
	}
	if len(fields) == 0 {
		return nil
	}

	// to reduce the amount of tags necessary, we look for:
	// - if any tag == false: include all except those fields
	// - if any tag == true:  include only those fields
	// if no fields have tagKey: we default to including all Exported fields
	switch {
	case len(fields.WithTagFalse(tagKey)) > 0:
		fields = fields.WithTagTrue(tagKey)
	case len(fields.WithTagTrue(tagKey)) > 0:
		fields = fields.WithTagTrue(tagKey)
	default:
		// do nothing - will include all fields
	}

	vm := ValueMap{}
	for _, field := range fields {
		var v any
		if field.Value.Kind() != reflect.Invalid {
			v = field.Value.Interface()
		}
		vm[field.Name] = v
	}
	return vm
}

func (v ValueMap) Hash() string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var bb bytes.Buffer
	if err := json.Compact(&bb, b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", sha1.Sum(bb.Bytes()))
}
