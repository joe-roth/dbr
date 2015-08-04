package ql

import (
	"database/sql"
	"reflect"
)

// Load loads any value from sql.Rows
func Load(rows *sql.Rows, value interface{}) error {
	defer rows.Close()

	column, err := rows.Columns()
	if err != nil {
		return err
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return ErrBadArgument
	}
	v = v.Elem()
	isSlice := v.Kind() == reflect.Slice && v.Type().Elem().Kind() != reflect.Uint8
	for rows.Next() {
		elem := v
		if isSlice {
			elem = reflect.New(v.Type().Elem()).Elem()
		}
		ptr, err := findPtr(column, elem)
		if err != nil {
			return err
		}
		rows.Scan(ptr...)
		if isSlice {
			v.Set(reflect.Append(v, elem))
		} else {
			break
		}
	}
	return nil
}

func findPtr(column []string, value reflect.Value) ([]interface{}, error) {
	switch value.Kind() {
	case reflect.Struct:
		var ptr []interface{}
		m := structMap(value)
		for _, key := range column {
			if val, ok := m[key]; ok {
				ptr = append(ptr, val.Addr().Interface())
			} else {
				ptr = append(ptr, nil)
			}
		}
		return ptr, nil
	case reflect.Ptr:
		if value.CanSet() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		return findPtr(column, value.Elem())
	}
	return []interface{}{value.Addr().Interface()}, nil
}
