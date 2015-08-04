package dialect

import (
	"fmt"
	"strings"
	"time"
)

type postgreSQL struct{}

func (d postgreSQL) QuoteIdent(column string) string {
	if isSQLFunc(column) {
		return column
	}
	part := strings.SplitN(column, ".", 2)
	if len(part) == 2 {
		return d.QuoteIdent(part[0]) + "." + d.QuoteIdent(part[1])
	}
	return `"` + column + `"`
}

func (d postgreSQL) EncodeString(s string) string {
	return MySQL.EncodeString(s)
}

func (d postgreSQL) EncodeBool(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
}

func (d postgreSQL) EncodeTime(t time.Time) string {
	return MySQL.EncodeTime(t)
}

func (d postgreSQL) EncodeBytes(b []byte) string {
	return d.EncodeString(`\x` + fmt.Sprintf("%x", b))
}

func (d postgreSQL) Placeholder() string {
	return "?"
}
