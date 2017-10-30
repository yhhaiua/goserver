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

func (der *decoder) bool() bool {
	x := der.buf[0]
	der.buf = der.buf[1:]
	return x != 0
}

func (der *encoder) bool(x bool) {
	if x {
		der.buf[0] = 1
	} else {
		der.buf[0] = 0
	}
	der.buf = der.buf[1:]
}

func (der *decoder) uint8() uint8 {
	x := der.buf[0]
	der.buf = der.buf[1:]
	return x
}

func (der *encoder) uint8(x uint8) {
	der.buf[0] = x
	der.buf = der.buf[1:]
}

func (der *decoder) uint16() uint16 {
	x := der.order.Uint16(der.buf[0:2])
	der.buf = der.buf[2:]
	return x
}

func (der *encoder) uint16(x uint16) {
	der.order.PutUint16(der.buf[0:2], x)
	der.buf = der.buf[2:]
}

func (der *decoder) uint32() uint32 {
	x := der.order.Uint32(der.buf[0:4])
	der.buf = der.buf[4:]
	return x
}

func (der *encoder) uint32(x uint32) {
	der.order.PutUint32(der.buf[0:4], x)
	der.buf = der.buf[4:]
}

func (der *decoder) uint64() uint64 {
	x := der.order.Uint64(der.buf[0:8])
	der.buf = der.buf[8:]
	return x
}

func (der *encoder) uint64(x uint64) {
	der.order.PutUint64(der.buf[0:8], x)
	der.buf = der.buf[8:]
}
func (der *decoder) bytes() []byte {
	l := der.int32()
	buf := make([]byte, l)

	copy(buf, der.buf[0:l])
	der.buf = der.buf[l:]
	return buf
}

func (der *encoder) bytes(x []byte) {
	l := len(x)
	der.int32(int32(l))
	copy(der.buf, []byte(x))
	der.buf = der.buf[l:]
}

func (der *decoder) int8() int8 { return int8(der.uint8()) }

func (der *encoder) int8(x int8) { der.uint8(uint8(x)) }

func (der *decoder) int16() int16 { return int16(der.uint16()) }

func (der *encoder) int16(x int16) { der.uint16(uint16(x)) }

func (der *decoder) int32() int32 { return int32(der.uint32()) }

func (der *encoder) int32(x int32) { der.uint32(uint32(x)) }

func (der *decoder) int64() int64 { return int64(der.uint64()) }

func (der *encoder) int64(x int64) { der.uint64(uint64(x)) }

func (der *decoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			der.value(v.Index(i))
		}

	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			// Note: Calling v.CanSet() below is an optimization.
			// It would be sufficient to check the field name,
			// but creating the StructField info for each field is
			// costly (run "go test -bench=ReadStruct" and compare
			// results when making changes to this code).
			vv := v.Field(i)
			der.value(vv)
		}
	case reflect.String:
		v.SetString(string(der.bytes()))

	case reflect.Slice:

		if v.Type().Elem().Kind() == reflect.Uint8 {

			v.SetBytes(der.bytes())

		} else {
			l := int(der.int32())
			slice := reflect.MakeSlice(v.Type(), l, l)

			for i := 0; i < l; i++ {

				sliceValue := reflect.New(slice.Type().Elem()).Elem()

				der.value(sliceValue)
				slice.Index(i).Set(sliceValue)
			}

			v.Set(slice)
		}

	case reflect.Bool:
		v.SetBool(der.bool())

	case reflect.Int8:
		v.SetInt(int64(der.int8()))
	case reflect.Int16:
		v.SetInt(int64(der.int16()))
	case reflect.Int32:
		v.SetInt(int64(der.int32()))
	case reflect.Int64:
		v.SetInt(der.int64())

	case reflect.Uint8:
		v.SetUint(uint64(der.uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(der.uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(der.uint32()))
	case reflect.Uint64:
		v.SetUint(der.uint64())

	case reflect.Float32:
		v.SetFloat(float64(math.Float32frombits(der.uint32())))
	case reflect.Float64:
		v.SetFloat(math.Float64frombits(der.uint64()))
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
