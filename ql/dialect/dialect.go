package dialect

import "strings"

var (
	// MySQL dialect
	MySQL = mysql{}
	// PostgreSQL dialect
	PostgreSQL = postgreSQL{}
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func isSQLFunc(s string) bool {
	return strings.ContainsAny(s, "()* ")
}
