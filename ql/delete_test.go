package ql

import (
	"testing"

	"github.com/gocraft/dbr/ql/dialect"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBuilder(t *testing.T) {
	builder := DeleteFrom("table").Where(Eq("a", 1))
	query, value, err := builder.Build(dialect.MySQL)
	assert.NoError(t, err)
	assert.Equal(t, "DELETE FROM `table` WHERE (`a` = ?)", query)
	assert.Equal(t, []interface{}{1}, value)
}

func BenchmarkDeleteSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DeleteFrom("table").Where(Eq("a", 1)).Build(dialect.MySQL)
	}
}
