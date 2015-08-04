package ql

import (
	"bytes"
	"reflect"
)

// InsertBuilder builds `INSERT INTO ...`
type InsertBuilder struct {
	raw

	Table  string
	Column []string
	Value  [][]interface{}
}

// Build builds `INSERT INTO ...` in dialect
func (b *InsertBuilder) Build(d Dialect) (string, []interface{}, error) {
	if b.raw.Query != "" {
		return b.raw.Query, b.raw.Value, nil
	}

	if b.Table == "" {
		return "", nil, ErrTableNotSpecified
	}

	if len(b.Column) == 0 {
		return "", nil, ErrColumnNotSpecified
	}

	buf := new(bytes.Buffer)
	var value []interface{}

	buf.WriteString("INSERT INTO ")
	buf.WriteString(d.QuoteIdent(b.Table))

	buf.WriteString(" (")

	placeholder := new(bytes.Buffer)
	placeholder.WriteRune('(')
	for i, col := range b.Column {
		if i > 0 {
			buf.WriteRune(',')
			placeholder.WriteRune(',')
		}
		buf.WriteString(d.QuoteIdent(col))
		placeholder.WriteString(d.Placeholder())
	}
	placeholder.WriteRune(')')
	placeholderStr := placeholder.String()

	buf.WriteString(") VALUES ")

	for i, list := range b.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(placeholderStr)

		for _, v := range list {
			value = append(value, v)
		}
	}

	return buf.String(), value, nil
}

// InsertInto creates an InsertBuilder
func InsertInto(table string) *InsertBuilder {
	return &InsertBuilder{
		Table: table,
	}
}

// InsertBySQL creates an InsertBuilder from raw query
func InsertBySQL(query string, value ...interface{}) *InsertBuilder {
	return &InsertBuilder{
		raw: raw{
			Query: query,
			Value: value,
		},
	}
}

// Columns adds columns
func (b *InsertBuilder) Columns(column ...string) *InsertBuilder {
	b.Column = column
	return b
}

// Values adds a tuple for columns
func (b *InsertBuilder) Values(value ...interface{}) *InsertBuilder {
	b.Value = append(b.Value, value)
	return b
}

// Record adds a tuple for columns from a struct
func (b *InsertBuilder) Record(structValue interface{}) *InsertBuilder {
	v := reflect.ValueOf(structValue)
	v = reflect.Indirect(v)

	if v.Kind() == reflect.Struct {
		var value []interface{}
		m := structMap(v)
		for _, key := range b.Column {
			if val, ok := m[key]; ok {
				value = append(value, val.Interface())
			} else {
				value = append(value, nil)
			}
		}
		b.Values(value...)
	}
	return b
}
