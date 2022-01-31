package qgenda

import (
	"fmt"
	"reflect"
)

// t.Elem panics if t is not Array, Chan, Map, Pointer, or Slice
func ReflectionInfo[T any](a T) {

	// let's make some sample values
	sampleValue := "string"
	samplePointer := Pointer("stringPointer")
	sampleSliceValue := []string{"one", "two", "three"}
	sampleSliceOfPointers := []*string{Pointer("one"), Pointer("two"), Pointer("three")}
	sampleSlicePointer := &[]string{"one", "two", "three"}
	sampleMapValue := map[string]string{
		"one":   "one",
		"two":   "two",
		"three": "three",
	}
	sampleMapOfPointers := map[string]*string{
		"one":   Pointer("one"),
		"two":   Pointer("two"),
		"three": Pointer("three"),
	}
	sampleMapPointer := &map[string]string{
		"one":   "one",
		"two":   "two",
		"three": "three",
	}
	sampleStruct := struct {
		Value           string
		Pointer         *string
		SliceValue      []string
		SliceOfPointers []*string
		SlicePointer    *[]string
		MapValue        map[string]string
		MapOfPointers   map[string]*string
		MapPointer      *map[string]string
	}{
		Value:           sampleValue,
		Pointer:         &(*samplePointer),
		SliceValue:      sampleSliceValue,
		SliceOfPointers: sampleSliceOfPointers,
		SlicePointer:    &(*sampleSlicePointer),
		MapValue:        *&sampleMapValue,
		MapOfPointers:   *&sampleMapOfPointers,
		MapPointer:      &*sampleMapPointer,
	}
	fmt.Sprint(sampleStruct)
	// value
	v := reflect.ValueOf(a)
	vt := v.Type() // == fmt.Printf("%T", a) == reflect.TypeOf(a)
	vk := v.Kind() // if v is zero Value, v.Kind() == reflect.Invalid
	vtk := vt.Kind()
	vString := fmt.Sprintf("v<%3s %-6.6s %-20.20s>", " ", vk.String(), vt)

	// indirect
	iv := reflect.Indirect(v)
	ivt := iv.Type()
	ivtk := ivt.Kind()
	if vk == reflect.Pointer {
		vString = fmt.Sprintf("v<%3s %-6.6s %-20.20s>", vk.String(), " ", vt)
	}

	ivString := fmt.Sprintf("iv<%-6.6s %-20.20s %-5t>", ivtk, ivt, iv.CanSet())

	// v.Elem panics if v is not interface or pointer
	ve := reflect.Value{}
	if vtk == reflect.Interface || vtk == reflect.Pointer {
		ve = v.Elem()
	}
	veString := fmt.Sprintf("%37s", " ")
	if ve.IsValid() {
		veString = fmt.Sprintf("el<%-6.6s %-20.20s %-5t>", ve.Kind(), ve.Type(), ve.CanSet())
	}

	fmt.Println(
		vString,
		ivString,
		veString,
		StructReflectionInfo(iv),
		SliceReflection(iv),
	)

}

func StructReflectionInfo(v reflect.Value) string {
	if v.Kind() != reflect.Struct {
		// fmt.Println("Not a struct")
		return ""
	}
	sfString := ""
	if v.NumField() > 0 {
		f := v.Field(0)
		// iv := reflect.Indirect(f)

		sfString = fmt.Sprintf("sf<%3s %-6.6s %-20.20s %-5t iv-set %-5t>", " ", f.Kind(), f.Type(), f.CanSet(), reflect.Indirect(f).CanSet())
		if f.Kind() == reflect.Pointer {
			sfString = fmt.Sprintf("sf<%3s %-6.6s %-20.20s %-5t iv-set %-5t>", f.Kind(), " ", f.Type(), f.CanSet(), reflect.Indirect(f).CanSet())

		}
	}
	return sfString
}

func SliceReflection(v reflect.Value) string {
	if v.Kind() != reflect.Slice {
		return ""
	}

	seString := ""
	if v.Len() > 0 {
		f := v.Index(0)
		iv := reflect.Indirect(f)
		seString = fmt.Sprintf("se<%3s %-6.6s %-20.20s %-5t iv-set %-5t>", " ", f.Kind(), f.Type(), f.CanSet(), iv.CanSet())
		if f.Kind() == reflect.Pointer {
			seString = fmt.Sprintf("se<%3s %-6.6s %-20.20s %-5t iv-set %-5t>", f.Kind(), " ", f.Type(), f.CanSet(), iv.CanSet())

		}
	}
	return seString

}
