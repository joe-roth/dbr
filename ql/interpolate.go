package ql

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Interpolate replaces placeholder in query with corresponding value in dialect
func Interpolate(query string, value []interface{}, d Dialect) (string, error) {
	buf := new(bytes.Buffer)
	err := interpolate(query, value, d, buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func interpolate(query string, value []interface{}, d Dialect, w StringWriter) error {
	placeholder := d.Placeholder()

	if strings.Count(query, placeholder) != len(value) {
		return ErrBadArgument
	}

	valueIndex := 0

	for {
		index := strings.Index(query, placeholder)
		if index == -1 {
			break
		}
		w.WriteString(query[:index])
		query = query[index+len(placeholder):]

		err := encodePlaceholder(value[valueIndex], d, w)
		if err != nil {
			return err
		}

		valueIndex++
	}

	// placeholder not found; write remaining query
	w.WriteString(query)

	return nil
}

func encodePlaceholder(value interface{}, d Dialect, w StringWriter) error {
	if builder, ok := value.(Builder); ok {
		buf := NewBuffer()
		err := builder.Build(d, buf)
		if err != nil {
			return err
		}
		// subquery
		_, ok := value.(*SelectBuilder)
		if ok {
			w.WriteString("(")
		}
		err = interpolate(buf.String(), buf.Value(), d, w)
		if err != nil {
			return err
		}
		if ok {
			w.WriteString(")")
		}
		return nil
	}

	if valuer, ok := value.(driver.Valuer); ok {
		// get driver.Valuer's data
		var err error
		value, err = valuer.Value()
		if err != nil {
			return err
		}
	}

	if value == nil {
		w.WriteString("NULL")
		return nil
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		w.WriteString(d.EncodeString(v.String()))
		return nil
	case reflect.Bool:
		w.WriteString(d.EncodeBool(v.Bool()))
		return nil
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		// TODO: verify this works
		w.WriteString(fmt.Sprint(v.Interface()))
		return nil
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			w.WriteString(d.EncodeTime(v.Interface().(time.Time).UTC()))
			return nil
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			w.WriteString(d.EncodeBytes(v.Bytes()))
			return nil
		}
		if v.Len() == 0 {
			// This will never match, since nothing is equal to null (not even null itself.)
			w.WriteString("(NULL)")
		} else {
			w.WriteString("(")
			for i := 0; i < v.Len(); i++ {
				if i > 0 {
					w.WriteString(",")
				}
				err := encodePlaceholder(v.Index(i).Interface(), d, w)
				if err != nil {
					return err
				}
			}
			w.WriteString(")")
		}
		return nil
	case reflect.Ptr:
		return encodePlaceholder(v.Elem().Interface(), d, w)
	}
	return ErrNotSupported
}
