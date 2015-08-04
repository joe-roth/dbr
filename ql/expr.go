package ql

// XxxBuilders all support raw query
type raw struct {
	Query string
	Value []interface{}
}

// Expr should be used when sql syntax is not supported
func Expr(query string, value ...interface{}) Builder {
	return &raw{Query: query, Value: value}
}

func (raw *raw) Build(_ Dialect) (string, []interface{}, error) {
	return raw.Query, raw.Value, nil
}
