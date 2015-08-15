package ql

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

func Order(column string, dir Direction) Builder {
	return &order{
		Column:    column,
		Direction: dir,
	}
}

func (order *order) Build(d Dialect, buf Buffer) error {
	buf.WriteString(d.QuoteIdent(order.Column))
	buf.WriteString(" ")
	switch order.Direction {
	case ASC:
		buf.WriteString("ASC")
	case DESC:
		buf.WriteString("DESC")
	}
	return nil
}
