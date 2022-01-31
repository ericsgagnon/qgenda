package qgenda

type ValueProcessor[T any] interface {
	ProcessValue() (T, error)
}

// func IsValueProcessor[T any](a T) bool {
// 	var iv interface{} = a
// 	_, ok := iv.(ValueProcessor[T])
// 	return ok
// }

// func AsValueProcessor[T any](a T) (ValueProcessor[T], error) {
// 	if IsProcessor(a) {
// 		var iv interface{} = a
// 		return (iv).(ValueProcessor[T]), nil
// 	}
// 	return nil, errors.New(fmt.Sprintf("%T does not implement ValueProcessor", a))
// }

// func ProcessValue[T any](a T) (T, error) {
// 	switch {
// 	case IsValueProcessor(a):
// 		p, err := AsValueProcessor(a)
// 		if err != nil {
// 			return *new(T), err
// 		}
// 		return p.ProcessValue()
// 	case IsStruct(a):
// 		return ProcessStructValue(a)
// 	// case IsSlice(a):
// 	// 	return ProcessSliceValue(a)
// 	case IsMap(a):
// 		out := ToMap(a)
// 		fmt.Println("out:    ", out)
// 		out, err := ProcessMapValue(out)
// 		fmt.Println("out:    ", out)
// 		if err != nil {
// 			return *new(T), err
// 		}
// 		outValue := MapToAny(out, a)
// 		fmt.Println("outValue:    ", outValue)
// 		return outValue, nil
// 	default:
// 		// Process ignores any fields that dont' need processing
// 		return *new(T), nil
// 	}
// 	// return errors.New(fmt.Sprintf("%T is not a Processor", a))
// }

// func ProcessStructValue[T any](a T) (T, error) {
// 	// sf := StructFields(a)

// 	return a, nil
// }

// // ProcessSliceValue is used as part of the ValueProcessor interface
// // for top level variables, you should generally use ProcessValue
// func ProcessSliceValue[T any](a []T) ([]T, error) {
// 	out := []T{}
// 	for _, v := range a {
// 		nv, err := ProcessValue(v)
// 		if err != nil {
// 			return nil, err
// 		}
// 		out = append(out, nv)
// 	}
// 	return out, nil
// }

// // ProcessMapValue is used as part of the ValueProcessor interface
// // for top level variables, you should generally use ProcessValue
// func ProcessMapValue[K comparable, V any, M map[K]V](a M) (M, error) {
// 	out := M{}
// 	for k, v := range a {
// 		nv, err := ProcessValue(v)
// 		if err != nil {
// 			return nil, err
// 		}
// 		out[k] = nv
// 	}
// 	return out, nil
// }
