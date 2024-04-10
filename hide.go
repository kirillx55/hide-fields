package hide_fields

import (
	"errors"
	"math"
	"reflect"
	"strconv"
)

const (
	tagHide = "hide"
)

var (
	errNotAPointer = errors.New("not a pointer")
	errUnknownType = errors.New("unknown type")
)

func HideFields(v any) error {
	if v == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errNotAPointer
	}

	return hideFields(reflect.ValueOf(v), false, "", false, reflect.Value{}, reflect.Value{})
}

func hideFields(vOf reflect.Value, hide bool, value string, isMap bool, m, key reflect.Value) error {
	switch vOf.Kind() {
	case reflect.Pointer:
		if err := hideFields(vOf.Elem(), hide, value, isMap, m, key); err != nil {
			return err
		}

	case reflect.Struct:
		vOfType := vOf.Type()

		if isMap {
			emptyStruct := reflect.New(vOf.Type()).Elem()
			var wasHide bool

			for i := 0; i < vOf.NumField(); i++ {
				emptyStruct.Field(i).Set(vOf.Field(i))

				hideValue, hideOk := vOfType.Field(i).Tag.Lookup(tagHide)

				v := value
				if hide || hideOk {
					wasHide = true
					if hideOk {
						v = hideValue
					}

					if err := hideFields(emptyStruct.Field(i), hide || hideOk, v, false, m, key); err != nil {
						return err
					}
				}
			}

			if wasHide {
				m.SetMapIndex(key, emptyStruct)
			}

			return nil
		}

		for i := 0; i < vOf.NumField(); i++ {
			hideValue, hideOk := vOfType.Field(i).Tag.Lookup(tagHide)

			v := value
			if hideOk {
				v = hideValue
			}

			if err := hideFields(vOf.Field(i), hide || hideOk, v, isMap, m, key); err != nil {
				return err
			}
		}

	case reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer, reflect.Uintptr, reflect.Invalid:

	case reflect.Map:
		keys := vOf.MapKeys()
		for i := range keys {
			if err := hideFields(vOf.MapIndex(keys[i]), hide, value, true, vOf, keys[i]); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < vOf.Len(); i++ {
			if err := hideFields(vOf.Index(i), hide, value, false, m, key); err != nil {
				return err
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Bool, reflect.String:
		if !hide {
			return nil
		}

		if isMap {
			m.SetMapIndex(key, getDefaultValue(m.MapIndex(key), value))

			return nil
		}

		if vOf.Equal(reflect.New(vOf.Type()).Elem()) {
			return nil
		}

		vOf.Set(getDefaultValue(vOf, value))

	default:
		return errUnknownType
	}

	return nil
}

func getDefaultValue(vOf reflect.Value, value string) reflect.Value {
	vOfKind := vOf.Kind()
	switch vOfKind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(value, 10, 64)

		switch vOfKind {
		case reflect.Int:
			if v >= int64(math.MinInt) && v <= int64(math.MaxInt) {
				return reflect.ValueOf(int(v))
			}

		case reflect.Int8:
			if v >= math.MinInt8 && v <= math.MaxInt8 {
				return reflect.ValueOf(int8(v))
			}

		case reflect.Int16:
			if v >= math.MinInt16 && v <= math.MaxInt16 {
				return reflect.ValueOf(int16(v))
			}

		case reflect.Int32:
			if v >= math.MinInt32 && v <= math.MaxInt32 {
				return reflect.ValueOf(int32(v))
			}

		default:
			return reflect.ValueOf(v)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(value, 10, 64)

		switch vOfKind {
		case reflect.Uint:
			if v <= uint64(math.MaxUint) {
				return reflect.ValueOf(uint(v))
			}

		case reflect.Uint8:
			if v <= math.MaxUint8 {
				return reflect.ValueOf(uint8(v))
			}

		case reflect.Uint16:
			if v <= math.MaxUint16 {
				return reflect.ValueOf(uint16(v))
			}

		case reflect.Uint32:
			if v <= math.MaxUint32 {
				return reflect.ValueOf(uint32(v))
			}

		default:
			return reflect.ValueOf(v)
		}

	case reflect.Float32, reflect.Float64:
		v, _ := strconv.ParseFloat(value, 64)

		switch vOfKind {
		case reflect.Float32:
			if v <= math.MaxFloat32 {
				return reflect.ValueOf(float32(v))
			}
		default:
			return reflect.ValueOf(v)
		}

	case reflect.Complex64, reflect.Complex128:
		v, _ := strconv.ParseComplex(value, 64)

		switch vOfKind {
		case reflect.Complex64:

		default:
			return reflect.ValueOf(v)
		}

	case reflect.Bool:
		v, _ := strconv.ParseBool(value)
		return reflect.ValueOf(v)

	case reflect.String:
		return reflect.ValueOf(value)
	}

	return reflect.New(vOf.Type()).Elem()
}
