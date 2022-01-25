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

func AsProcessor(a any) Processor {
	if IsProcessor(a) {
		return a.(Processor)
	}
	return nil
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
		// if f.Kind() == reflect.Invalid {
		// 	f = v.Field(i)
		// }
		// f := (v.Field(i))
		// f = reflect.Indirect(f)
		// if f.Kind() == reflect.Map {
		// 	err := ProcessMap(f)
		// 	if err != nil {
		// 		return err
		// 	}
		// } else
		if CanSet(f) && sf.IsExported() {
			err := Process(f)
			if err != nil {
				return err
			}
		}
		// else if f.Kind() == reflect.Invalid {
		// 	nv := reflect.New(v.Field(i).Type())
		// 	fmt.Printf("New Value: %T\n", nv.Type().Name())
		// }
		if sf.IsExported() {
			fmt.Printf("%20s %T:  %s is settable: %t\n", fields[i].Name, f, f.Kind(), f.CanSet())
			fmt.Println("--------------------------------------")

		}
		// fmt.Println(f.CanAddr())

	}
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
			// fmt.Printf("Field %d - %T: canset: %t\n", i, f.Kind(), f.CanSet())
			// implementsProcessor := f.Type().Implements(reflect.TypeOf(new(Processor)).Elem())
			// if implementsProcessor {
			// 	result := f.MethodByName("Process").Call([]reflect.Value{})
			// 	err := result[0].Interface().(error)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	// fmt.Println(check[0].IsNil())

			// }
			// fmt.Printf("%#v\n", implementsProcessor)
			// if err := Process(fv); err != nil {
			// 	return err
			// }
			// fmt.Printf("Slice field %d: %#v\n", i, v.Index(i))

			// 	func main() {
			// 		fmt.Println(check(new(Hello), new(Person))) // false
			//    }

			//    func check(i interface{}, n interface{}) bool {
			// 	   ti := reflect.TypeOf(i).Elem()
			// 	   return reflect.TypeOf(n).Implements(ti)
			//    }
		}
	default:
		return errors.New(fmt.Sprintf("%T is not a slice", v.Kind()))
	}
	// fmt.Println("I'm a slice")
	return nil
}

func MapTest(a any) error {
	v := reflect.ValueOf(a)
	// ve := reflect.Indirect(v)
	nm := reflect.MakeMap(v.Type())
	iter := v.MapRange()
	for iter.Next() {
		mk := iter.Key()
		mv := iter.Value()
		// mvi := reflect.Indirect(mv)
		mvt := mv.Type()
		// mva := mv.Interface()
		nv := reflect.New(mvt)
		// var p Processor
		p := (nv.Interface()).(Processor)
		// fmt.Println(nv)
		if err := p.Process(); err != nil {
			return err
		}
		// fmt.Println(p)
		nm.SetMapIndex(mk, reflect.ValueOf(p).Elem())
		// nm.SetMapIndex(mk, reflect.ValueOf(p).Elem())
		// fmt.Println(mk, mv, mvi, mvt, mva, nv.Type(), nv.Elem().CanSet(), nv.Elem().Type(), nv.Interface())
	}
	out := reflect.ValueOf(&a).Elem()
	out.Set(nm)
	// fmt.Println(out)
	// ve.Set(nm)
	return nil
}

func ProcessMap(a any) error {
	v := IndirectReflectionValue(a)
	fmt.Printf("%s\n", reflect.ValueOf(a).Pointer())
	if v.Kind() != reflect.Map {
		return errors.New("Value is not a map")
	}
	iter := v.MapRange()
	for iter.Next() {
		// mk := iter.Key()
		mv := iter.Value()
		mvi := mv.Interface()
		err := (mvi).(Processor).Process()
		// err := Process(mv)
		if err != nil {
			return err
		}
		// fmt.Printf("%s: %s - %T\n", mk, mv, mv)
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
func ImplementsInterface[Reference any](value interface{}) bool {
	_, ok := value.(Reference)
	return ok

}
