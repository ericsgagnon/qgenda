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

func Process(a any) error {
	var p Processor
	if ImplementsInterface[Processor](a) {
		// fmt.Printf("%T is a Processor!!\n", a)
		p = a.(Processor)

		fmt.Printf("%T is settable: %t\n", a, reflect.ValueOf(a).Elem().CanSet())
		return p.Process()
	} else if IsStruct(a) {
		// fmt.Printf("%T is a struct\n", a)
		ProcessStructFields(a)
	}
	return errors.New(fmt.Sprintf("%T is not a Processor", a))
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
		// v = reflect.Indirect(v)
		// v = v.Slice(0, 0)
		// fmt.Printf("%T: let's get this done - canset: %t\n", v.Kind(), v.CanSet())
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

func ProcessMap(m reflect.Value) error {

	if m.Kind() != reflect.Map {
		return errors.New("Value is not a map")
	}
	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		fmt.Printf("%s: %s\n", k, v)
	}
	return nil
}

// func something(f interface{}) error {
//     v := reflect.ValueOf(f)
//     if v.Kind() != reflect.Ptr  {
//         return fmt.Errorf("not ptr; is %T", f)
//     }
//     v := v.Elem() // dereference the pointer
//     if v.Kind() != reflect.Struct  {
//         return fmt.Errorf("not struct; is %T", f)
//     }
//     t := v.Type()
//     for i := 0; i < t.NumField(); i++ {
//         sf := t.Field(i)
//         fmt.Println(sf.Name, v.Field(i).Interface())
//     }
//     return nil
// }

// IsKind returns true if a's reflect.Kind == t
func IsKind(a any, t string) bool {
	v := reflect.Indirect(reflect.ValueOf(a))
	k := v.Type().Kind()
	return (k.String() == t)
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
	// v := reflect.Indirect(reflect.ValueOf(a))
	// k := v.Type().Kind()
	// return (k.String() == "struct")
}

// StructFields de-references as
func StructFields(a any) []reflect.StructField {
	var structFields []reflect.StructField
	if IsStruct(a) {

		v := reflect.Indirect(reflect.ValueOf(a))
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
		v := reflect.Indirect(reflect.ValueOf(a))
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
