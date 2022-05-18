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
var testDB *sql.DB

func TestMain(m *testing.M) {
	// todo: Add models for Recipies and Inventory
	// todo: Add test files for the rest of the models
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatalln("Cannot connect to db", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
