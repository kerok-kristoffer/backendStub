package db

import (
	"database/sql"
	"github.com/jaswdr/faker"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:eloh@localhost:5454/formulating?sslmode=disable"
)

var testQueries *Queries
var F = faker.New()

func TestMain(m *testing.M) {
	// todo: Add test files for the rest of the models
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatalln("Cannot connect to db", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
