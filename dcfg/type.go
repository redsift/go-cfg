package dcfg

import (
	"reflect"
	"strings"
)

// Type is a custom string alias used to represent a type as a predictable string.
type Type string

// TypeOfValue takes a value of any type and returns its predictable Type representation.
// This generic function uses a workaround (wrapping the value in a slice) to retrieve the actual
// type given. Without this, even if explicitly instantiated with an interface type, a concrete
// type is returned by reflect.TypeOfValue.
func TypeOfValue[T any](t T) Type {
	return TypeOfReflect(
		reflect.TypeOf([]T{t}).Elem(),
	)
}

// TypeOf is a convenience function to generate a type without a value.
func TypeOf[T any]() Type {
	return TypeOfReflect(
		reflect.TypeOf([]T{}).Elem(),
	)
}

// TypeOfReflect returns a predictable type string for the given reflect.Type.
func TypeOfReflect(t reflect.Type) Type {
	switch t.Kind() {
	case reflect.Pointer:
		return TypeOfReflect(t.Elem())

	case reflect.Array, reflect.Slice:
		return TypeOfReflect(t.Elem()) + "_slice"

	case reflect.Map:
		return "map[" + TypeOfReflect(t.Key()) + "]" + TypeOfReflect(t.Elem())

	default:
		name := t.Name()
		if name == "" {
			panic("cannot get type name")
		}
		pkgPath := strings.TrimPrefix(t.PkgPath(), "github.com/redsift/")
		if pkgPath != "" {
			pkgPath += "/"
		}
		return Type(pkgPath + name)
	}
}
