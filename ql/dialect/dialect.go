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

func isExpr(s string) bool {
	for _, c := range s {
		switch {
		case c == '.':
		case c == '_':
		case 'a' <= c && c <= 'z':
		case 'A' <= c && c <= 'Z':
		case '0' <= c && c <= '9':
		default:
			return true
		}
	}
	return false
}

func quoteIdent(s, quote string) string {
	part := strings.SplitN(s, ".", 2)
	if len(part) == 2 {
		return quoteIdent(part[0], quote) + "." + quoteIdent(part[1], quote)
	}
	return quote + s + quote
}
