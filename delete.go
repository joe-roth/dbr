package dbr

// DeleteBuilder builds `DELETE ...`
type DeleteBuilder struct {
	raw

	Table string

	WhereCond []Condition
}

// Build builds `DELETE ...` in dialect
func (b *DeleteBuilder) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	buf.WriteString("DELETE FROM ")
	buf.WriteString(d.QuoteIdent(b.Table))

	if len(b.WhereCond) > 0 {
		buf.WriteString(" WHERE ")
		err := And(b.WhereCond...).Build(d, buf)
		if err != nil {
			return err
		}
	}
	return nil
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
