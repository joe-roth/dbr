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
	placeholder := d.Placeholder()

	if strings.Count(query, placeholder) != len(value) {
		return "", ErrBadArgument
	}

	buf := new(bytes.Buffer)
	valueIndex := 0

	for {
		index := strings.Index(query, placeholder)
		if index == -1 {
			break
		}
		buf.WriteString(query[:index])
		query = query[index+len(placeholder):]

		s, err := interpolateWithDialect(value[valueIndex], d)
		if err != nil {
			return "", err
		}
		buf.WriteString(s)

		valueIndex++
	}

	// placeholder not found; write remaining query
	buf.WriteString(query)

	return buf.String(), nil
}

// return literal for different dialect
func interpolateWithDialect(value interface{}, d Dialect) (string, error) {
	if builder, ok := value.(Builder); ok {
		s, v, err := builder.Build(d)
		if err != nil {
			return "", err
		}
		s, err = Interpolate(s, v, d)
		if err != nil {
			return "", err
		}
		// subquery
		if _, ok := value.(*SelectBuilder); ok {
			return "(" + s + ")", nil
		}
		return s, err
	}

	if valuer, ok := value.(driver.Valuer); ok {
		// get driver.Valuer's data
		var err error
		value, err = valuer.Value()
		if err != nil {
			return "", err
		}
	}

	if value == nil {
		return "NULL", nil
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return d.EncodeString(v.String()), nil
	case reflect.Bool:
		return d.EncodeBool(v.Bool()), nil
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
		return fmt.Sprint(v.Interface()), nil
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return d.EncodeTime(v.Interface().(time.Time)), nil
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			return d.EncodeBytes(v.Bytes()), nil
		}
		buf := new(bytes.Buffer)
		buf.WriteRune('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteRune(',')
			}
			s, err := interpolateWithDialect(v.Index(i).Interface(), d)
			if err != nil {
				return "", err
			}
			buf.WriteString(s)
		}
		buf.WriteRune(')')
		return buf.String(), nil
	case reflect.Ptr:
		return interpolateWithDialect(v.Elem().Interface(), d)
	}
	return "", ErrNotSupported
}
