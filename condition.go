package dbr

import "github.com/gocraft/dbr/ql"

type And map[string]interface{}

func (and And) Build(d ql.Dialect, buf ql.Buffer) error {
	var cond []ql.Condition
	for col, val := range and {
		cond = append(cond, ql.Eq(col, val))
	}
	return ql.And(cond...).Build(d, buf)
}

type Or map[string]interface{}

func (or Or) Build(d ql.Dialect, buf ql.Buffer) error {
	var cond []ql.Condition
	for col, val := range or {
		cond = append(cond, ql.Eq(col, val))
	}
	return ql.Or(cond...).Build(d, buf)
}
