package ql

import "bytes"

type JoinType uint8

const (
	Inner JoinType = iota
	Left
	Right
	Full
)

type join struct {
	Table interface{}
	Type  JoinType
	On    []Condition
}

func (join *join) Build(d Dialect) (string, []interface{}, error) {
	buf := new(bytes.Buffer)
	var value []interface{}
	buf.WriteRune(' ')
	switch join.Type {
	case Left:
		buf.WriteString("LEFT ")
	case Right:
		buf.WriteString("RIGHT ")
	case Full:
		buf.WriteString("FULL ")
	}
	buf.WriteString("JOIN ")
	switch table := join.Table.(type) {
	case string:
		buf.WriteString(d.QuoteIdent(table))
	default:
		buf.WriteString(d.Placeholder())
		value = append(value, table)
	}
	buf.WriteString(" ON ")
	buf.WriteString(d.Placeholder())
	value = append(value, And(join.On...))

	return buf.String(), value, nil
}

type on struct {
	Column1, Column2 string
}

func On(col1, col2 string) Condition {
	return &on{
		Column1: col1,
		Column2: col2,
	}
}

func (on *on) Build(d Dialect) (string, []interface{}, error) {
	buf := new(bytes.Buffer)
	buf.WriteString(d.QuoteIdent(on.Column1))
	buf.WriteString(" = ")
	buf.WriteString(d.QuoteIdent(on.Column2))
	return buf.String(), nil, nil
}
