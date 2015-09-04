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

// And creates AND from a list of conditions
func And(cond ...Condition) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		return buildCond(d, buf, "AND", cond...)
	})
}

// Or creates OR from a list of conditions
func Or(cond ...Condition) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		return buildCond(d, buf, "OR", cond...)
	})
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

// Eq is `=`.
// When value is nil, it will be translated to `IS NULL`.
// When value is a slice, it will be translated to `IN`.
// Otherwise it will be translated to `=`.
func Eq(column string, value interface{}) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		if value == nil {
			buf.WriteString(d.QuoteIdent(column))
			buf.WriteString(" IS NULL")
			return nil
		}
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			return buildCmp(d, buf, "IN", column, value)
		}
		return buildCmp(d, buf, "=", column, value)
	})
}

// Neq is `!=`.
// When value is nil, it will be translated to `IS NOT NULL`.
// When value is a slice, it will be translated to `NOT IN`.
// Otherwise it will be translated to `!=`.
func Neq(column string, value interface{}) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		if value == nil {
			buf.WriteString(d.QuoteIdent(column))
			buf.WriteString(" IS NOT NULL")
			return nil
		}
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			return buildCmp(d, buf, "NOT IN", column, value)
		}
		return buildCmp(d, buf, "!=", column, value)
	})
}

// Gt is `>`.
func Gt(column string, value interface{}) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		return buildCmp(d, buf, ">", column, value)
	})
}

// Gte is '>='.
func Gte(column string, value interface{}) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		return buildCmp(d, buf, ">=", column, value)
	})
}

// Lt is '<'.
func Lt(column string, value interface{}) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		return buildCmp(d, buf, "<", column, value)
	})
}

// Lte is `<=`.
func Lte(column string, value interface{}) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		return buildCmp(d, buf, "<=", column, value)
	})
}
