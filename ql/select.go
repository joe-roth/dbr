package ql

import (
	"bytes"
	"fmt"
)

// SelectBuilder builds `SELECT ...`
type SelectBuilder struct {
	raw

	IsDistinct bool

	Column    []interface{}
	Table     interface{}
	JoinTable []Builder

	WhereCond  []Condition
	Group      []string
	HavingCond []Condition
	Order      []Builder

	LimitCount  int64
	OffsetCount int64
}

// Build builds `SELECT ...` in dialect
func (b *SelectBuilder) Build(d Dialect) (string, []interface{}, error) {
	if b.raw.Query != "" {
		return b.raw.Query, b.raw.Value, nil
	}

	if len(b.Column) == 0 {
		return "", nil, ErrColumnNotSpecified
	}

	buf := new(bytes.Buffer)
	var value []interface{}

	buf.WriteString("SELECT ")

	if b.IsDistinct {
		buf.WriteString("DISTINCT ")
	}

	for i, col := range b.Column {
		if i > 0 {
			buf.WriteString(", ")
		}
		switch col := col.(type) {
		case string:
			buf.WriteString(d.QuoteIdent(col))
		default:
			buf.WriteString(d.Placeholder())
			value = append(value, col)
		}
	}

	if b.Table != nil {
		buf.WriteString(" FROM ")
		switch table := b.Table.(type) {
		case string:
			buf.WriteString(d.QuoteIdent(table))
		default:
			buf.WriteString(d.Placeholder())
			value = append(value, table)
		}
		if len(b.JoinTable) > 0 {
			for _, join := range b.JoinTable {
				query, _, err := join.Build(d)
				if err != nil {
					return "", nil, err
				}
				buf.WriteString(query)
			}
		}
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

	if len(b.Group) > 0 {
		buf.WriteString(" GROUP BY ")
		for i, col := range b.Group {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(d.QuoteIdent(col))
		}
	}

	if len(b.HavingCond) > 0 {
		query, v, err := And(b.HavingCond...).Build(d)
		if err != nil {
			return "", nil, err
		}
		if query != "" {
			buf.WriteString(" HAVING ")
			buf.WriteString(query)

			value = append(value, v...)
		}
	}

	if len(b.Order) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, order := range b.Order {
			if i > 0 {
				buf.WriteString(", ")
			}
			query, _, err := order.Build(d)
			if err != nil {
				return "", nil, err
			}
			buf.WriteString(query)
		}
	}

	if b.LimitCount >= 0 {
		buf.WriteString(" LIMIT ")
		fmt.Fprint(buf, b.LimitCount)
	}

	if b.OffsetCount >= 0 {
		buf.WriteString(" OFFSET ")
		fmt.Fprint(buf, b.OffsetCount)
	}
	return buf.String(), value, nil
}

// Select creates a SelectBuilder
func Select(column ...interface{}) *SelectBuilder {
	return &SelectBuilder{
		Column:      column,
		LimitCount:  -1,
		OffsetCount: -1,
	}
}

// From specifies table
func (b *SelectBuilder) From(table interface{}) *SelectBuilder {
	b.Table = table
	return b
}

// SelectBySQL creates a SelectBuilder from raw query
func SelectBySQL(query string, value ...interface{}) *SelectBuilder {
	return &SelectBuilder{
		raw: raw{
			Query: query,
			Value: value,
		},
		LimitCount:  -1,
		OffsetCount: -1,
	}
}

// Distinct adds `DISTINCT`
func (b *SelectBuilder) Distinct() *SelectBuilder {
	b.IsDistinct = true
	return b
}

// Where adds a where condition
func (b *SelectBuilder) Where(query interface{}, value ...interface{}) *SelectBuilder {
	switch query := query.(type) {
	case string:
		b.WhereCond = append(b.WhereCond, Expr(query, value...))
	case Condition:
		b.WhereCond = append(b.WhereCond, query)
	}
	return b
}

// Having adds a having condition
func (b *SelectBuilder) Having(query interface{}, value ...interface{}) *SelectBuilder {
	switch query := query.(type) {
	case string:
		b.HavingCond = append(b.HavingCond, Expr(query, value...))
	case Condition:
		b.HavingCond = append(b.HavingCond, query)
	}
	return b
}

// GroupBy specifies columns for grouping
func (b *SelectBuilder) GroupBy(col ...string) *SelectBuilder {
	b.Group = col
	return b
}

// OrderBy specifies columns for ordering
func (b *SelectBuilder) OrderBy(col string, dir Direction) *SelectBuilder {
	b.Order = append(b.Order, Order(col, dir))
	return b
}

// Limit adds limit
func (b *SelectBuilder) Limit(n uint64) *SelectBuilder {
	b.LimitCount = int64(n)
	return b
}

// Offset adds offset
func (b *SelectBuilder) Offset(n uint64) *SelectBuilder {
	b.OffsetCount = int64(n)
	return b
}

// Join joins table on condition
func (b *SelectBuilder) Join(t JoinType, table interface{}, cond ...Condition) *SelectBuilder {
	b.JoinTable = append(b.JoinTable, Join(t, table, cond...))
	return b
}
