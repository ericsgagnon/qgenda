package meta

import (
	"reflect"
)

type Fields []Field

// and returns the slice of Fields, including unexported fields
// ToFields takes any struct type, slice of structs, or pointer to structs
func ToFields(value any) Fields {
	var fields Fields

	rv, rt, _ := ToIndirectReflectValue(value)
	if !rv.IsValid() {
		return fields
	}

	if rt.Kind() != reflect.Struct {
		return fields
	}

	sfs := reflect.VisibleFields(rt)
	for _, sf := range sfs {
		rfv, rft, rfPointer := ToIndirectReflectValue(rv.FieldByName(sf.Name))
		s, _ := ToStruct(rfv)

		sf.Type = rft
		field := Field{
			Name:        sf.Name,
			StructField: sf,
			Value:       rfv,
			pointer:     rfPointer,
			Struct:      s,
		}
		fields = append(fields, field)
	}
	return fields
}

// // and returns the slice of Fields, including unexported fields
// // ToFields takes any struct type, slice of structs, or pointer to structs
// func ToFields(value any) Fields {
// 	var fields Fields
// 	// var rt reflect.Type
// 	// var rv reflect.Value

// 	// switch v := value.(type) {
// 	// case nil:
// 	// 	return fields
// 	// case reflect.Type:
// 	// 	rv = reflect.New(v).Elem()
// 	// case reflect.Value:
// 	// 	rv = v
// 	// default:
// 	// 	rv = reflect.ValueOf(v)
// 	// }

// 	// // var pointer bool
// 	// switch {
// 	// case rv.Kind() == reflect.Invalid:
// 	// 	return fields
// 	// case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
// 	// 	// pointer = true
// 	// 	rt = rv.Type().Elem()
// 	// 	rv = reflect.New(rt).Elem()
// 	// case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
// 	// 	// pointer = true
// 	// 	rt = rv.Type().Elem()
// 	// 	rv = rv.Elem()
// 	// default:
// 	// 	rt = rv.Type()
// 	// }

// 	rv, rt, _ := ToIndirectReflectValue(value)
// 	if !rv.IsValid() {
// 		return fields
// 	}

// 	if rt.Kind() != reflect.Struct {
// 		return fields
// 	}

// 	sfs := reflect.VisibleFields(rt)
// 	for _, sf := range sfs {
// 		// for i := 0; i < len(sfs); i++ {
// 		var pointer bool
// 		// sf := sfs[i]
// 		rfv := rv.FieldByName(sf.Name)

// 		if !rfv.IsValid() {
// 			continue
// 		}
// 		if rfv.Kind() == reflect.Pointer {
// 			pointer = true
// 			rfv = rfv.Elem()
// 		}
// 		if sf.Type.Kind() == reflect.Pointer {
// 			pointer = true
// 			sf.Type = sf.Type.Elem()
// 		}

// 		field := Field{
// 			Name:        sf.Name,
// 			StructField: sf,
// 			Value:       rfv,
// 			pointer:     pointer,
// 		}
// 		fields = append(fields, field)
// 	}
// 	return fields
// }

// for i, field := range fields {
// 	if cfg, ok := config[field.Name]; ok {
// if cfg.Name != ""
// 	}
// }
// 	return nil
// }

// WithTag returns a subset of Fields with the key
func (f Fields) WithTag(key string) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTag(key) {
			fields = append(fields, field)
		}
	}
	return fields
}

// WithTagValue returns a subset of Fields with both the key and value
func (f Fields) WithTagValue(key, value string) Fields {
	fields := Fields{}
	for _, field := range f {
		if field.HasTagValue(key, value) {
			fields = append(fields, field)
		}
	}
	return fields

}

// WithoutTag returns a subset of Fields that do not have the key
func (f Fields) WithoutTag(key string) Fields {
	fields := Fields{}
	for _, field := range f {
		if !field.HasTag(key) {
			fields = append(fields, field)
		}
	}

	return fields
}

// WithoutTagValue returns a subset of Fields that do not have both the key and value
func (f Fields) WithoutTagValue(key, value string) Fields {
	fields := Fields{}
	for _, field := range f {
		if !field.HasTagValue(key, value) {
			fields = append(fields, field)
		}
	}

	return fields
}

// Names returns a slice of field names
func (f Fields) Names() []string {
	names := []string{}
	for _, field := range f {
		names = append(names, field.Name)
	}
	return names
}

func (f Fields) SetUUID(id string) {
	for i, _ := range f {
		f[i].SetUUID(id)
	}
}
