package qgenda

import (
	"errors"
	"reflect"
)

// func ptr[T any](a T) *T {
// 	return &a
// }

// func val[T any](a *T) T {
// 	return *a
// }

// consolidate these somewhere
func typeName(a any) string {
	return reflect.TypeOf(a).Name()
}

var (
	ErrNope    = errors.New("Nope! Can't be done...")
	ErrMissing = errors.New("Missing or nil argument.")
)
