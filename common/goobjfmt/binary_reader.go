package goobjfmt

import (
	"encoding/binary"
	"reflect"
)

//BinaryRead 从data读取
func BinaryRead(data []byte, obj interface{}) error {

	if len(data) == 0 {
		return errGetDatalen
	}

	v := reflect.ValueOf(obj)

	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	}

	size, indefinite := dataSize(v, nil)
	if size < 0 {
		return errInvalidType
	}

	if (len(data) != size && !indefinite) || (indefinite && len(data) < size) {
		return errOutOfData
	}

	d := &decoder{order: binary.LittleEndian, buf: data}
	ok := d.value(v)
	if !ok {
		return errOfData
	}
	return nil
}
