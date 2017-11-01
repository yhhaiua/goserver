package goobjfmt

import (
	"encoding/binary"
	"math"
	"reflect"
)

type coder struct {
	order binary.ByteOrder
	buf   []byte
}

type decoder coder
type encoder coder

func (der *decoder) check(needlen int) bool {
	if needlen < 0 || len(der.buf) < needlen {
		return false
	}
	return true
}
func (der *decoder) bool() (bool, bool) {
	if !der.check(1) {
		return false, false
	}
	x := der.buf[0]
	der.buf = der.buf[1:]
	return x != 0, true
}

func (der *encoder) bool(x bool) {
	if x {
		der.buf[0] = 1
	} else {
		der.buf[0] = 0
	}
	der.buf = der.buf[1:]
}

func (der *decoder) uint8() (uint8, bool) {
	if !der.check(1) {
		return 0, false
	}
	x := der.buf[0]
	der.buf = der.buf[1:]
	return x, true
}

func (der *encoder) uint8(x uint8) {
	der.buf[0] = x
	der.buf = der.buf[1:]
}

func (der *decoder) uint16() (uint16, bool) {
	if !der.check(2) {
		return 0, false
	}
	x := der.order.Uint16(der.buf[0:2])
	der.buf = der.buf[2:]
	return x, true
}

func (der *encoder) uint16(x uint16) {
	der.order.PutUint16(der.buf[0:2], x)
	der.buf = der.buf[2:]
}

func (der *decoder) uint32() (uint32, bool) {
	if !der.check(4) {
		return 0, false
	}
	x := der.order.Uint32(der.buf[0:4])
	der.buf = der.buf[4:]
	return x, true
}

func (der *encoder) uint32(x uint32) {
	der.order.PutUint32(der.buf[0:4], x)
	der.buf = der.buf[4:]
}

func (der *decoder) uint64() (uint64, bool) {
	if !der.check(8) {
		return 0, false
	}
	x := der.order.Uint64(der.buf[0:8])
	der.buf = der.buf[8:]
	return x, true
}

func (der *encoder) uint64(x uint64) {
	der.order.PutUint64(der.buf[0:8], x)
	der.buf = der.buf[8:]
}
func (der *decoder) bytes() ([]byte, bool) {
	l, ok := der.int32()
	if l > 0 && ok {
		if !der.check(int(l)) {
			return nil, false
		}
		buf := make([]byte, l)

		copy(buf, der.buf[0:l])
		der.buf = der.buf[l:]
		return buf, true
	}

	return nil, false
}

func (der *encoder) bytes(x []byte) {
	l := len(x)
	der.int32(int32(l))
	copy(der.buf, []byte(x))
	der.buf = der.buf[l:]
}

func (der *decoder) int8() (int8, bool) {
	data, ok := der.uint8()
	return int8(data), ok
}

func (der *encoder) int8(x int8) { der.uint8(uint8(x)) }

func (der *decoder) int16() (int16, bool) {
	data, ok := der.uint16()
	return int16(data), ok
}

func (der *encoder) int16(x int16) { der.uint16(uint16(x)) }

func (der *decoder) int32() (int32, bool) {
	data, ok := der.uint32()
	return int32(data), ok
}

func (der *encoder) int32(x int32) { der.uint32(uint32(x)) }

func (der *decoder) int64() (int64, bool) {
	data, ok := der.uint64()
	return int64(data), ok
}

func (der *encoder) int64(x int64) { der.uint64(uint64(x)) }

func (der *decoder) value(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			ok := der.value(v.Index(i))
			if !ok {
				return false
			}
		}
		return true

	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			// Note: Calling v.CanSet() below is an optimization.
			// It would be sufficient to check the field name,
			// but creating the StructField info for each field is
			// costly (run "go test -bench=ReadStruct" and compare
			// results when making changes to this code).
			vv := v.Field(i)
			ok := der.value(vv)
			if !ok {
				return false
			}
		}
		return true
	case reflect.String:
		data, ok := der.bytes()
		if ok {
			v.SetString(string(data))
		}
		return ok

	case reflect.Slice:
		var ok bool
		if v.Type().Elem().Kind() == reflect.Uint8 {
			data, ok := der.bytes()
			if ok {
				v.SetBytes(data)
			}

		} else {
			l, ok := der.int32()
			if ok && l > 0 {
				tempslice := reflect.MakeSlice(v.Type(), 1, 1)
				onesize, _ := dataSize(tempslice.Index(0), nil)
				ok = der.check(onesize * int(l))
				if ok {
					if l == 1 {
						sliceValue := reflect.New(tempslice.Type().Elem()).Elem()
						ok = der.value(sliceValue)
						if ok {
							tempslice.Index(0).Set(sliceValue)
							v.Set(tempslice)
						}

					} else {

						slice := reflect.MakeSlice(v.Type(), int(l), int(l))

						for i := 0; i < int(l); i++ {
							sliceValue := reflect.New(slice.Type().Elem()).Elem()
							ok = der.value(sliceValue)
							if !ok {
								return false
							}
							slice.Index(i).Set(sliceValue)
						}
						v.Set(slice)
					}

				}

			}

		}
		return ok
	case reflect.Bool:
		data, ok := der.bool()
		if ok {
			v.SetBool(data)
		}
		return ok
	case reflect.Int8:
		data, ok := der.int8()
		if ok {
			v.SetInt(int64(data))
		}
		return ok
	case reflect.Int16:
		data, ok := der.int16()
		if ok {
			v.SetInt(int64(data))
		}
		return ok
	case reflect.Int32:
		data, ok := der.int32()
		if ok {
			v.SetInt(int64(data))
		}
		return ok
	case reflect.Int64:
		data, ok := der.int64()
		if ok {
			v.SetInt(data)
		}
		return ok

	case reflect.Uint8:
		data, ok := der.uint8()
		if ok {
			v.SetUint(uint64(data))
		}
		return ok
	case reflect.Uint16:
		data, ok := der.uint16()
		if ok {
			v.SetUint(uint64(data))
		}
		return ok
	case reflect.Uint32:
		data, ok := der.uint32()
		if ok {
			v.SetUint(uint64(data))
		}
		return ok
	case reflect.Uint64:
		data, ok := der.uint64()
		if ok {
			v.SetUint(data)
		}
		return ok

	case reflect.Float32:
		data, ok := der.uint32()
		if ok {
			v.SetFloat(float64(math.Float32frombits(data)))
		}
		return ok

	case reflect.Float64:
		data, ok := der.uint64()
		if ok {
			v.SetFloat(math.Float64frombits(data))
		}
		return ok
	case reflect.Interface:

	//case reflect.Ptr:
	//
	//	valuePtr := reflect.New(v.Type().Elem())
	//
	//	der.value(valuePtr.Elem())
	//
	//	v.Set(valuePtr)

	default:
		panic("encode: unsupport kind: " + v.Kind().String())
	}
	return false
}

func (der *encoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			der.value(v.Index(i))
		}
	case reflect.Interface:
		der.value(v.Elem())
	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {

			v := v.Field(i)

			der.value(v)
		}

	case reflect.Slice:

		if v.Type().Elem().Kind() == reflect.Uint8 {

			der.bytes(v.Bytes())

		} else {
			l := v.Len()
			der.int32(int32(l))
			for i := 0; i < l; i++ {
				der.value(v.Index(i))
			}
		}

	case reflect.String:
		der.bytes([]byte(v.String()))

	case reflect.Bool:
		der.bool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v.Type().Kind() {
		case reflect.Int8:
			der.int8(int8(v.Int()))
		case reflect.Int16:
			der.int16(int16(v.Int()))
		case reflect.Int32:
			der.int32(int32(v.Int()))
		case reflect.Int64:
			der.int64(v.Int())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v.Type().Kind() {
		case reflect.Uint8:
			der.uint8(uint8(v.Uint()))
		case reflect.Uint16:
			der.uint16(uint16(v.Uint()))
		case reflect.Uint32:
			der.uint32(uint32(v.Uint()))
		case reflect.Uint64:
			der.uint64(v.Uint())
		}

	case reflect.Float32, reflect.Float64:
		switch v.Type().Kind() {
		case reflect.Float32:
			der.uint32(math.Float32bits(float32(v.Float())))
		case reflect.Float64:
			der.uint64(math.Float64bits(v.Float()))
		}
	case reflect.Ptr:
		der.value(v.Elem())
	//case reflect.Invalid:

	default:
		panic("encode: unsupport kind: " + v.Kind().String())
	}
}
