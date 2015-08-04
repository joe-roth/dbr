package ql

import "bytes"

// DeleteBuilder builds `DELETE ...`
type DeleteBuilder struct {
	raw

	Table string

	WhereCond []Condition
}

// Build builds `DELETE ...` in dialect
func (b *DeleteBuilder) Build(d Dialect) (string, []interface{}, error) {
	if b.raw.Query != "" {
		return b.raw.Query, b.raw.Value, nil
	}

	if b.Table == "" {
		return "", nil, ErrTableNotSpecified
	}

	buf := new(bytes.Buffer)
	var value []interface{}

	buf.WriteString("DELETE FROM ")
	buf.WriteString(d.QuoteIdent(b.Table))

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

// DeleteFrom creates a DeleteBuilder
func DeleteFrom(table string) *DeleteBuilder {
	return &DeleteBuilder{
		Table: table,
	}
}

// DeleteBySQL creates a DeleteBuilder from raw query
func DeleteBySQL(query string, value ...interface{}) *DeleteBuilder {
	return &DeleteBuilder{
		raw: raw{
			Query: query,
			Value: value,
		},
	}
}

// Where adds a where condition
func (b *DeleteBuilder) Where(query interface{}, value ...interface{}) *DeleteBuilder {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Condition:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}
