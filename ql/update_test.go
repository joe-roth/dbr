package ql

import (
	"testing"

	"github.com/gocraft/dbr/ql/dialect"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBuilder(t *testing.T) {
	builder := Update("table").Set("a", 1).Where(Eq("b", 2))
	query, value, err := builder.Build(dialect.MySQL)
	assert.NoError(t, err)

	assert.Equal(t, "UPDATE `table` SET `a` = ? WHERE (`b` = ?)", query)
	assert.Equal(t, []interface{}{1, 2}, value)
}

func BenchmarkUpdateValuesSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Update("table").Set("a", 1).Build(dialect.MySQL)
	}
}

func BenchmarkUpdateMapSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Update("table").SetMap(map[string]interface{}{"a": 1, "b": 2}).Build(dialect.MySQL)
	}
}
