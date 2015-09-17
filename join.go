package dbr

type JoinType uint8

const (
	Inner JoinType = iota
	Left
	Right
	Full
)

func Join(t JoinType, table interface{}, on ...Condition) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString(" ")
		switch t {
		case Left:
			buf.WriteString("LEFT ")
		case Right:
			buf.WriteString("RIGHT ")
		case Full:
			buf.WriteString("FULL ")
		}
		buf.WriteString("JOIN ")
		switch table := table.(type) {
		case string:
			buf.WriteString(d.QuoteIdent(table))
		default:
			buf.WriteString(d.Placeholder())
			buf.WriteValue(table)
		}
		buf.WriteString(" ON ")
		buf.WriteString(d.Placeholder())
		buf.WriteValue(And(on...))

		return nil
	})
}

func On(col1, col2 string) Condition {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString(d.QuoteIdent(col1))
		buf.WriteString(" = ")
		buf.WriteString(d.QuoteIdent(col2))
		return nil
	})
}
