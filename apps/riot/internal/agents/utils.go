package agents

import (
	"reflect"
)

// DeepCopy creates a deep clone of any value using reflection.
func DeepCopy[T any](src T) T {
	orig := reflect.ValueOf(src)
	return deepCopyValue(orig).Interface().(T)
}

func deepCopyValue(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		clonedPtr := reflect.New(v.Type().Elem())
		clonedPtr.Elem().Set(deepCopyValue(v.Elem()))
		return clonedPtr

	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		clonedSlice := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			clonedSlice.Index(i).Set(deepCopyValue(v.Index(i)))
		}
		return clonedSlice

	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		clonedMap := reflect.MakeMapWithSize(v.Type(), v.Len())
		for _, key := range v.MapKeys() {
			clonedMap.SetMapIndex(deepCopyValue(key), deepCopyValue(v.MapIndex(key)))
		}
		return clonedMap

	case reflect.Struct:
		clonedStruct := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if clonedStruct.Field(i).CanSet() {
				clonedStruct.Field(i).Set(deepCopyValue(v.Field(i)))
			}
		}
		return clonedStruct

	case reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		return deepCopyValue(v.Elem())

	default:
		// Primitive types (int, string, bool, etc.) are passed by value
		return v
	}
}
