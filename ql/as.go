package ql

import "bytes"

type as struct {
	Expr  interface{}
	Alias string
}

// As creates an alias for expr. e.g. SELECT `a1` AS `a2`
func As(expr interface{}, alias string) Builder {
	return &as{
		Expr:  expr,
		Alias: alias,
	}
}

func (as *as) Build(d Dialect) (string, []interface{}, error) {
	buf := new(bytes.Buffer)
	var value []interface{}
	switch expr := as.Expr.(type) {
	case string:
		buf.WriteString(d.QuoteIdent(expr))
	default:
		buf.WriteString(d.Placeholder())
		value = append(value, expr)
	}
	buf.WriteString(" AS ")
	buf.WriteString(d.QuoteIdent(as.Alias))
	return buf.String(), value, nil
}
