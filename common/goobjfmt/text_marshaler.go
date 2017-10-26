package goobjfmt

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
)

// TextMarshaler is a configurable text format marshaler.
type TextMarshaler struct {
	Compact        bool // use compact text format (one line).
	IgnoreDefault  bool // Do not output value when value equals its default value
	OriginalString bool // Do not output string as byte
}

//Marshal 数据写入
func (der *TextMarshaler) Marshal(w io.Writer, obj interface{}) error {

	val := reflect.ValueOf(obj)

	if obj == nil || (val.Kind() == reflect.Ptr && val.IsNil()) {
		w.Write([]byte("<nil>"))
		return nil
	}

	var bw *bufio.Writer
	ww, ok := w.(writer)
	if !ok {
		bw = bufio.NewWriter(w)
		ww = bw
	}
	aw := &textWriter{
		w:        ww,
		complete: true,
		compact:  der.Compact,
	}

	// Dereference the received pointer so we don't have outer < and >.
	v := reflect.Indirect(val)
	if err := der.writeStruct(aw, v); err != nil {
		return err
	}
	if bw != nil {
		return bw.Flush()
	}

	return nil
}

func writeName(w *textWriter, name string) error {
	if _, err := w.WriteString(name); err != nil {
		return err
	}

	return w.WriteByte(':')
}

func (der *TextMarshaler) ignoreField(w *textWriter, v reflect.Value) bool {

	if !der.IgnoreDefault {
		return false
	}

	v = reflect.Indirect(v)

	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return true
		}
	case reflect.Int32, reflect.Int64, reflect.Int:
		if v.Int() == 0 {
			return true
		}
	case reflect.Uint32, reflect.Uint64, reflect.Uint:
		if v.Uint() == 0 {
			return true
		}
	case reflect.Bool:
		if !v.Bool() {
			return true
		}

	case reflect.Slice, reflect.String:
		if v.Len() == 0 {
			return true
		}
	}

	return false
}

func (der *TextMarshaler) writeAny(w *textWriter, v reflect.Value) error {
	v = reflect.Indirect(v)

	// Floats have special cases.
	if v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
		x := v.Float()
		var b []byte
		switch {
		case math.IsInf(x, 1):
			b = posInf
		case math.IsInf(x, -1):
			b = negInf
		case math.IsNaN(x):
			b = nan
		}
		if b != nil {
			_, err := w.Write(b)
			return err
		}
		// Other values are handled below.
	}

	// We don't attempt to serialise every possible value type; only those
	// that can occur in protocol buffers.

	switch v.Kind() {
	case reflect.Slice:
		// Should only be a []byte; repeated fields are handled in writeStruct.
		if err := writeString(w, string(v.Bytes())); err != nil {
			return err
		}
	case reflect.String:

		if err := writeString(w, v.String()); err != nil {
			return err
		}
	case reflect.Struct:
		// Required/optional group/message.
		var bra, ket byte = '<', '>'

		if err := w.WriteByte(bra); err != nil {
			return err
		}
		if !w.compact {
			if err := w.WriteByte('\n'); err != nil {
				return err
			}
		}
		w.indent()
		if etm, ok := v.Interface().(encoding.TextMarshaler); ok {
			text, err := etm.MarshalText()
			if err != nil {
				return err
			}
			if _, err = w.Write(text); err != nil {
				return err
			}
		} else if err := der.writeStruct(w, v); err != nil {
			return err
		}
		w.unindent()
		if err := w.WriteByte(ket); err != nil {
			return err
		}
	default:
		_, err := fmt.Fprint(w, v.Interface())
		return err
	}

	return nil
}

func (der *TextMarshaler) writeStruct(w *textWriter, sv reflect.Value) error {

	st := sv.Type()

	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)

		stf := st.Field(i)

		if tag := stf.Tag.Get("text"); tag == "-" {
			continue
		}

		name := stf.Name

		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			// Field not filled in. This could be an optional field or
			// a required field that wasn't filled in. Either way, there
			// isn't anything we can show for it.
			continue
		}

		if fv.Kind() == reflect.Slice {

			// Repeated field that is empty, or a bytes field that is unused.
			if fv.IsNil() {
				continue
			}

			// []byte
			if fv.Len() > 0 && fv.Index(0).Kind() == reflect.Uint8 {

				if err := writeName(w, name); err != nil {
					return err
				}

				if err := w.WriteByte('['); err != nil {
					return err
				}
			}

			// Repeated field.
			for j := 0; j < fv.Len(); j++ {

				v := fv.Index(j)

				if v.Kind() != reflect.Uint8 {
					if err := writeName(w, name); err != nil {
						return err
					}
				}

				if !w.compact {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
				}

				if v.Kind() == reflect.Ptr && v.IsNil() {
					// A nil message in a repeated field is not valid,
					// but we can handle that more gracefully than panicking.
					if _, err := w.Write([]byte("<nil>\n")); err != nil {
						return err
					}
					continue
				}
				if err := der.writeAny(w, v); err != nil {
					return err
				}
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}

			if fv.Len() > 0 && fv.Index(0).Kind() == reflect.Uint8 {
				if _, err := w.WriteString("] "); err != nil {
					return err
				}
			}

			continue
		}

		if fv.Kind() == reflect.Map {
			// Map fields are rendered as a repeated struct with key/value fields.
			keys := fv.MapKeys()
			sort.Sort(mapKeys(keys))
			for _, key := range keys {
				val := fv.MapIndex(key)
				if err := writeName(w, name); err != nil {
					return err
				}
				if !w.compact {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
				}
				// open struct
				if err := w.WriteByte('<'); err != nil {
					return err
				}
				if !w.compact {
					if err := w.WriteByte('\n'); err != nil {
						return err
					}
				}
				w.indent()
				// key
				if _, err := w.WriteString("key:"); err != nil {
					return err
				}
				if !w.compact {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
				}
				if err := der.writeAny(w, key); err != nil {
					return err
				}
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
				// nil values aren't legal, but we can avoid panicking because of them.
				if val.Kind() != reflect.Ptr || !val.IsNil() {
					// value
					if _, err := w.WriteString("value:"); err != nil {
						return err
					}
					if !w.compact {
						if err := w.WriteByte(' '); err != nil {
							return err
						}
					}
					if err := der.writeAny(w, val); err != nil {
						return err
					}
					if err := w.WriteByte('\n'); err != nil {
						return err
					}
				}
				// close struct
				w.unindent()
				if err := w.WriteByte('>'); err != nil {
					return err
				}
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
			continue
		}

		if !der.ignoreField(w, fv) {
			if err := writeName(w, name); err != nil {
				return err
			}
			if !w.compact {
				if err := w.WriteByte(' '); err != nil {
					return err
				}
			}

			// Enums have a String method, so writeAny will work fine.
			if err := der.writeAny(w, fv); err != nil {
				return err
			}

			if err := w.WriteByte('\n'); err != nil {
				return err
			}

		}

	}

	return nil
}

//Text 转成strings
func (der *TextMarshaler) Text(obj interface{}) string {
	var buf bytes.Buffer
	der.Marshal(&buf, obj)
	return buf.String()
}

var (
	defaultTextMarshaler = TextMarshaler{}
	compactTextMarshaler = TextMarshaler{Compact: true, IgnoreDefault: true}
)

//MarshalTextString String
func MarshalTextString(obj interface{}) string {
	return defaultTextMarshaler.Text(obj)
}

//CompactTextString String
func CompactTextString(obj interface{}) string { return compactTextMarshaler.Text(obj) }
