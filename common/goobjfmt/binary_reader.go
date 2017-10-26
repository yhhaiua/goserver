package goobjfmt

import (
	"encoding/binary"
	"reflect"
)

//BinaryRead 从data读取
func BinaryRead(data []byte, obj interface{}) error {

	if len(data) == 0 {
		return nil
	}

	v := reflect.ValueOf(obj)

	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	}

	size := dataSize(v, nil)
	if size < 0 {
		return errInvalidType
	}

	if len(data) < size {
		return errOutOfData
	}

	d := &decoder{order: binary.LittleEndian, buf: data}
	d.value(v)

	return nil
}
