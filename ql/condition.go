package ql

import "reflect"

// Condition abstracts AND, OR and simple conditions like eq.
type Condition interface {
	Builder
}

func buildCond(d Dialect, buf Buffer, pred string, cond ...Condition) error {
	for i, c := range cond {
		if i > 0 {
			buf.WriteString(" ")
			buf.WriteString(pred)
			buf.WriteString(" ")
		}
		buf.WriteString("(")
		err := c.Build(d, buf)
		if err != nil {
			return err
		}
		buf.WriteString(")")
	}
	return nil
}

type and []Condition

func (and and) Build(d Dialect, buf Buffer) error {
	return buildCond(d, buf, "AND", and...)
}

// And creates AND from a list of conditions
func And(cond ...Condition) Condition {
	return and(cond)
}

type or []Condition

func (or or) Build(d Dialect, buf Buffer) error {
	return buildCond(d, buf, "OR", or...)
}

// Or creates OR from a list of conditions
func Or(cond ...Condition) Condition {
	return or(cond)
}

type cmp struct {
	Column string
	Value  interface{}
}

func buildCmp(d Dialect, buf Buffer, pred string, column string, value interface{}) error {
	buf.WriteString(d.QuoteIdent(column))
	buf.WriteString(" ")
	buf.WriteString(pred)
	buf.WriteString(" ")
	buf.WriteString(d.Placeholder())

	buf.WriteValue(value)
	return nil
}

type eq cmp

// Eq is `=`.
// When value is nil, it will be translated to `IS NULL`.
// When value is a slice, it will be translated to `IN`.
// Otherwise it will be translated to `=`.
func Eq(column string, value interface{}) Condition {
	return &eq{
		Column: column,
		Value:  value,
	}
}

func (eq *eq) Build(d Dialect, buf Buffer) error {
	if eq.Value == nil {
		buf.WriteString(d.QuoteIdent(eq.Column))
		buf.WriteString(" IS NULL")
		return nil
	}
	if reflect.ValueOf(eq.Value).Kind() == reflect.Slice {
		return buildCmp(d, buf, "IN", eq.Column, eq.Value)
	}
	return buildCmp(d, buf, "=", eq.Column, eq.Value)
}

type neq cmp

// Neq is `!=`.
// When value is nil, it will be translated to `IS NOT NULL`.
// When value is a slice, it will be translated to `NOT IN`.
// Otherwise it will be translated to `!=`.
func Neq(column string, value interface{}) Condition {
	return &neq{
		Column: column,
		Value:  value,
	}
}

func (neq *neq) Build(d Dialect, buf Buffer) error {
	if neq.Value == nil {
		buf.WriteString(d.QuoteIdent(neq.Column))
		buf.WriteString(" IS NOT NULL")
		return nil
	}
	if reflect.ValueOf(neq.Value).Kind() == reflect.Slice {
		return buildCmp(d, buf, "NOT IN", neq.Column, neq.Value)
	}
	return buildCmp(d, buf, "!=", neq.Column, neq.Value)
}

type gt cmp

// Gt is `>`.
func Gt(column string, value interface{}) Condition {
	return &gt{
		Column: column,
		Value:  value,
	}
}

func (gt *gt) Build(d Dialect, buf Buffer) error {
	return buildCmp(d, buf, ">", gt.Column, gt.Value)
}

type gte cmp

// Gte is '>='.
func Gte(column string, value interface{}) Condition {
	return &gte{
		Column: column,
		Value:  value,
	}
}

func (gte *gte) Build(d Dialect, buf Buffer) error {
	return buildCmp(d, buf, ">=", gte.Column, gte.Value)
}

type lt cmp

// Lt is '<'.
func Lt(column string, value interface{}) Condition {
	return &lt{
		Column: column,
		Value:  value,
	}
}

func (lt *lt) Build(d Dialect, buf Buffer) error {
	return buildCmp(d, buf, "<", lt.Column, lt.Value)
}

type lte cmp

// Lte is `<=`.
func Lte(column string, value interface{}) Condition {
	return &lte{
		Column: column,
		Value:  value,
	}
}

func (lte *lte) Build(d Dialect, buf Buffer) error {
	return buildCmp(d, buf, "<=", lte.Column, lte.Value)
}
