package dbr

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
func (b *InsertBuilder) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	if len(b.Column) == 0 {
		return ErrColumnNotSpecified
	}

	buf.WriteString("INSERT INTO ")
	buf.WriteString(d.QuoteIdent(b.Table))

	buf.WriteString(" (")

	placeholder := new(bytes.Buffer)
	placeholder.WriteRune('(')
	for i, col := range b.Column {
		if i > 0 {
			buf.WriteString(",")
			placeholder.WriteString(",")
		}
		buf.WriteString(d.QuoteIdent(col))
		placeholder.WriteString(d.Placeholder())
	}
	placeholder.WriteString(")")
	placeholderStr := placeholder.String()

	buf.WriteString(") VALUES ")

	for i, tuple := range b.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(placeholderStr)

		buf.WriteValue(tuple...)
	}

	return nil
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
	v := reflect.Indirect(reflect.ValueOf(structValue))

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
