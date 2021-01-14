package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// global var because it will be extensively used for testing
var testQueries *Queries
var testDB *sql.DB

// test main function is the main entry point of all unit tests inside the package (e.g package db)
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	// New function is defined in db/sqlc/db.go
	testQueries = New(testDB)

	// Run the unit test; returns whether test passes or fails
	// Report it back to test runner via os.Exit command
	os.Exit(m.Run())
}
