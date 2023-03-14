package meta

import (
	"reflect"
)

// type Structs []Struct

// Structs represents a slice, array, map, or channel of Structs
type Structs struct {
	Value       reflect.Value
	Type        reflect.Type
	IndexType   reflect.Type
	ElementType reflect.Type
	Structs     []Struct
}

// func toStructs(value any) (Structs, error) {
// 	var s Structs
// 	return s, nil
// }

// func ToStructs[S ~[]T, T any](value T) (Structs, error) {
// 	var s Structs
// 	rv, rt, pointer := ToIndirectReflectValue(value)
// 	if !slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, rv.Kind()) {
// 		return s, fmt.Errorf("value %T is a %v, need a slice, array, map, or channel", value, rv.Kind())
// 	}
// 	if rv.Len() < 1 {
// 		value = *new(T)
// 	}
// 	var ss []Struct

// 	switch kind := rv.Kind(); {
// 	case kind == reflect.Slice:
// 	case kind == reflect.Map:
// 		iter := rv.MapRange()
// 		for iter.Next() {
// 			// k := iter.Key()oi
// 			v := iter.Value()
// 			ssi, err := ToStruct(v.Interface())
// 			if err != nil {
// 				return s, err
// 			}
// 			ss = append(ss, ssi)
// 		}
// 	case kind == reflect.Chan:
// 	case kind == reflect.Array:
// 	default:
// 	}

// 	// if len(value) < 1 {
// 	// 	// fmt.Printf("%T\n", *new(T))
// 	// 	s0, err := ToStruct(*new(T))
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	return append(s, s0), nil
// 	// }
// 	for _, v := range value {
// 		// fmt.Printf("%T\t%s\t%s\n", value, value, reflect.ValueOf(value).IsValid())
// 		// fmt.Printf("%T\t%s\t%s\n", v, v, reflect.ValueOf(v).IsValid())

// 		si, err := ToStruct(v)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s = append(s, si)
// 	}
// 	return s, nil
// }

// func NewStructs[S ~[]T, T any](value S, cfg StructConfig) (Structs, error) {
// 	s := []Struct{}
// 	if len(value) < 1 {
// 		s0, err := NewStruct(*new(T), cfg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return append(s, s0), nil
// 	}
// 	for _, v := range value {
// 		si, err := NewStruct(v, cfg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s = append(s, si)
// 	}
// 	return s, nil
// }

func (s Structs) Names() []string {
	var names []string
	for _, v := range s.Structs {
		names = append(names, v.Name)
	}
	return names
}

// set children:
// name
// *parent
// uuid
// Handle child slices of structs
// handle tags that indicate child/not child
//
