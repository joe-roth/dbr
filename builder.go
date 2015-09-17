package dbr

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Builder builds sql in one dialect like MySQL/PostgreSQL
// e.g. XxxBuilder, Condition
type Builder interface {
	Build(Dialect, Buffer) error
}

type BuildFunc func(Dialect, Buffer) error

func (b BuildFunc) Build(d Dialect, buf Buffer) error {
	return b(d, buf)
}

type SelectBuilderSession struct {
	runner
	EventReceiver
	Dialect Dialect

	*SelectBuilder
}

func prepareSelect(a []string) []interface{} {
	b := make([]interface{}, len(a))
	for i := range a {
		b[i] = a[i]
	}
	return b
}

func (sess *Session) Select(column ...string) *SelectBuilderSession {
	return &SelectBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		SelectBuilder: Select(prepareSelect(column)...),
	}
}

func (tx *Tx) Select(column ...string) *SelectBuilderSession {
	return &SelectBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		SelectBuilder: Select(prepareSelect(column)...),
	}
}

func (sess *Session) SelectBySql(query string, value ...interface{}) *SelectBuilderSession {
	return &SelectBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		SelectBuilder: SelectBySQL(query, value...),
	}
}

func (tx *Tx) SelectBySql(query string, value ...interface{}) *SelectBuilderSession {
	return &SelectBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		SelectBuilder: SelectBySQL(query, value...),
	}
}

func (b *SelectBuilderSession) ToSql() (string, []interface{}) {
	buf := NewBuffer()
	err := b.Build(b.Dialect, buf)
	if err != nil {
		panic(err)
	}
	return buf.String(), buf.Value()
}

func (b *SelectBuilderSession) Load(value interface{}) (int, error) {
	return query(b.runner, b.EventReceiver, b, b.Dialect, value)
}

func (b *SelectBuilderSession) LoadStruct(value interface{}) error {
	count, err := query(b.runner, b.EventReceiver, b, b.Dialect, value)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (b *SelectBuilderSession) LoadStructs(value interface{}) (int, error) {
	return query(b.runner, b.EventReceiver, b, b.Dialect, value)
}

func (b *SelectBuilderSession) LoadValue(value interface{}) error {
	count, err := query(b.runner, b.EventReceiver, b, b.Dialect, value)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (b *SelectBuilderSession) LoadValues(value interface{}) (int, error) {
	return query(b.runner, b.EventReceiver, b, b.Dialect, value)
}

func (b *SelectBuilderSession) Join(t JoinType, table interface{}, cond ...Condition) *SelectBuilderSession {
	b.Join(t, table, cond...)
	return b
}

func (b *SelectBuilderSession) Distinct() *SelectBuilderSession {
	b.SelectBuilder.Distinct()
	return b
}

func (b *SelectBuilderSession) From(table interface{}) *SelectBuilderSession {
	b.SelectBuilder.From(table)
	return b
}

func (b *SelectBuilderSession) GroupBy(col ...string) *SelectBuilderSession {
	b.SelectBuilder.GroupBy(col...)
	return b
}

func (b *SelectBuilderSession) Having(query interface{}, value ...interface{}) *SelectBuilderSession {
	b.SelectBuilder.Having(query, value...)
	return b
}

func (b *SelectBuilderSession) Limit(n uint64) *SelectBuilderSession {
	b.SelectBuilder.Limit(n)
	return b
}

func (b *SelectBuilderSession) Offset(n uint64) *SelectBuilderSession {
	b.SelectBuilder.Offset(n)
	return b
}

func (b *SelectBuilderSession) OrderDir(col string, isAsc bool) *SelectBuilderSession {
	if isAsc {
		b.SelectBuilder.OrderBy(col, ASC)
	} else {
		b.SelectBuilder.OrderBy(col, DESC)
	}
	return b
}

func (b *SelectBuilderSession) Paginate(page, perPage uint64) *SelectBuilderSession {
	b.Limit(perPage)
	b.Offset((page - 1) * perPage)
	return b
}

func (b *SelectBuilderSession) OrderBy(col string) *SelectBuilderSession {
	b.SelectBuilder.Order = append(b.SelectBuilder.Order, Expr(col))
	return b
}

func (b *SelectBuilderSession) Where(query interface{}, value ...interface{}) *SelectBuilderSession {
	b.SelectBuilder.Where(query, value...)
	return b
}

type InsertBuilderSession struct {
	runner
	EventReceiver
	Dialect Dialect

	RecordID reflect.Value

	*InsertBuilder
}

func (sess *Session) InsertInto(table string) *InsertBuilderSession {
	return &InsertBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		InsertBuilder: InsertInto(table),
	}
}

func (tx *Tx) InsertInto(table string) *InsertBuilderSession {
	return &InsertBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		InsertBuilder: InsertInto(table),
	}
}

func (sess *Session) InsertBySql(query string, value ...interface{}) *InsertBuilderSession {
	return &InsertBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		InsertBuilder: InsertBySQL(query, value...),
	}
}

func (tx *Tx) InsertBySql(query string, value ...interface{}) *InsertBuilderSession {
	return &InsertBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		InsertBuilder: InsertBySQL(query, value...),
	}
}

func (b *InsertBuilderSession) ToSql() (string, []interface{}) {
	buf := NewBuffer()
	err := b.Build(b.Dialect, buf)
	if err != nil {
		panic(err)
	}
	return buf.String(), buf.Value()
}

func (b *InsertBuilderSession) Pair(column string, value interface{}) *InsertBuilderSession {
	b.Column = append(b.Column, column)
	switch len(b.Value) {
	case 0:
		b.InsertBuilder.Values(value)
	case 1:
		b.Value[0] = append(b.Value[0], value)
	default:
		panic("pair only allows one record to insert")
	}
	return b
}

func (b *InsertBuilderSession) Exec() (sql.Result, error) {
	result, err := exec(b.runner, b.EventReceiver, b, b.Dialect)
	if err != nil {
		return nil, err
	}

	if b.RecordID.IsValid() {
		if id, err := result.LastInsertId(); err == nil {
			b.RecordID.SetInt(id)
		}
	}

	return result, nil
}

func (b *InsertBuilderSession) Columns(column ...string) *InsertBuilderSession {
	b.InsertBuilder.Columns(column...)
	return b
}

func (b *InsertBuilderSession) Record(structValue interface{}) *InsertBuilderSession {
	v := reflect.Indirect(reflect.ValueOf(structValue))
	if v.Kind() == reflect.Struct && v.CanSet() {
		// ID is recommended by golint here
		for _, name := range []string{"Id", "ID"} {
			field := v.FieldByName(name)
			if field.IsValid() && field.Kind() == reflect.Int64 {
				b.RecordID = field
				break
			}
		}
	}

	b.InsertBuilder.Record(structValue)
	return b
}

func (b *InsertBuilderSession) Values(value ...interface{}) *InsertBuilderSession {
	b.InsertBuilder.Values(value...)
	return b
}

type UpdateBuilderSession struct {
	runner
	EventReceiver
	Dialect Dialect

	*UpdateBuilder

	LimitCount int64
}

func (sess *Session) Update(table string) *UpdateBuilderSession {
	return &UpdateBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		UpdateBuilder: Update(table),
		LimitCount:    -1,
	}
}

func (tx *Tx) Update(table string) *UpdateBuilderSession {
	return &UpdateBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		UpdateBuilder: Update(table),
		LimitCount:    -1,
	}
}

func (sess *Session) UpdateBySql(query string, value ...interface{}) *UpdateBuilderSession {
	return &UpdateBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		UpdateBuilder: UpdateBySQL(query, value...),
		LimitCount:    -1,
	}
}

func (tx *Tx) UpdateBySql(query string, value ...interface{}) *UpdateBuilderSession {
	return &UpdateBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		UpdateBuilder: UpdateBySQL(query, value...),
		LimitCount:    -1,
	}
}

func (b *UpdateBuilderSession) ToSql() (string, []interface{}) {
	buf := NewBuffer()
	err := b.Build(b.Dialect, buf)
	if err != nil {
		panic(err)
	}
	return buf.String(), buf.Value()
}

func (b *UpdateBuilderSession) Exec() (sql.Result, error) {
	return exec(b.runner, b.EventReceiver, b, b.Dialect)
}

func (b *UpdateBuilderSession) Set(column string, value interface{}) *UpdateBuilderSession {
	b.UpdateBuilder.Set(column, value)
	return b
}

func (b *UpdateBuilderSession) SetMap(m map[string]interface{}) *UpdateBuilderSession {
	b.UpdateBuilder.SetMap(m)
	return b
}

func (b *UpdateBuilderSession) Where(query interface{}, value ...interface{}) *UpdateBuilderSession {
	b.UpdateBuilder.Where(query, value...)
	return b
}

func (b *UpdateBuilderSession) Limit(n uint64) *UpdateBuilderSession {
	b.LimitCount = int64(n)
	return b
}

func (b *UpdateBuilderSession) Build(d Dialect, buf Buffer) error {
	err := b.UpdateBuilder.Build(b.Dialect, buf)
	if err != nil {
		return err
	}
	if b.LimitCount >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(fmt.Sprint(b.LimitCount))
	}
	return nil
}

type DeleteBuilderSession struct {
	runner
	EventReceiver
	Dialect Dialect

	*DeleteBuilder

	LimitCount int64
}

func (sess *Session) DeleteFrom(table string) *DeleteBuilderSession {
	return &DeleteBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		DeleteBuilder: DeleteFrom(table),
		LimitCount:    -1,
	}
}

func (tx *Tx) DeleteFrom(table string) *DeleteBuilderSession {
	return &DeleteBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		DeleteBuilder: DeleteFrom(table),
		LimitCount:    -1,
	}
}

func (sess *Session) DeleteBySql(query string, value ...interface{}) *DeleteBuilderSession {
	return &DeleteBuilderSession{
		runner:        sess,
		EventReceiver: sess,
		Dialect:       sess.Dialect,
		DeleteBuilder: DeleteBySQL(query, value...),
		LimitCount:    -1,
	}
}

func (tx *Tx) DeleteBySql(query string, value ...interface{}) *DeleteBuilderSession {
	return &DeleteBuilderSession{
		runner:        tx,
		EventReceiver: tx,
		Dialect:       tx.Dialect,
		DeleteBuilder: DeleteBySQL(query, value...),
		LimitCount:    -1,
	}
}

func (b *DeleteBuilderSession) ToSql() (string, []interface{}) {
	buf := NewBuffer()
	err := b.Build(b.Dialect, buf)
	if err != nil {
		panic(err)
	}
	return buf.String(), buf.Value()
}

func (b *DeleteBuilderSession) Exec() (sql.Result, error) {
	return exec(b.runner, b.EventReceiver, b, b.Dialect)
}

func (b *DeleteBuilderSession) Where(query interface{}, value ...interface{}) *DeleteBuilderSession {
	b.DeleteBuilder.Where(query, value...)
	return b
}

func (b *DeleteBuilderSession) Limit(n uint64) *DeleteBuilderSession {
	b.LimitCount = int64(n)
	return b
}

func (b *DeleteBuilderSession) Build(d Dialect, buf Buffer) error {
	err := b.DeleteBuilder.Build(b.Dialect, buf)
	if err != nil {
		return err
	}
	if b.LimitCount >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(fmt.Sprint(b.LimitCount))
	}
	return nil
}
