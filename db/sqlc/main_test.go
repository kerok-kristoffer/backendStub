package db

import (
	"database/sql"
	"github.com/jaswdr/faker"
	"github.com/kerok-kristoffer/formulating/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var F = faker.New()
var testDB *sql.DB

func TestMain(m *testing.M) {
	// todo: Add models for Recipes and Inventory
	// todo: Add test files for the rest of the models

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalln("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatalln("Cannot connect to db", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
