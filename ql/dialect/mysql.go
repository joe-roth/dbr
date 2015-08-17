package dialect

import (
	"bytes"
	"fmt"
	"time"
)

type mysql struct{}

func (d mysql) QuoteIdent(s string) string {
	return quoteIdent(s, "`")
}

func (d mysql) EncodeString(s string) string {
	buf := new(bytes.Buffer)

	buf.WriteRune('\'')

	for _, char := range s {
		switch char {
		case '\'': // single quote: ' -> \'
			buf.WriteString("\\'")
		case '"': // double quote: " -> \"
			buf.WriteString("\\\"")
		case '\\': // slash: \ -> "\\"
			buf.WriteString("\\\\")
		case '\n': // control: newline: \n -> "\n"
			buf.WriteString("\\n")
		case '\r': // control: return: \r -> "\r"
			buf.WriteString("\\r")
		case 0: // control: NUL: 0 -> "\x00"
			buf.WriteString("\\x00")
		case 0x1a: // control: \x1a -> "\x1a"
			buf.WriteString("\\x1a")
		default:
			buf.WriteRune(char)
		}
	}

	buf.WriteRune('\'')
	return buf.String()
}

func (d mysql) EncodeBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (d mysql) EncodeTime(t time.Time) string {
	return d.EncodeString(t.UTC().Format(timeFormat))
}

func (d mysql) EncodeBytes(b []byte) string {
	return fmt.Sprintf(`0x%x`, b)
}

func (d mysql) Placeholder() string {
	return "?"
}
