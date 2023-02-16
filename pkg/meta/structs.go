package meta

import (
	"fmt"
	"reflect"
)

type Structs []Struct

func toStructs(value any) (Structs, error) {

	return nil, nil
}

func ToStructs[S ~[]T, T any](value S) (Structs, error) {
	fmt.Printf("%s\n", reflect.ValueOf(value).Type())
	if len(value) < 1 {
		value = S{*new(T)}
	}
	s := []Struct{}
	// if len(value) < 1 {
	// 	// fmt.Printf("%T\n", *new(T))
	// 	s0, err := ToStruct(*new(T))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return append(s, s0), nil
	// }
	for _, v := range value {
		// fmt.Printf("%T\t%s\t%s\n", value, value, reflect.ValueOf(value).IsValid())
		// fmt.Printf("%T\t%s\t%s\n", v, v, reflect.ValueOf(v).IsValid())

		si, err := ToStruct(v)
		if err != nil {
			return nil, err
		}
		s = append(s, si)
	}
	return s, nil
}

func NewStructs[S ~[]T, T any](value S, cfg StructConfig) (Structs, error) {
	s := []Struct{}
	if len(value) < 1 {
		s0, err := NewStruct(*new(T), cfg)
		if err != nil {
			return nil, err
		}
		return append(s, s0), nil
	}
	for _, v := range value {
		si, err := NewStruct(v, cfg)
		if err != nil {
			return nil, err
		}
		s = append(s, si)
	}
	return s, nil
}
