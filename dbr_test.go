package dbr

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/ql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

//
// Test helpers
//

var (
	currID uint64 = 256
)

// create id
func nextID() uint64 {
	currID++
	return currID
}

const (
	mysqlDSN    = "root:unprotected@unix(/tmp/mysql.sock)/uservoice_development?charset=utf8&parseTime=true"
	postgresDSN = "postgres://postgres:unprotected@localhost:5432/uservoice_development?sslmode=disable"
)

func createSession(driver, dsn string) *Session {
	var testDSN string
	switch driver {
	case "mysql":
		testDSN = os.Getenv("DBR_TEST_MYSQL_DSN")
	case "postgres":
		testDSN = os.Getenv("DBR_TEST_POSTGRES_DSN")
	}
	if testDSN != "" {
		dsn = testDSN
	}
	conn, err := Open(driver, dsn, nil)
	if err != nil {
		log.Fatal(err)
	}
	reset(conn)
	sess := conn.NewSession(nil)
	return sess
}

var (
	mysqlSession    = createSession("mysql", mysqlDSN)
	postgresSession = createSession("postgres", postgresDSN)
)

type dbrPerson struct {
	ID    uint64
	Name  string
	Email string
}

type nullTypedRecord struct {
	ID         uint64
	StringVal  NullString
	Int64Val   NullInt64
	Float64Val NullFloat64
	TimeVal    NullTime
	BoolVal    NullBool
}

func reset(conn *Connection) {
	// serial = BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE
	// the following sql should work for both mysql and postgres
	createPeopleTable := `
		CREATE TABLE dbr_people (
			id serial PRIMARY KEY,
			name varchar(255) NOT NULL,
			email varchar(255)
		)
	`

	createNullTypesTable := `
		CREATE TABLE null_types (
			id serial PRIMARY KEY,
			string_val varchar(255) NULL,
			int64_val integer NULL,
			float64_val float NULL,
			time_val timestamp NULL ,
			bool_val bool NULL
		)
	`

	for _, v := range []string{
		"DROP TABLE IF EXISTS dbr_people",
		createPeopleTable,

		"DROP TABLE IF EXISTS null_types",
		createNullTypesTable,
	} {
		_, err := conn.Exec(v)
		if err != nil {
			log.Fatal("Failed to execute statement: ", v, " Got error: ", err)
		}
	}
}

func TestBasicCRUD(t *testing.T) {
	jonathan := dbrPerson{
		ID:    nextID(),
		Name:  "jonathan",
		Email: "jonathan@uservoice.com",
	}
	for _, sess := range []SessionRunner{mysqlSession, postgresSession} {
		// insert
		result, err := sess.InsertInto("dbr_people").Columns("id", "name", "email").Record(jonathan).Exec()
		assert.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.EqualValues(t, 1, rowsAffected)

		// select
		var people []dbrPerson
		err = sess.Select("*").From("dbr_people").Where(ql.Eq("id", jonathan.ID)).Load(&people)
		assert.NoError(t, err)
		assert.Equal(t, len(people), 1)
		assert.Equal(t, jonathan.ID, people[0].ID)
		assert.Equal(t, jonathan.Name, people[0].Name)
		assert.Equal(t, jonathan.Email, people[0].Email)

		// update
		result, err = sess.Update("dbr_people").Where(ql.Eq("id", jonathan.ID)).Set("name", "jonathan1").Exec()
		assert.NoError(t, err)

		rowsAffected, err = result.RowsAffected()
		assert.NoError(t, err)
		assert.EqualValues(t, 1, rowsAffected)

		// delete
		result, err = sess.DeleteFrom("dbr_people").Where(ql.Eq("id", jonathan.ID)).Exec()
		assert.NoError(t, err)

		rowsAffected, err = result.RowsAffected()
		assert.NoError(t, err)
		assert.EqualValues(t, 1, rowsAffected)
	}
}
