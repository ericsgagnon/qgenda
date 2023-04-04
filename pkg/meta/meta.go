package meta

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

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

type Value struct {
	Name    any // optional - intended for struct field names, map keys, or slice/array indexes
	Type    reflect.Type
	Value   reflect.Value
	Pointer bool
	Parent  *Value // only populated if it is a member of a struct or elements of collections - slice, array, map, or channel
}

// ToValue 'crawls' reflection values/types, forming a linked list of Values with their parent & children (where applicable)
func ToValue(value any) (Value, error) {
	var v Value
	// children := map[string]Value{}
	rv, rt, pointer := ToIndirectReflectValue(value)
	if rt == nil {
		return v, fmt.Errorf("invalid value: %v", value)
	}

	v = Value{
		Type:    rt,
		Value:   rv,
		Pointer: pointer,
	}

	return v, nil
}

func (v Value) Kind() reflect.Kind {
	return v.Value.Kind()
}

func (v Value) Valid() bool {
	return v.Value.IsValid()
}

func (v *Value) Children() []Value {
	var children []Value
	switch kind := v.Kind(); {
	case kind == reflect.Invalid:
		return children
	case kind == reflect.Chan:
		// treat channel's a 0 length
		fallthrough
	case slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, kind) &&
		v.Value.Len() == 0:

		return children
	case kind == reflect.Map:
		iter := v.Value.MapRange()
		for iter.Next() {
			child, err := ToValue(iter.Value().Interface())
			if child.Value.Kind() != reflect.Invalid && err != nil {
				return children
			}
			child.Name = iter.Key().String()
			// children[iter.Key().String()] = child
			children = append(children, child)
		}
	case kind == reflect.Slice || kind == reflect.Array:
		// precision := len(fmt.Sprint(rv.Len())) + 2
		for i := 0; i < v.Value.Len(); i++ {
			child, err := ToValue(v.Value.Index(i).Interface())
			if child.Value.Kind() != reflect.Invalid && err != nil {
				return children
			}
			// key := fmt.Sprintf("%.[1]*d\n", precision, i)
			// children[key] = child
			// child.Name = key
			children = append(children, child)
		}
	case kind == reflect.Struct:
		for _, field := range reflect.VisibleFields(v.Type) {
			if field.Anonymous || !field.IsExported() {
				continue
			}
			fieldValue := v.Value.FieldByName(field.Name)
			if fieldValue.Kind() == reflect.Invalid {
				fieldValue = reflect.New(field.Type).Elem()
			}
			fmt.Println(field.Name, ":\t", fieldValue.Type())
			child, err := ToValue(fieldValue.Interface())
			if err != nil {
				return children
			}
			// children[field.Name] = child
			child.Name = field.Name
			children = append(children, child)
		}
	default:
	}

	for _, child := range children {
		child.Parent = v
	}
	return children

}

func (v Value) Child(a any) (Value, error) {
	children := v.Children()
	if len(children) == 0 {
		return Value{}, fmt.Errorf("%v has no children", v)
	}
	switch t := a.(type) {
	case int:
		if slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Struct}, v.Kind()) && len(children) >= t {
			return children[t], nil
		}
		// return Value{}, nil
	case string:
		if slices.Contains([]reflect.Kind{reflect.Map, reflect.Struct}, v.Kind()) {
			for _, child := range children {
				if child.Name == t {
					return child, nil
				}
			}
		}
	default:
		return Value{}, fmt.Errorf("method child(a any) expects an int or string, got %v", t)
	}
	return Value{}, fmt.Errorf("no matching child for %v", a)
}

// NewElement returns a blank element for slices, arrays, maps and channels.
// For everything else it returns an (invalid) Value and a non-nil error.
func (v Value) NewElement() (Value, error) {
	var element Value
	if !v.Valid() {
		return element, fmt.Errorf("cannot determine child type of invalid value")
	}
	if slices.Contains([]reflect.Kind{reflect.Slice, reflect.Array, reflect.Map, reflect.Chan}, v.Kind()) {
		elemRT := v.Value.Type().Elem()
		return ToValue(elemRT)
	}
	return element, fmt.Errorf("value %v is a %s, not slice, array, map, or channel", v, v.Kind())
}
