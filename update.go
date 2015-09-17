package dbr

// UpdateBuilder builds `UPDATE ...`
type UpdateBuilder struct {
	raw

	Table string
	Value map[string]interface{}

	WhereCond []Condition
}

// Build builds `UPDATE ...` in dialect
func (b *UpdateBuilder) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	if len(b.Value) == 0 {
		return ErrColumnNotSpecified
	}

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

		buf.WriteValue(v)
		i++
	}

	if len(b.WhereCond) > 0 {
		buf.WriteString(" WHERE ")
		err := And(b.WhereCond...).Build(d, buf)
		if err != nil {
			return err
		}
	}
	return nil
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
