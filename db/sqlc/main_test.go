package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/keremakillioglu/simplebank/util"

	_ "github.com/lib/pq"
)

// global var because it will be extensively used for testing
var testQueries *Queries
var testDB *sql.DB

// test main function is the main entry point of all unit tests inside the package (e.g package db)
func TestMain(m *testing.M) {
	// go to parent folder
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	// New function is defined in db/sqlc/db.go
	testQueries = New(testDB)

	// Run the unit test; returns whether test passes or fails
	// Report it back to test runner via os.Exit command
	os.Exit(m.Run())
}
