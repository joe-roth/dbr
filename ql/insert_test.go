package ql

import (
	"testing"

	"github.com/gocraft/dbr/ql/dialect"
	"github.com/stretchr/testify/assert"
)

type insertTest struct {
	A int
	C string `db:"b"`
}

func TestInsertBuilder(t *testing.T) {
	builder := InsertInto("table").Columns("a", "b").Values(1, "one").Record(&insertTest{
		A: 2,
		C: "two",
	})
	query, value, err := builder.Build(dialect.MySQL)
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO `table` (`a`,`b`) VALUES (?,?), (?,?)", query)
	assert.Equal(t, []interface{}{1, "one", 2, "two"}, value)
}

func BenchmarkInsertValuesSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InsertInto("table").Columns("a", "b").Values(1, "one").Build(dialect.MySQL)
	}
}

func BenchmarkInsertRecordSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InsertInto("table").Columns("a", "b").Record(&insertTest{
			A: 2,
			C: "two",
		}).Build(dialect.MySQL)
	}
}
