package meta

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// Meta is intended to capture the de-referenced, underlying value of a variable or type.
// If it is constructed from a map, slice, array, channel, it will have the container type.
// It also indicates if it was originally a pointer
type Meta struct {
	Type      reflect.Type
	Value     reflect.Value
	Container reflect.Type
	Pointer   bool
}

func ToMeta(value any) Meta {
	var rv reflect.Value
	var rt reflect.Type

	switch v := value.(type) {
	case nil:
		// do nothing
	case reflect.Type:
		rv = reflect.New(v).Elem()
	case reflect.Value:
		rv = v
	default:
		rv = reflect.ValueOf(v)
	}

	var pointer bool
	switch {
	case rv.Kind() == reflect.Invalid:
		// do nothing
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = reflect.New(rt).Elem()
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = rv.Elem()
	default:
		rt = rv.Type()
	}
	containerKinds := []reflect.Kind{
		reflect.Array, reflect.Chan, reflect.Map, reflect.Slice,
	}
	if slices.Contains(containerKinds, rt.Kind()) {
		fmt.Println("I'm a container")
	}
	m := Meta{
		Type:    rt,
		Value:   rv,
		Pointer: pointer,
	}
	return m

}

func IndirectReflectValue(value any) reflect.Value {
	if v, ok := value.(reflect.Value); ok {
		if v.Kind() == reflect.Pointer {
			return reflect.Indirect(v)
		}
		return v
	}
	if t, ok := value.(reflect.Type); ok {
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		return reflect.New(t)
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		return reflect.Indirect(v)
	}
	return v
}

// ToReflectValue attempts to convert any value to a reflect.Value. If value is nil,
// an invalid reflect.Value, or a nil pointer, the returned value will be invalid.
func ToReflectValue(value any) reflect.Value {
	switch v := value.(type) {
	case nil:
		return reflect.Value{}
	case reflect.Type:
		return reflect.New(v).Elem()
	case reflect.Value:
		return v
	default:
		return reflect.ValueOf(v)
	}
}

func DereferenceReflectValue(value reflect.Value) (reflect.Value, reflect.Type, bool) {
	var rv reflect.Value
	var rt reflect.Type
	var pointer bool

	switch {
	case rv.Kind() == reflect.Invalid:
		// do nothing
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = reflect.New(rt).Elem()
	// case rv.Kind() == reflect.Interface && rv.Elem().Kind() == reflect.Invalid:
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = rv.Elem()
	default:
		rt = rv.Type()
	}

	return rv, rt, pointer
}

func DereferenceYeah(value reflect.Value) (reflect.Value, reflect.Type, bool, bool) {
	var rv reflect.Value
	var rt reflect.Type
	var pointer bool
	var boxed bool

	switch {
	case rv.Kind() == reflect.Invalid:
		// do nothing
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = reflect.New(rt).Elem()
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = rv.Elem()
	case rv.Kind() == reflect.Interface && rv.Elem().Kind() == reflect.Invalid:
		boxed = true
		rt = rv.Type().Elem()
		rv = reflect.New(rt).Elem()
	case rv.Kind() == reflect.Interface && rv.Elem().Kind() != reflect.Invalid:

	default:
		rt = rv.Type()
	}

	return rv, rt, pointer, boxed
}

// ToIndirectReflectValue attempts to convert any value to a reflect.Value and indirect it. If value is nil,
// an invalid reflect.Value, or a nil pointer, the returned value will be invalid. It returns the indrected value
// and whether the original value was a pointer
func ToIndirectReflectValue(value any) (reflect.Value, reflect.Type, bool) {
	var rv reflect.Value
	var rt reflect.Type

	switch v := value.(type) {
	case nil:
		// do nothing
	case reflect.Type:
		rv = reflect.New(v).Elem()
	case reflect.Value:
		rv = v
	default:
		rv = reflect.ValueOf(v)
	}

	var pointer bool
	switch {
	case rv.Kind() == reflect.Invalid:
		// do nothing
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = reflect.New(rt).Elem()
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
		pointer = true
		rt = rv.Type().Elem()
		rv = rv.Elem()
	default:
		rt = rv.Type()
	}

	return rv, rt, pointer
}

// ToIndirectReflectValue attempts to convert any value to a reflect.Value and indirect it. If value is nil,
// an invalid reflect.Value, or a nil pointer, the returned value will be invalid. It returns the indrected value
// and whether the original value was a pointer
func ToIndirect(value any) (reflect.Value, reflect.Type, bool, bool) {
	var rv reflect.Value
	var rt reflect.Type

	switch v := value.(type) {
	case nil:
		// do nothing
	case reflect.Type:
		rv = reflect.New(v).Elem()
	case reflect.Value:
		rv = v
	default:
		rv = reflect.ValueOf(v)
	}

	var pointer bool
	// switch {
	// case rv.Kind() == reflect.Invalid:
	// 	// do nothing
	// case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
	// 	pointer = true
	// 	rt = rv.Type().Elem()
	// 	rv = reflect.New(rt).Elem()
	// case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
	// 	pointer = true
	// 	rt = rv.Type().Elem()
	// 	rv = rv.Elem()
	// default:
	// 	rt = rv.Type()
	// }

	var boxed bool
	switch {
	case rv.Kind() == reflect.Invalid:
		// do nothing
	case rv.Kind() == reflect.Interface && rv.Elem().Kind() == reflect.Invalid:
		boxed = true
	case rv.Kind() == reflect.Interface && rv.Elem().Kind() != reflect.Invalid:
		boxed = true
		rv, rt, pointer, _ = ToIndirect(rv.Elem().Interface())
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Invalid:
		pointer = true
	case rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Invalid:
		pointer = true
		rv, rt, _, boxed = ToIndirect(rv.Elem().Interface())
	default:
		rt = rv.Type()
	}

	return rv, rt, pointer, boxed
}

type Children[Key comparable] interface {
	map[Key]Value
}

// Collection is the Value type for slices, arrays, maps, and channels
type Collection struct {
	Type     reflect.Type
	Value    reflect.Value
	Parent   *Value  // only populated if it is a member of a struct or elements of collections - slice, array, map, or channel
	Children []Value // members of structs or elements of collections - slice, array, map, channel
	Pointer  bool
}

// Value is a (recursively) dereferenced merging of reflect Value and Type.
type Value struct {
	Type     reflect.Type
	Value    reflect.Value
	Parent   *Value  // only populated if it is a member of a struct or elements of collections - slice, array, map, or channel
	Children []Value // members of structs or elements of collections - slice, array, map, channel
	Pointer  bool
	// IndexType reflect.Type // index type for slices, arrays, maps, and channels
}

// ToValue 'crawls' reflection values/types, forming a linked list of Values with their parent & children (where applicable)
func ToValue(value any) (Value, error) {
	var v Value
	var children []Value
	rv, rt, pointer := ToIndirectReflectValue(value)
	// fmt.Printf("value: %v\ttype: %v\tpointer: %t\n", rv, rt, pointer)
	if rt == nil {
		return v, fmt.Errorf("invalid value: %v", value)
	}

	switch kind := rt.Kind(); {
	case kind == reflect.Invalid:
		return v, fmt.Errorf("invalid value: Kind() == reflect.Invalid: %v", value)
	case kind == reflect.Chan:
		// treat channel's a 0 length
		fallthrough
	case slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, kind) &&
		rv.Len() == 0:
		// fmt.Println("WHAT???????????????????????????????????? ", rv.Len())
		// do nothing - length 0 doesn't have children
		// childType := rt.Elem()
		// fmt.Println("child type is...\t", childType)
		// childValue := reflect.New(childType).Elem().Interface()
		// child, err := ToValue(childValue)
		// if err != nil {
		// 	return v, err
		// }
		// children = append(children, child)
	case kind == reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			child, err := ToValue(iter.Value().Interface())
			if child.Value.Kind() != reflect.Invalid && err != nil {
				return v, err
			}
			children = append(children, child)
		}
	case kind == reflect.Slice || kind == reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			child, err := ToValue(rv.Index(i).Interface())
			if child.Value.Kind() != reflect.Invalid && err != nil {
				return v, err
			}
			children = append(children, child)
		}
	case kind == reflect.Struct:
		// Children are StructFields -> Fields
		// var fieldNames []string
		// for _, field := range reflect.VisibleFields(rt) {
		// 	if field.IsExported() && !field.Anonymous {
		// 		fieldNames = append(fieldNames, field.Name)
		// 	}
		// }
		for _, field := range reflect.VisibleFields(rt) {
			if field.Anonymous || !field.IsExported() {
				continue
			}
			fieldValue := rv.FieldByName(field.Name)
			if fieldValue.Kind() == reflect.Invalid {
				fieldValue = reflect.New(field.Type).Elem()
			}
			fmt.Println(field.Name, ":\t", fieldValue.Type())
			child, err := ToValue(fieldValue.Interface())
			if err != nil {
				return v, err
			}
			children = append(children, child)
		}
	default:
	}
	v = Value{
		Type:     rt,
		Value:    rv,
		Pointer:  pointer,
		Children: children,
	}

	for _, child := range children {
		child.Parent = &v
	}

	return v, nil
}

func (v Value) ChildTypes(index ...int) []reflect.Type {
	var types []reflect.Type
	switch length := len(v.Children); {
	case length == 0:
		for range index {
			types = append(types, v.Type.Elem())
		}
	default:
		for i := range index {
			child := v.Children[i]
			types = append(types, child.Type)
		}
	}
	return types
}

func (v Value) Kind() reflect.Kind {
	return v.Type.Kind()
}

func Unbox(value any) {
	rv := reflect.ValueOf(value)
	// rt := reflect.TypeOf(value)
	if rv.Kind() == reflect.Interface {
		fmt.Println("value was boxed")
		rv = rv.Elem()
	}

	fmt.Println(rv.Kind())
}

// func MapToValue(value any) (Value, error) {

// 	rt = rt.Elem()
// 	iter := rv.MapRange()
// 	for iter.Next() {
// 		k := iter.Key()
// 		v := iter.Value()
// 		fmt.Printf("%v\t%v\t%v\t%T\t%v\n", k, v.Kind(), v.Interface(), v.Interface(), reflect.TypeOf(v.Interface()))
// 		fmt.Println(v.Elem().Kind(), v.Elem().Type())
// 		break
// 	}
// 	if rv.Len() > 0 {
// 		mv0 := rv.MapIndex(rv.MapKeys()[0])
// 		if mv0.Kind() == reflect.Interface {
// 			mv0 = mv0.Elem()
// 		}
// 		fmt.Println(mv0.Type())
// 	}
// 	return Value{}, nil
// }

// switch {
// case rv.Kind() == reflect.Map:
// case rv.Kind() == reflect.Slice:
// 	if rv.Len() > 0 {
// 		rv0 := rv.Index(0)
// 		if rv0.Kind() == reflect.Interface {
// 			rv0 = rv0.Elem()
// 		}
// 		fmt.Printf("%v\t%v\n", rv0.Kind(), rv0.Type())
// 	}

// 	fmt.Println("slice")
// case rv.Kind() == reflect.Array:
// 	fmt.Println("array")
// case rv.Kind() == reflect.Struct:
// 	fmt.Println("struct")
// case rv.Kind() == reflect.Pointer:
// 	fmt.Println("pointer")
// case rv.Kind() == reflect.Interface:
// 	fmt.Println("interface")
// }
