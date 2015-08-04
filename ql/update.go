package ql

import "bytes"

// UpdateBuilder builds `UPDATE ...`
type UpdateBuilder struct {
	raw

	Table string
	Value map[string]interface{}

	WhereCond []Condition
}

// Build builds `UPDATE ...` in dialect
func (b *UpdateBuilder) Build(d Dialect) (string, []interface{}, error) {
	if b.raw.Query != "" {
		return b.raw.Query, b.raw.Value, nil
	}

	if b.Table == "" {
		return "", nil, ErrTableNotSpecified
	}

	if len(b.Value) == 0 {
		return "", nil, ErrColumnNotSpecified
	}

	buf := new(bytes.Buffer)
	var value []interface{}

	buf.WriteString("UPDATE ")
	buf.WriteString(d.QuoteIdent(b.Table))
	buf.WriteString(" SET ")

	i := 0
	for col, v := range b.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(d.QuoteIdent(col))
		buf.WriteString(" = ")
		buf.WriteString(d.Placeholder())

		value = append(value, v)
		i++
	}

	if len(b.WhereCond) > 0 {
		query, v, err := And(b.WhereCond...).Build(d)
		if err != nil {
			return "", nil, err
		}
		if query != "" {
			buf.WriteString(" WHERE ")
			buf.WriteString(query)

			value = append(value, v...)
		}
	}
	return buf.String(), value, nil
}

// Update creates an UpdateBuilder
func Update(table string) *UpdateBuilder {
	return &UpdateBuilder{
		Table: table,
		Value: make(map[string]interface{}),
	}
}

// UpdateBySQL creates an UpdateBuilder with raw query
func UpdateBySQL(query string, value ...interface{}) *UpdateBuilder {
	return &UpdateBuilder{
		raw: raw{
			Query: query,
			Value: value,
		},
		Value: make(map[string]interface{}),
	}
}

// Where adds a where condition
func (b *UpdateBuilder) Where(query interface{}, value ...interface{}) *UpdateBuilder {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Condition:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}

// Set specifies a key-value pair
func (b *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	b.Value[column] = value
	return b
}

// SetMap specifies a list of key-value pair
func (b *UpdateBuilder) SetMap(m map[string]interface{}) *UpdateBuilder {
	for col, val := range m {
		b.Set(col, val)
	}
	return b
}
