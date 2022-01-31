package qgenda

import (
	"fmt"
	"log"
	"reflect"
)

func ToSlice[T any](a T) []any {
	v := reflect.ValueOf(a)
	var s []any
	if v.Kind() != reflect.Slice {
		fmt.Printf("%T is not a slice\n", a)
	}
	if v.Kind() == reflect.Slice {
		fmt.Printf("%T is a slice\n", a)
		iv := reflect.Indirect(v)
		sliceType := reflect.TypeOf(a).Elem()
		out := reflect.MakeSlice(reflect.SliceOf(sliceType), iv.Len(), iv.Len())
		fmt.Printf("%T type of out var", out)
		fmt.Println(out)
		for i := 0; i < iv.Len(); i++ {
			f := reflect.Indirect(iv.Index(i))
			out.Index(i).Set(f)
			s = append(s, f.Interface())
		}
		fmt.Println(out)
		fmt.Printf("%T type of out var", out.Interface())
	}
	// fmt.Printf("\n\n%T\n%v\n", s, s)
	return s
}

func ToMap[T any](a T) map[any]any {
	v := reflect.Indirect(reflect.ValueOf(a))
	if v.Kind() != reflect.Map {
		return nil
	}
	out := map[any]any{}
	iter := v.MapRange()
	for iter.Next() {
		k := iter.Key().Interface()
		v := iter.Value().Interface()
		out[k] = v
	}
	return out
}

func ToMapStringAny[T any](a T) map[string]any {
	v := reflect.Indirect(reflect.ValueOf(a))
	if v.Kind() != reflect.Map {
		return nil
	}
	out := map[string]any{}
	iter := v.MapRange()
	for iter.Next() {
		// iv := reflect.Indirect(iter.Value())
		k := iter.Key()
		v := iter.Value()

		out[k.Interface().(string)] = v.Interface()
	}
	return out
}

func MapToAny[M map[K]V, T any, K comparable, V any](m M, a T) T {
	v := reflect.Indirect(reflect.ValueOf(a))
	mv := reflect.Indirect(reflect.ValueOf(m))
	if v.Kind() != reflect.Map {
		return *new(T)
	}
	keyType := reflect.TypeOf(a).Key()
	elementType := reflect.TypeOf(a).Elem()
	out := reflect.Indirect(reflect.MakeMapWithSize(reflect.MapOf(keyType, elementType), v.Len()))

	iter := mv.MapRange()
	for iter.Next() {
		fmt.Printf("key: %#v\n", iter.Key())
		fmt.Printf("value: %#v\n", iter.Value())
		out.SetMapIndex(iter.Key(), iter.Value())

	}
	fmt.Printf("MapToAny out: %#v\n", out)
	outValue, ok := out.Interface().(T)
	fmt.Printf("MapToAny outValue: %#v\n", outValue)
	if !ok {
		log.Printf("Could not convert %T to %T\n", out, a)
	}
	return outValue
}
