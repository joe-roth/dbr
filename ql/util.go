package ql

import (
	"bytes"
	"reflect"
	"unicode"
)

func camelCaseToSnakeCase(name string) string {
	buf := new(bytes.Buffer)

	runes := []rune(name)

	for i := 0; i < len(runes); i++ {
		buf.WriteRune(unicode.ToLower(runes[i]))
		if !unicode.IsUpper(runes[i]) && i != len(runes)-1 && unicode.IsUpper(runes[i+1]) {
			buf.WriteRune('_')
		}
	}

	return buf.String()
}

func structMap(value reflect.Value) map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	structValue(m, value)
	return m
}

func structValue(m map[string]reflect.Value, value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		structValue(m, value.Elem())
	case reflect.Struct:
		t := value.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				// unexported
				continue
			}
			tag := field.Tag.Get("db")
			if tag == "-" {
				// ignore
				continue
			}
			if tag == "" {
				// no tag, but we can record the field name
				tag = camelCaseToSnakeCase(field.Name)
			}
			if _, ok := m[tag]; !ok {
				m[tag] = value.Field(i)
			}
			structValue(m, value.Field(i))
		}
	}
}
