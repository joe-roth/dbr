package dialect

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
	for _, c := range s {
		switch {
		case c == '.':
		case 'a' <= c && c <= 'z':
		case 'A' <= c && c <= 'Z':
		case '0' <= c && c <= '9':
		default:
			return true
		}
	}
	return false
}
