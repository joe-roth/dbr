package ql

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

func Join(t JoinType, table interface{}, on ...Condition) Builder {
	return &join{
		Table: table,
		Type:  t,
		On:    on,
	}
}

func (join *join) Build(d Dialect, buf Buffer) error {
	buf.WriteString(" ")
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
		buf.WriteValue(table)
	}
	buf.WriteString(" ON ")
	buf.WriteString(d.Placeholder())
	buf.WriteValue(And(join.On...))

	return nil
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

func (on *on) Build(d Dialect, buf Buffer) error {
	buf.WriteString(d.QuoteIdent(on.Column1))
	buf.WriteString(" = ")
	buf.WriteString(d.QuoteIdent(on.Column2))
	return nil
}
