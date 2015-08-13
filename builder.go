package dbr

import (
	"database/sql"

	"github.com/gocraft/dbr/ql"
)

type SelectBuilder struct {
	runner
	EventReceiver
	Dialect ql.Dialect

	*ql.SelectBuilder
}

func (sess *Session) Select(column ...interface{}) *SelectBuilder {
	return &SelectBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		SelectBuilder: ql.Select(column...),
	}
}

func (tx *Tx) Select(column ...interface{}) *SelectBuilder {
	return &SelectBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		SelectBuilder: ql.Select(column...),
	}
}

func (sess *Session) SelectBySql(query string, value ...interface{}) *SelectBuilder {
	return &SelectBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		SelectBuilder: ql.SelectBySQL(query, value...),
	}
}

func (tx *Tx) SelectBySql(query string, value ...interface{}) *SelectBuilder {
	return &SelectBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		SelectBuilder: ql.SelectBySQL(query, value...),
	}
}

func (b *SelectBuilder) ToSql() (string, []interface{}) {
	s, v, err := b.Build(b.Dialect)
	panic(err)
	return s, v
}

func (b *SelectBuilder) Load(value interface{}) (int, error) {
	return query(b.runner, b.EventReceiver, b, b.Dialect, value)
}

func (b *SelectBuilder) LoadStruct(value interface{}) error {
	_, err := query(b.runner, b.EventReceiver, b, b.Dialect, value)
	if err != nil {
		return err
	}
	return nil
}

func (b *SelectBuilder) LoadStructs(value interface{}) (int, error) {
	return query(b.runner, b.EventReceiver, b, b.Dialect, value)
}

func (b *SelectBuilder) LoadValue(value interface{}) error {
	_, err := query(b.runner, b.EventReceiver, b, b.Dialect, value)
	if err != nil {
		return err
	}
	return nil
}

func (b *SelectBuilder) LoadValues(value interface{}) (int, error) {
	return query(b.runner, b.EventReceiver, b, b.Dialect, value)
}

func (b *SelectBuilder) Join(t ql.JoinType, table interface{}, cond ...ql.Condition) *SelectBuilder {
	b.Join(t, table, cond...)
	return b
}

func (b *SelectBuilder) Distinct() *SelectBuilder {
	b.SelectBuilder.Distinct()
	return b
}

func (b *SelectBuilder) From(table interface{}) *SelectBuilder {
	b.SelectBuilder.From(table)
	return b
}

func (b *SelectBuilder) GroupBy(col ...string) *SelectBuilder {
	b.SelectBuilder.GroupBy(col...)
	return b
}

func (b *SelectBuilder) Having(query interface{}, value ...interface{}) *SelectBuilder {
	b.SelectBuilder.Having(query, value...)
	return b
}

func (b *SelectBuilder) Limit(n uint64) *SelectBuilder {
	b.SelectBuilder.Limit(n)
	return b
}

func (b *SelectBuilder) Offset(n uint64) *SelectBuilder {
	b.SelectBuilder.Offset(n)
	return b
}

func (b *SelectBuilder) OrderDir(col string, isAsc bool) *SelectBuilder {
	if isAsc {
		b.SelectBuilder.OrderBy(col, ql.ASC)
	} else {
		b.SelectBuilder.OrderBy(col, ql.DESC)
	}
	return b
}

func (b *SelectBuilder) Paginate(page, perPage uint64) *SelectBuilder {
	b.Limit(perPage)
	b.Offset((page - 1) * perPage)
	return b
}

func (b *SelectBuilder) OrderBy(col string) *SelectBuilder {
	b.SelectBuilder.OrderBy(col, ql.ASC)
	return b
}

func (b *SelectBuilder) Where(query interface{}, value ...interface{}) *SelectBuilder {
	b.SelectBuilder.Where(query, value...)
	return b
}

type InsertBuilder struct {
	runner
	EventReceiver
	Dialect ql.Dialect

	*ql.InsertBuilder
}

func (sess *Session) InsertInto(table string) *InsertBuilder {
	return &InsertBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		InsertBuilder: ql.InsertInto(table),
	}
}

func (tx *Tx) InsertInto(table string) *InsertBuilder {
	return &InsertBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		InsertBuilder: ql.InsertInto(table),
	}
}

func (sess *Session) InsertBySql(query string, value ...interface{}) *InsertBuilder {
	return &InsertBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		InsertBuilder: ql.InsertBySQL(query, value...),
	}
}

func (tx *Tx) InsertBySql(query string, value ...interface{}) *InsertBuilder {
	return &InsertBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		InsertBuilder: ql.InsertBySQL(query, value...),
	}
}

func (b *InsertBuilder) ToSql() (string, []interface{}) {
	s, v, err := b.Build(b.Dialect)
	panic(err)
	return s, v
}

func (b *InsertBuilder) Exec() (sql.Result, error) {
	return exec(b.runner, b.EventReceiver, b, b.Dialect)
}

func (b *InsertBuilder) Columns(column ...string) *InsertBuilder {
	b.InsertBuilder.Columns(column...)
	return b
}

func (b *InsertBuilder) Record(structValue interface{}) *InsertBuilder {
	b.InsertBuilder.Record(structValue)
	return b
}

func (b *InsertBuilder) Values(value ...interface{}) *InsertBuilder {
	b.InsertBuilder.Values(value...)
	return b
}

type UpdateBuilder struct {
	runner
	EventReceiver
	Dialect ql.Dialect

	*ql.UpdateBuilder
}

func (sess *Session) Update(table string) *UpdateBuilder {
	return &UpdateBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		UpdateBuilder: ql.Update(table),
	}
}

func (tx *Tx) Update(table string) *UpdateBuilder {
	return &UpdateBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		UpdateBuilder: ql.Update(table),
	}
}

func (sess *Session) UpdateBySql(query string, value ...interface{}) *UpdateBuilder {
	return &UpdateBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		UpdateBuilder: ql.UpdateBySQL(query, value...),
	}
}

func (tx *Tx) UpdateBySql(query string, value ...interface{}) *UpdateBuilder {
	return &UpdateBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		UpdateBuilder: ql.UpdateBySQL(query, value...),
	}
}

func (b *UpdateBuilder) ToSql() (string, []interface{}) {
	s, v, err := b.Build(b.Dialect)
	panic(err)
	return s, v
}

func (b *UpdateBuilder) Exec() (sql.Result, error) {
	return exec(b.runner, b.EventReceiver, b, b.Dialect)
}

func (b *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	b.UpdateBuilder.Set(column, value)
	return b
}

func (b *UpdateBuilder) SetMap(m map[string]interface{}) *UpdateBuilder {
	b.UpdateBuilder.SetMap(m)
	return b
}

func (b *UpdateBuilder) Where(query interface{}, value ...interface{}) *UpdateBuilder {
	b.UpdateBuilder.Where(query, value...)
	return b
}

type DeleteBuilder struct {
	runner
	EventReceiver
	Dialect ql.Dialect

	*ql.DeleteBuilder
}

func (sess *Session) DeleteFrom(table string) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		DeleteBuilder: ql.DeleteFrom(table),
	}
}

func (tx *Tx) DeleteFrom(table string) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		DeleteBuilder: ql.DeleteFrom(table),
	}
}

func (sess *Session) DeleteBySql(query string, value ...interface{}) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		DeleteBuilder: ql.DeleteBySQL(query, value...),
	}
}

func (tx *Tx) DeleteBySql(query string, value ...interface{}) *DeleteBuilder {
	return &DeleteBuilder{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		DeleteBuilder: ql.DeleteBySQL(query, value...),
	}
}

func (b *DeleteBuilder) ToSql() (string, []interface{}) {
	s, v, err := b.Build(b.Dialect)
	panic(err)
	return s, v
}

func (b *DeleteBuilder) Exec() (sql.Result, error) {
	return exec(b.runner, b.EventReceiver, b, b.Dialect)
}

func (b *DeleteBuilder) Where(query interface{}, value ...interface{}) *DeleteBuilder {
	b.DeleteBuilder.Where(query, value...)
	return b
}
