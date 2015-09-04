package ql

type Direction bool

// orderby directions
// most databases by default use asc
const (
	ASC  Direction = false
	DESC           = true
)

func Order(column string, dir Direction) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString(d.QuoteIdent(column))
		buf.WriteString(" ")
		switch dir {
		case ASC:
			buf.WriteString("ASC")
		case DESC:
			buf.WriteString("DESC")
		}
		return nil
	})
}
