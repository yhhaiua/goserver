package goobjfmt

import (
	"encoding/binary"
	"errors"
	"reflect"
)

var (
	errInvalidType = errors.New("invalid type")
	errOutOfData   = errors.New("length inconsistent")
	errGetDatalen  = errors.New("length 0")
	errOfData      = errors.New("length error")
)

//BinaryWrite 写入data
func BinaryWrite(obj interface{}) ([]byte, error) {

	// Fallback to reflect-based encoding.
	v := reflect.Indirect(reflect.ValueOf(obj))
	size, _ := dataSize(v, nil)
	if size < 0 {
		return nil, errInvalidType
	}

	buf := make([]byte, size)

	e := &encoder{order: binary.LittleEndian, buf: buf}
	e.value(v)

	return buf, nil
}
