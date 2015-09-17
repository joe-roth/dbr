package dbr

import (
	"testing"

	"github.com/gocraft/dbr/dialect"
	"github.com/stretchr/testify/assert"
)

func TestSelectBuilder(t *testing.T) {
	buf := NewBuffer()
	sub, on := Select("a").From("table"), On("table.a1", "table2.a2")
	builder := Select("a", "b").From(sub).Join(Left, "table2", on).Distinct().Where(Eq("c", 1)).GroupBy("d").Having(Eq("e", 2)).OrderBy("f", ASC).Limit(3).Offset(4)
	err := builder.Build(dialect.MySQL, buf)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT DISTINCT a, b FROM ? LEFT JOIN `table2` ON ? WHERE (`c` = ?) GROUP BY d HAVING (`e` = ?) ORDER BY f ASC LIMIT 3 OFFSET 4", buf.String())
	// two functions cannot be compared
	assert.Equal(t, 4, len(buf.Value()))
}

func BenchmarkSelectSQL(b *testing.B) {
	buf := NewBuffer()
	for i := 0; i < b.N; i++ {
		Select("a", "b").From("table").Where(Eq("c", 1)).OrderBy("d", ASC).Build(dialect.MySQL, buf)
	}
}
