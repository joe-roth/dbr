package ql

// As creates an alias for expr. e.g. SELECT `a1` AS `a2`
func As(expr interface{}, alias string) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		switch expr := expr.(type) {
		case string:
			buf.WriteString(d.QuoteIdent(expr))
		default:
			buf.WriteString(d.Placeholder())
			buf.WriteValue(expr)
		}
		buf.WriteString(" AS ")
		buf.WriteString(d.QuoteIdent(alias))
		return nil
	})
}
