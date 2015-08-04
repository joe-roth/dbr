package ql

// Builder builds sql in one dialect like MySQL/PostgreSQL
// e.g. XxxBuilder, Condition
type Builder interface {
	Build(Dialect) (string, []interface{}, error)
}
