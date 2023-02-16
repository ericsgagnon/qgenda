package meta

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// Struct captures the properties of a struct and allows overriding properties using tags or by direct assignment
type Struct struct {
	Name               string
	NameSpace          []string
	NameSpaceSeparator string        // defaults to "."
	Type               reflect.Type  // pointers will be de-referenced (indirected)
	Value              reflect.Value // pointers will be de-referenced (indirected)
	UUID               string
	Attributes         map[string]string
	Fields             Fields
	Children           []Structs // fields that are slices of structs with more than one member
	Tags               map[string][]string
	Parent             *Struct
	pointer            bool // is the original variable a pointer
}

// SetUUID recursively sets s and its fields' struct UUID's to id
func (s *Struct) SetUUID(id string) {
	s.UUID = id
	if s.Fields != nil {
		s.Fields.SetUUID(id)
	}
}

func ToStruct(value any) (Struct, error) {
	var s Struct
	rv, rt, pointer := ToIndirectReflectValue(value)

	if rt.Kind() != reflect.Struct {
		return s, fmt.Errorf("invalid type: (%s) %s", rt.Kind(), rt)
	}

	s = Struct{
		Name:       rt.Name(),
		Type:       rt,
		Value:      rv,
		UUID:       strings.ReplaceAll(uuid.NewString(), "-", ""),
		Attributes: nil,
		Fields:     ToFields(value),
		pointer:    pointer,
	}
	s.SetUUID(s.UUID)

	return s, nil
}

type StructConfig struct {
	Name                     string
	NameSpace                []string
	NameSpaceSeparator       string // default is "."
	UUID                     string
	Tags                     map[string][]string // struct tags, not field tags, may optionally be parsed from tags labeled struct
	RemoveExistingTags       bool                // remove existing tags - false will simply apopend any included tags to current ones
	Attributes               map[string]string   // these should be table constraints, not field constraints
	RemoveExistingAttributes bool                // remove existing constraints - false will simply apopend any included constraints to current ones
}

// NewStruct enables configuration when parsing a struct to a Struct
func NewStruct(a any, cfg StructConfig) (Struct, error) {
	ds, err := ToStruct(a)
	if err != nil {
		return ds, err
	}
	if cfg.Name != "" {
		ds.Name = cfg.Name
	}

	if len(cfg.NameSpace) > 0 {
		ds.NameSpace = cfg.NameSpace
	}
	if cfg.UUID != "" {
		ds.UUID = cfg.UUID
	}

	if cfg.RemoveExistingTags {
		ds.Tags = nil
	}

	switch {
	case ds.Tags == nil && cfg.Tags != nil:
		ds.Tags = cfg.Tags
	case ds.Tags != nil && cfg.Tags != nil:
		for k, v := range cfg.Tags {
			ds.Tags[k] = v
		}
	}

	if cfg.RemoveExistingAttributes {
		ds.Attributes = nil
	}

	switch {
	case ds.Attributes == nil && cfg.Attributes != nil:
		ds.Attributes = cfg.Attributes
	case ds.Attributes != nil && cfg.Attributes != nil:
		for k, v := range cfg.Attributes {
			ds.Attributes[k] = v
		}
	}

	return ds, nil
}

func (s Struct) Identifier() string {
	ids := append(s.NameSpace, s.Name)
	for i, v := range ids {
		ids[i] = strings.ToLower(v)
	}
	return strings.Join(ids, s.NameSpaceSeparator)
}

// func ToStruct(value any) (Struct, error) {
// 	var pointer bool
// 	rt := reflect.TypeOf(value)
// 	// fmt.Printf("%T\t%s\t%s\n", value, value, rt)
// 	if rt.Kind() == reflect.Pointer {
// 		rt = rt.Elem()
// 	}
// 	rv := reflect.ValueOf(value)
// 	if t, ok := value.(reflect.Type); ok {
// 		rt = t
// 		if pointer = rt.Kind() == reflect.Pointer; pointer {
// 			rt = rt.Elem()
// 		}
// 		rv = reflect.New(rt).Elem()
// 	}
// 	if v, ok := value.(reflect.Value); ok {
// 		rv = v
// 		if !rv.IsValid() {
// 			return Struct{Value: rv}, fmt.Errorf("value is invalid: %T", value)
// 		}
// 		rt = v.Type()
// 		if rt.Kind() == reflect.Pointer {
// 			pointer = true
// 			rv = reflect.Indirect(rv)
// 			rt = rv.Type()
// 		}
// 	}
// 	if rv.IsZero() {
// 		rv = reflect.New(rt).Elem()
// 	}
// 	if rv.Kind() == reflect.Pointer {
// 		rv = rv.Elem()
// 	}

// 	if rt.Kind() != reflect.Struct {
// 		return Struct{Value: rv}, fmt.Errorf("%T is not a struct", value)
// 	}

// 	s := Struct{
// 		Name:       rt.Name(),
// 		Type:       rt,
// 		UUID:       strings.ReplaceAll(uuid.NewString(), "-", ""),
// 		Attributes: nil,
// 		Fields:     ToFields(value),
// 		pointer:    pointer,
// 	}
// 	children := []Structs{}
// 	for i, sf := range s.Fields {
// 		rvi := rv.Field(i)
// 		if rvi.Kind() == reflect.Pointer {
// 			rvi = rvi.Elem()
// 		}
// 		if rvi.Kind() == reflect.Invalid || rvi.IsZero() {
// 			rvi = reflect.New(sf.Type()).Elem()
// 		}
// 		// fmt.Printf("%s\t%T\t%s\n", fv.Type(), fv, rvi.Kind())
// 		if rvi.Kind() == reflect.Slice {
// 			// reflect.New(rt.Elem()).Elem()
// 			if rvi.Len() == 0 {
// 				reflect.Append(rvi, reflect.New(rvi.Type().Elem()).Elem())
// 			}
// 			fvs := []any{}
// 			for i := 0; i < rvi.Len(); i++ {
// 				fvs = append(fvs, rvi.Index(i).Interface())
// 			}
// 			// fmt.Printf("%s\t%T\t%s\t%s\n", sf.Name, fvs, len(fvs), rvi.Len())
// 			child, err := ToStructs(fvs)
// 			if err != nil {
// 				return Struct{}, err
// 			}
// 			children = append(children, child)
// 		}

// 	}
// 	s.Children = children

// 	return s, nil
// }
