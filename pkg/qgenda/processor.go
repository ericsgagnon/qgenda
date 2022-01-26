package qgenda

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type Processor interface {
	Process() error
}

func IsProcessor(a any) bool {
	return ImplementsInterface[Processor](a)
}

func AsProcessor[T any](a T) (Processor, error) {
	if IsProcessor(a) {
		var iv interface{} = a
		return (iv).(Processor), nil
	}
	return nil, errors.New(fmt.Sprintf("%T does not implement Processor", a))
}

func CanSet(a any) bool {
	v := IndirectReflectionValue(a)
	k := v.Kind()
	if k == reflect.Invalid || !v.CanSet() {
		return false
	}
	return v.CanSet()
}

func Process(a any) error {
	switch {
	case IsProcessor(a):
		return a.(Processor).Process()
	case IsStruct(a):
		return ProcessStruct(a)
	case IsSlice(a):
		return ProcessSlice(a)
	case IsMap(a):
		return ProcessMap(a)
	default:
		// Process ignores any fields that dont' need processing
		return nil
	}
	// return errors.New(fmt.Sprintf("%T is not a Processor", a))
}

// ProcessRecursively dive's into any member or element processing.
// It then attempts to call a' Process method, if applicable.
func ProcessRecursively(a any) error {
	// v := IndirectReflectionValue(a)
	// fmt.Printf("%#v\n", v)
	switch {
	case IsStruct(a):
		if err := ProcessStruct(a); err != nil {
			return err
		}
	case IsSlice(a):
		if err := ProcessSlice(a); err != nil {
			return err
		}
	case IsMap(a):
		if err := ProcessMap(a); err != nil {
			return err
		}
	default:
		// Process ignores any fields that dont' need processing
		// return nil
	}
	if IsProcessor(a) {
		if err := a.(Processor).Process(); err != nil {
			return err
		}
	}
	return nil
	// return errors.New(fmt.Sprintf("%T is not a Processor", a))
}

// ProcessStruct doesn't attempt to check/use the struct's Process method.
//  Instead it iterates through each member and attempts to Process them.
// It also makes no effort to process members that are nil pointers or
// otherwise result in reflect.Kind() == reflect.Invalid.
func ProcessStruct(a any) error {
	v := IndirectReflectionValue(a)
	fields := StructFields(v)
	for i := 0; i < v.NumField(); i++ {
		f := IndirectReflectionValue(v.Field(i))
		sf := fields[i]
		if CanSet(f) && sf.IsExported() {
			fv := f.Interface()
			if IsProcessor(fv) {
				p := fv.(Processor)
				if err := p.Process(); err != nil {
					return err
				}
			}

			// err := Process(f)
			// if err != nil {
			// 	return err
			// }
		}
		// if sf.IsExported() {
		// 	fmt.Printf("%20s %T:  %s is settable: %t\n", fields[i].Name, f, f.Kind(), f.CanSet())
		// 	fmt.Println("--------------------------------------")

		// }
		// fmt.Println(f.CanAddr())

	}

	// for i := 0; i < v.Len(); i++ {
	// 	f := v.Index(i)
	// 	fv := f.Interface()
	// 	if ImplementsInterface[Processor](fv) {
	// 		// var p Processor
	// 		p, _ := fv.(Processor)
	// 		if err := p.Process(); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func ProcessStructFields(a any) {
	fmt.Println("--------------------------------------------------------------------------")
	v := reflect.ValueOf(a)
	fmt.Printf("%T\n", v)
	vi := reflect.Indirect(v)
	fmt.Printf("%T\n", vi)
	// reflect.ValueOf(&t).MethodByName("GFG").Call([]reflect.Value{})
	fields := StructFields(a)
	// for _, f := range fields {
	for i := 0; i < v.NumField(); i++ {
		// fmt.Println("--------------------------------------")
		// sf := v.Type().Field(i)
		f := v.Field(i) //.Addr()
		f = reflect.Indirect(f)
		if f.Kind() == reflect.Invalid || !f.CanSet() {
			fmt.Printf("%#v\n", fields[i])
		}
		if f.Kind() == reflect.Slice {
			// f = f.Elem()
			// fmt.Printf("%#v\n", f)
			ProcessSlice(f)
		}
		fmt.Printf("%T:  %s is settable: %t\n", a, f.Kind(), f.CanSet())
		fmt.Println("--------------------------------------")
		// fmt.Printf("%#v\n", f)
		// if f.Kind() == reflect.Pointer {
		// 	fv := f.Elem()
		// 	fv = reflect.Indirect(fv)
		// 	// fmt.Printf("%s: %#v\n", sf.Name, fv)

		// 	// if f.IsNil() {
		// 	// 	fmt.Printf("%T.%s is nil but settable: %t\n", a, f.Type(), fv.CanSet())
		// 	// 	// fmt.Println(f.Addr())
		// 	// 	fmt.Println("--------------------------------------")
		// 	// 	continue
		// 	// }
		// 	fmt.Printf("%T.%s is settable: %t\n", a, fv.Type(), fv.CanSet())
		// }
		// fmt.Printf("%T.%s is settable: %t\n", a, f.Type(), (reflect.Indirect(f)).CanSet())
		// fmt.Println("--------------------------------------")

		// fmt.Printf("%T is pointer: %t\n", a, f.Kind() == reflect.Pointer)

		// f.MethodByName("Process").Call([]reflect.Value{})
		// fv := v.FieldByName(f.Name)

		// fmt.Printf("%#v\n", fv)
		// res := fv.MethodByName("Process").Call([]reflect.Value{})
		// fmt.Println(res)
		// fmt.Printf("%#v\n", f)

	}
	fmt.Printf("%s\n", vi.Type())
}

func ProcessSlice(a any) error {
	v := reflect.Indirect(reflect.ValueOf(a))
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			fv := f.Interface()
			if ImplementsInterface[Processor](fv) {
				// var p Processor
				p, _ := fv.(Processor)
				if err := p.Process(); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New(fmt.Sprintf("%s is not a slice", v.Kind()))
	}
	// fmt.Println("I'm a slice")
	return nil
}

// func ProcessMap[K string, V any](m map[K]V) error {
// 	for _, v := range m {
// 		var i interface{} = v
// 		p := (i).(Processor)
// 		if err := p.Process(); err != nil {
// 			return err
// 		}

// 	}
// 	return nil
// }

// ProcessMap can currently only handle maps of pointers
// (or methods that can modify their receiver)
func ProcessMap(a any) error {
	v := IndirectReflectionValue(a)
	if v.Kind() != reflect.Map {
		return errors.New(fmt.Sprintf("%s is not a map", v.Kind()))
	}
	iter := v.MapRange()
	for iter.Next() {
		mv := iter.Value()
		mvi := mv.Interface()
		err := (mvi).(Processor).Process()
		if err != nil {
			return err
		}
	}
	return nil
}

// IndirectReflectionValue attempts to convert a to
// an indirect reflection value and return it
func IndirectReflectionValue(a any) reflect.Value {
	var v reflect.Value
	if reflect.ValueOf(a).Type().String() != "reflect.Value" {
		v = reflect.ValueOf(a)
	} else { // reflect.Value.Type == "reflect.Value"
		v = a.(reflect.Value)
	}
	if v.Kind() == reflect.Pointer {
		v = reflect.Indirect(v)
	}
	return v
}

// IndirectReflectionKind attempts to convert a to
// an indirect reflection kind and return it
func IndirectReflectionKind(a any) reflect.Kind {
	return IndirectReflectionValue(a).Kind()
}

// IsKind returns true if a's reflect.Kind == t
func IsKind(a any, t string) bool {
	// v := reflect.Indirect(reflect.ValueOf(a))
	v := IndirectReflectionValue(a)
	k := v.Type().Kind()
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

// StructFields de-references as
func StructFields(a any) []reflect.StructField {
	var structFields []reflect.StructField
	if IsStruct(a) {
		v := IndirectReflectionValue(a)
		// v := reflect.Indirect(reflect.ValueOf(a))
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			structFields = append(structFields, f)
		}
	}
	return structFields
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

// ImplementsInterface returns true if value implements Reference interface
// note: Reference must be passed as a type parameter
func ImplementsInterface[Reference any](value any) bool {
	_, ok := value.(Reference)
	return ok

}

func AnyProcess[T any](a T) (T, error) {
	switch {
	case IsProcessor(a):
		p, err := AsProcessor(a)
		if err != nil {
			return a, err
		}
		if err := p.Process(); err != nil {
			return a, err
		}
		out := (p).(T)
		return out, nil
	case IsStruct(a):
		out, err := StructProcess(a)
		if err != nil {
			return a, err
		}
		return out, nil
	case IsSlice(a):
		// v := reflect.ValueOf(a)
		// if v.CanConvert(reflect.Slice) {

		// }

		// out, err := SliceProcess(a)
		// if err != nil {
		// 	return a, err
		// }
		// return out, nil
	case IsMap(a):
		// out, err := MapProcess(a)
		// if err != nil {
		// 	return a, err
		// }
		// return out, nil
	default:
	}
	return a, nil
}

func MapProcess(a any) error {
	// v := reflect.ValueOf(a)
	return nil
}

func ReflectionStuff[T any, S []T](a T) (T, error) {
	fmt.Println("----------------------------------------------------")
	fmt.Printf("a's type is %T\n", a)
	v := reflect.ValueOf(a)
	fmt.Printf("Valueof: %s\n", v.String())
	fmt.Printf("Valueof.Kind: %s\n", v.Kind())
	t := v.Type()
	fmt.Printf("Valueof.Type: %s\n", t.String())
	fmt.Printf("Valueof.Type.Kind: %s\n", t.Kind())
	vi := reflect.Indirect(v)
	fmt.Printf("Indirect(ValueOf): %s\n", vi)
	fmt.Printf("Indirect(ValueOf).Type: %s\n", vi.Type())
	// goal: make a []T from a
	var iv interface{} = &a
	fmt.Printf("iv's type is %T\n", iv)
	// fmt.Println(iv)
	out, ok := (iv).(T)
	fmt.Println(ok)
	fmt.Printf("out: %T\n", out)
	fmt.Printf("%T\n", out)
	fmt.Println("------------")
	rs := reflect.SliceOf(t)
	fmt.Printf("reflect.SliceOf(reflect.ValueOf(a).Type()): %v\n", rs)
	fmt.Println("------------")
	TypeTest(a)
	fmt.Println("------------")
	// TypeTest[T](out)
	return out, nil
}

func TypeTest[T any](a T) (T, error) {
	// func TypeTest(a interface{}) (any, error) {
	// fmt.Println("------------------------------------------------------------")
	// fmt.Printf("Type of T: %T\n", new(T))
	// fmt.Printf("Type of S: %T\n", new(S))
	v := reflect.ValueOf(a)
	// fmt.Printf("Valueof: %s\n", v.String())
	// fmt.Printf("ValueOf.Elem: %s\n", v.Elem())
	// fmt.Printf("Valueof.Kind: %s\n", v.Kind())
	// fmt.Printf("Valueof.Index(0): %s\n", v.Index(0))
	// fmt.Printf("reflect.TypeOf(a).Elem(): %s\n", reflect.TypeOf(a).Elem())
	// sliceType := v.Index(0).Type()
	sliceType := reflect.TypeOf(a).Elem()
	// fmt.Printf("slicer type: %s\n", sliceType)
	// fmt.Printf("Valueof.Index(0).Type: %s\n", v.Index(0).Type())
	slicer := reflect.MakeSlice(reflect.SliceOf(sliceType), v.Len(), v.Len())
	fmt.Printf("slicer: %s\n", slicer.Type())
	t := v.Type()
	fmt.Printf("Valueof.Type: %s\n", t.String())
	fmt.Printf("Valueof.Type.Kind: %s\n", t.Kind())
	// rvi := reflect.Indirect(v)
	// fmt.Printf("Indirect(ValueOf): %s\n", rvi)
	// fmt.Printf("Indirect(ValueOf).Type: %s\n", rvi.Type())
	// fmt.Printf("Indirect(ValueOf).Kind: %s\n", rvi.Kind())
	// fmt.Printf("Indirect(ValueOf).Type.Kind: %s\n", rvi.Type().Kind())

	// var vi interface{} = a
	// fmt.Printf("vi: %#v\n", vi)
	// return (vi).(S), nil
	fmt.Println("------------")
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		// fmt.Println("I'm a little slicey-dicey")
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			fmt.Sprintf("v.Index(%d): %s\n", i, f)
			// fv := f.Interface()
		}
	default:
		// return errors.New(fmt.Sprintf("%s is not a slice", v.Kind()))
	}
	fmt.Println("------------------------------------------------------------")
	return slicer.Interface().(T), nil
	// return nil, nil
}

// v := reflect.Indirect(reflect.ValueOf(a))
// switch v.Kind() {
// case reflect.Array, reflect.Slice:
// 	for i := 0; i < v.Len(); i++ {
// 		f := v.Index(i)
// 		fv := f.Interface()
// 		if ImplementsInterface[Processor](fv) {
// 			// var p Processor
// 			p, _ := fv.(Processor)
// 			if err := p.Process(); err != nil {
// 				return err
// 			}
// 		}
// 	}
// default:
// 	return errors.New(fmt.Sprintf("%s is not a slice", v.Kind()))
// }
// // fmt.Println("I'm a slice")
// return nil

func SliceProcess[T any](a []T) ([]T, error) {
	// v := reflect.ValueOf(a[0])
	// fmt.Printf("Kind is %s\n", v.Kind())
	fmt.Println(IsStruct(a[0]))
	if !IsProcessor(a[0]) {
		return a, errors.New(fmt.Sprintf("%T does not implement processor", a[0]))
	}
	for i, v := range a {
		var iv interface{} = v
		p := (iv).(Processor)
		if err := p.Process(); err != nil {
			return nil, err
		}
		a[i] = p.(T)
	}
	return a, nil
}

func StructProcess[T any](a T) (T, error) {
	return a, nil
}
func GenericTests[T any](t T) *T {
	return &t
}

func ProcessTest[T any](a T) (T, error) {
	var i interface{} = a
	_, ok := i.(Processor)
	if ok {

	}
	return a, nil
}

// if reflect.ValueOf(a).Type().String() != "reflect.Value" {
// 	v = reflect.ValueOf(a)
// } else { // reflect.Value.Type == "reflect.Value"
// 	v = a.(reflect.Value)
// }
// if v.Kind() == reflect.Pointer {
// 	v = reflect.Indirect(v)
// }
