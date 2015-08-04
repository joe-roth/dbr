package ql

import "bytes"

type Direction bool

// orderby directions
const (
	ASC  Direction = false
	DESC           = true
)

type order struct {
	Column string
	// most databases by default use asc
	Direction Direction
}

func (order *order) Build(d Dialect) (string, []interface{}, error) {
	buf := new(bytes.Buffer)
	buf.WriteString(d.QuoteIdent(order.Column))
	buf.WriteRune(' ')
	switch order.Direction {
	case ASC:
		buf.WriteString("ASC")
	case DESC:
		buf.WriteString("DESC")
	}
	return buf.String(), nil, nil
}
