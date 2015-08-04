# gocraft/dbr (database records) [![GoDoc](https://godoc.org/github.com/gocraft/web?status.png)](https://godoc.org/github.com/gocraft/dbr)

gocraft/dbr provides additions to Go's database/sql for super fast performance and convenience.

## Getting Started

```go
// create a connection
conn, _ := dbr.Open("postgres", "...")

// create a session for each business unit of execution (e.g. a web request or goworkers job)
sess := conn.NewSession(nil)

// get a record
var suggestion Suggestion
sess.Select("id", "title").From("suggestions").Where("id = ?", 1).Load(&suggestion)

// JSON-ready, with dbr.Null* types serialized like you want
b, _ := json.Marshal(&suggestion)
fmt.Println(string(b))
```

## Feature highlights

### Automatically map results to structs
Querying is the heart of gocraft/dbr. Automatically map results to structs:

```go
var suggestion Suggestion
sess.Select("id", "title", "body").From("suggestions").Where("id = ?", 1).Load(&suggestion)
```

Additionally, easily query a single value or a slice of values:

```go
var suggestions []Suggestion
sess.Select("id", "title", "body").From("suggestions").OrderBy("id", ql.ASC).Load(&suggestions)
```


See below for many more examples.

### Use a Sweet Query Builder or use Plain SQL
gocraft/dbr supports both.

Sweet Query Builder:
```go

builder := ql.Select("title", "body").
	From("suggestions").
	OrderBy("id", ql.ASC).
	Limit(10)
```

Plain SQL:

```go
builder := ql.SelectBySQL("SELECT `title`, `body` FROM `suggestions` ORDER BY `id` ASC LIMIT 10")
```

### IN queries that aren't horrible
Traditionally, database/sql uses prepared statements, which means each argument in an IN clause needs its own question mark. gocraft/dbr, on the other hand, handles interpolation itself so that you can easily use a single question mark paired with a dynamically sized slice.

```go
ids := []int64{1, 2, 3, 4, 5}
builder.Where("id IN ?", ids) // `id` IN ?
```

### Amazing instrumentation
Writing instrumented code is a first-class concern for gocraft/dbr. We instrument each query to emit to a gocraft/health-compatible EventReceiver interface.

### Faster performance than using using database/sql directly
Every time you call database/sql's db.Query("SELECT ...") method, under the hood, the mysql driver will create a prepared statement, execute it, and then throw it away. This has a big performance cost.

gocraft/dbr doesn't use prepared statements. We ported mysql's query escape functionality directly into our package, which means we interpolate all of those question marks with their arguments before they get to MySQL. The result of this is that it's way faster, and just as secure.

Check out these [benchmarks](https://github.com/tyler-smith/golang-sql-benchmark).

### JSON Friendly
Every try to JSON-encode a sql.NullString? You get:
```json
{
	"str1": {
		"Valid": true,
		"String": "Hi!"
	},
	"str2": {
		"Valid": false,
		"String": ""
  }
}
```

Not quite what you want. gocraft/dbr has dbr.NullString (and the rest of the Null* types) that encode correctly, giving you:

```json
{
	"str1": "Hi!",
	"str2": null
}
```

## Driver support

* MySQL
* PostgreSQL

## Usage Examples

### Making a session
All queries in gocraft/dbr are made in the context of a session. This is because when instrumenting your app, it's important to understand which business action the query took place in. See gocraft/health for more detail.

Here's an example web endpoint that makes a session:

### Simple Record CRUD

See `TestBasicCRUD`.

### Overriding Column Names With Struct Tags

```go
// By default dbr converts CamelCase property names to snake_case column_names
// You can override this with struct tags, just like with JSON tags
// This is especially helpful while migrating from legacy systems
type Suggestion struct {
	Id        int64
	Title     dbr.NullString `db:"subject"` // subjects are called titles now
	CreatedAt dbr.NullTime
}
```

### Embedded structs

```go
// columns are mapped by tag then by field
type Suggestion struct {
	ID uint64  // id
	Title string // title
	Body dbr.NullString `db:"content"` // content
	User User
}

type User struct {
	ID // user_id
	Name `db:"name"` // name; without tag, it will be user_name
}
```

### JSON encoding of Null* types
```go
// dbr.Null* types serialize to JSON like you want
suggestion := &Suggestion{Id: 1, Title: "Test Title"}
jsonBytes, err := json.Marshal(&suggestion)
fmt.Println(string(jsonBytes)) // {"id":1,"title":"Test Title","created_at":null}
```

### Inserting Multiple Records

```
sess.InsertInto("suggestions").Columns("title", "body")
	.Record(suggestion1)
	.Record(suggestion2)
```

### Updating Records

```go
sess.Update("suggestions").
	Set("title", "Gopher").
	Set("body", "I love go.").
	Where("id = ?", 1)
```

### Transactions

```go
tx, err := sess.Begin()
tx.Rollback()
```

## gocraft

gocraft offers a toolkit for building web apps. Currently these packages are available:

* [gocraft/web](https://github.com/gocraft/web) - Go Router + Middleware. Your Contexts.
* [gocraft/dbr](https://github.com/gocraft/dbr) - Additions to Go's database/sql for super fast performance and convenience.
* [gocraft/health](https://github.com/gocraft/health) -  Instrument your web apps with logging and metrics.

These packages were developed by the [engineering team](https://eng.uservoice.com) at [UserVoice](https://www.uservoice.com) and currently power much of its infrastructure and tech stack.

## Thanks & Authors
Inspiration from these excellent libraries:
*  [sqlx](https://github.com/jmoiron/sqlx) - various useful tools and utils for interacting with database/sql.
*  [Squirrel](https://github.com/lann/squirrel) - simple fluent query builder.

Authors:
*  Jonathan Novak -- [https://github.com/cypriss](https://github.com/cypriss)
*  Tyler Smith -- [https://github.com/tyler-smith](https://github.com/tyler-smith)
*  Sponsored by [UserVoice](https://eng.uservoice.com)
