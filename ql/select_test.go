package ql

import (
	"testing"

	"github.com/gocraft/dbr/ql/dialect"
	"github.com/stretchr/testify/assert"
)

func TestSelectBuilder(t *testing.T) {
	sub := Select("a").From("table")
	builder := Select("a", "b").From(sub).Join(Left, "table2", On("table.a1", "table2.a2")).Distinct().Where(Eq("c", 1)).GroupBy("d").Having(Eq("e", 2)).OrderBy("f", ASC).Limit(3).Offset(4)
	query, value, err := builder.Build(dialect.MySQL)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT DISTINCT `a`, `b` FROM ? LEFT JOIN `table2` ON ? WHERE (`c` = ?) GROUP BY `d` HAVING (`e` = ?) ORDER BY `f` ASC LIMIT 3 OFFSET 4", query)
	assert.Equal(t, []interface{}{sub, 1, 2}, value)
}

func BenchmarkSelectSQL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Select("a", "b").From("table").Where(Eq("c", 1)).OrderBy("d", ASC).Build(dialect.MySQL)
	}
}
