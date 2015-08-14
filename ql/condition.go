package ql

import (
	"bytes"
	"reflect"
)

// Condition abstracts AND, OR and simple conditions like eq.
type Condition interface {
	Builder
}

func buildCond(d Dialect, pred string, cond ...Condition) (string, []interface{}, error) {
	buf := new(bytes.Buffer)
	var value []interface{}
	for i, c := range cond {
		if i > 0 {
			buf.WriteRune(' ')
			buf.WriteString(pred)
			buf.WriteRune(' ')
		}
		query, v, err := c.Build(d)
		if err != nil {
			return "", nil, err
		}
		buf.WriteRune('(')
		buf.WriteString(query)
		buf.WriteRune(')')
		value = append(value, v...)
	}
	return buf.String(), value, nil
}

type and []Condition

func (and and) Build(d Dialect) (string, []interface{}, error) {
	return buildCond(d, "AND", and...)
}

// And creates AND from a list of conditions
func And(cond ...Condition) Condition {
	return and(cond)
}

type or []Condition

func (or or) Build(d Dialect) (string, []interface{}, error) {
	return buildCond(d, "OR", or...)
}

// Or creates OR from a list of conditions
func Or(cond ...Condition) Condition {
	return or(cond)
}

type cmp struct {
	Column string
	Value  interface{}
}

func buildCmp(d Dialect, pred string, column string, value interface{}) (string, []interface{}, error) {
	buf := new(bytes.Buffer)
	buf.WriteString(d.QuoteIdent(column))
	buf.WriteRune(' ')
	buf.WriteString(pred)
	buf.WriteRune(' ')
	buf.WriteString(d.Placeholder())
	return buf.String(), []interface{}{value}, nil
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

func (eq *eq) Build(d Dialect) (string, []interface{}, error) {
	if eq.Value == nil {
		return d.QuoteIdent(eq.Column) + " IS NULL", nil, nil
	}
	if reflect.ValueOf(eq.Value).Kind() == reflect.Slice {
		return buildCmp(d, "IN", eq.Column, eq.Value)
	}
	return buildCmp(d, "=", eq.Column, eq.Value)
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

func (neq *neq) Build(d Dialect) (string, []interface{}, error) {
	if neq.Value == nil {
		return d.QuoteIdent(neq.Column) + " IS NOT NULL", nil, nil
	}
	if reflect.ValueOf(neq.Value).Kind() == reflect.Slice {
		return buildCmp(d, "NOT IN", neq.Column, neq.Value)
	}
	return buildCmp(d, "!=", neq.Column, neq.Value)
}

type gt cmp

// Gt is `>`.
func Gt(column string, value interface{}) Condition {
	return &gt{
		Column: column,
		Value:  value,
	}
}

func (gt *gt) Build(d Dialect) (string, []interface{}, error) {
	return buildCmp(d, ">", gt.Column, gt.Value)
}

type gte cmp

// Gte is '>='.
func Gte(column string, value interface{}) Condition {
	return &gte{
		Column: column,
		Value:  value,
	}
}

func (gte *gte) Build(d Dialect) (string, []interface{}, error) {
	return buildCmp(d, ">=", gte.Column, gte.Value)
}

type lt cmp

// Lt is '<'.
func Lt(column string, value interface{}) Condition {
	return &lt{
		Column: column,
		Value:  value,
	}
}

func (lt *lt) Build(d Dialect) (string, []interface{}, error) {
	return buildCmp(d, "<", lt.Column, lt.Value)
}

type lte cmp

// Lte is `<=`.
func Lte(column string, value interface{}) Condition {
	return &lte{
		Column: column,
		Value:  value,
	}
}

func (lte *lte) Build(d Dialect) (string, []interface{}, error) {
	return buildCmp(d, "<=", lte.Column, lte.Value)
}
