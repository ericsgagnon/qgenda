package qgenda

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func Hash[V []byte | string](b V) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(b)))
}

// func Value[T any](a *T) T {
// 	return *a
// }

// // Pointer simply returns a pointer to a value. It is useful
// // when using literals for pointer assignments.
// func Pointer[T any](t T) *T {
// 	return &t
// }

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

// IndirectReflectionValue attempts to convert a to
// an indirect reflection value, dereference it, and return it
func IndirectReflectionValue(a any) reflect.Value {
	rv, ok := a.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(a)
	}
	if rv.Kind() == reflect.Pointer {
		return reflect.Indirect(rv)
	}
	return rv
}

// // IndirectReflectionKind attempts to convert a to
// // an indirect reflection kind and return it
// func IndirectReflectionKind(a any) reflect.Kind {
// 	return IndirectReflectionValue(a).Kind()
// }

// IsKind returns true if a's reflect.Kind == t
func IsKind(a any, t string) bool {
	// v := reflect.Indirect(reflect.ValueOf(a))
	v := IndirectReflectionValue(a)
	// k := v.Type().Kind()
	k := v.Kind()

	return (k.String() == t)
}

// IsMap returns true if a's kind is a map (or the ill-advised pointer to a map)
func IsMap(a any) bool {
	return IsKind(a, "map")
}

// IsSlice returns true if a's kind is a slice/array or pointer to a slice/array
func IsSlice(a any) bool {
	isSlice := IsKind(a, "slice")
	isArray := IsKind(a, "array")
	return isSlice || isArray
	// v := reflect.Indirect(reflect.ValueOf(a))
	// k := v.Type().Kind()
	// return (k.String() == "struct")
}

// IsStruct returns true if a's kind is a struct or a pointer to a struct
func IsStruct(a any) bool {
	return IsKind(a, "struct")
}

// ImplementsInterface returns true if value implements Reference interface
// note: Reference must be passed as a type parameter
func ImplementsInterface[Reference any](value any) bool {
	_, ok := value.(Reference)
	return ok

}

func CanSet(a any) bool {
	v := IndirectReflectionValue(a)
	k := v.Kind()
	if k == reflect.Invalid || !v.CanSet() {
		return false
	}
	return v.CanSet()
}

// func StructFieldValues[St any](st St) []reflect.Value {
// 	var v reflect.Value
// 	if IsStruct(st) {

// 	}
// }

// StructFields de-references as
func StructFields(a any) []reflect.StructField {
	var sf []reflect.StructField
	switch rt := reflect.TypeOf(a); {
	case rt.Kind() != reflect.Struct:
		return nil
	case rt.NumField() == 0:
		return nil
	default:
		for i := 0; i < rt.NumField(); i++ {
			sf = append(sf, rt.Field(i))
		}
	}
	return sf
}

func StructFieldNames(a any) []string {
	var fieldNames []string
	if IsStruct(a) {
		// v := reflect.Indirect(reflect.ValueOf(a))
		v := IndirectReflectionValue(a)
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i).Name
			fieldNames = append(fieldNames, f)

			// fv := reflect.ValueOf(f)
			// if _, ok := afMap[f.Name]; ok {
			// 	fmt.Printf("%2d:\t%s\t%t\n", i, f.Name, v.Field(i).IsNil())
			// }
		}

	}
	return fieldNames
}

func StructFieldByName(a any, s string) reflect.StructField {
	if IsStruct(a) {
		v := reflect.Indirect(reflect.ValueOf(a))
		f, ok := v.Type().FieldByName(s)
		if !ok {
			log.Printf("Type %s doesn't have a field %#v\n", v.Type().String(), s)
		}
		return f

	}
	return reflect.StructField{}
}

func DynamicType(v reflect.Value) reflect.Type {
	if v.Kind() == reflect.Pointer {
		return reflect.TypeOf(v.Interface()).Elem()
	}
	return v.Type()
}

// ExpandEnvVars substitutes environment variables of the form ${ENV_VAR_NAME}
// if you have characters that need to be escaped, they should be surrounded in
// quotes in the source string.
func ExpandEnvVars(s string) string {
	re := regexp.MustCompile(`\$\{[^}]+\}`)

	envvars := map[string]string{}
	for _, m := range re.FindAllString(s, -1) {
		mre := regexp.MustCompile(`[${}]`)
		mtrimmed := mre.ReplaceAllString(m, "")
		// fmt.Printf("%s:\t%s\n", mtrimmed, os.Getenv(mtrimmed))
		envvars[m] = os.Getenv(mtrimmed)
	}

	for k, v := range envvars {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

// ExpandFileContents substitutes the placeholder with the contents on the first
// line of a file. It only accepts the pattern {file:/path/to/file}
// if you have characters that need to be escaped, they should be surrounded in
// quotes in the source string.
func ExpandFileContents(s string) string {
	re := regexp.MustCompile(`\{file:[^}]+\}`)

	files := map[string]string{}
	for _, filename := range re.FindAllString(s, -1) {
		idpattern := regexp.MustCompile(`(^\{file:)|(\}$)`)
		fn := idpattern.ReplaceAllString(filename, "")
		// fmt.Printf("%s:\t%s\n", mtrimmed, os.Getenv(mtrimmed))
		b, err := os.ReadFile(fn)
		if err != nil {
			panic(err)
		}
		fc := string(b)
		fca := strings.Split(fc, "\n")
		files[filename] = fca[0]
	}

	for k, v := range files {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

// since we're trying pointer-members, it's better to pass
// literal's as needed rather than create temporary variables
// that we might inadvertently tamper with
func stringPointer(s string) *string     { return &s }
func intPointer(i int) *int              { return &i }
func boolPointer(b bool) *bool           { return &b }
func timePointer(t time.Time) *time.Time { return &t }
func floatPointer(f float64) *float64    { return &f }

func stringFromPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func intFromPointer(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
func boolFromPointer(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
func timeFromPointer(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
func floatFromPointer(f *float64) float64 {
	if f != nil {
		return *f
	}
	return float64(0)
}
