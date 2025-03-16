package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Use test database URL from environment if provided
	dbURL := os.Getenv("SIMPLE_BANK_DB_SOURCE")
	if dbURL == "" {
		dbURL = dbSource
	}

	conn, err := sql.Open(dbDriver, dbURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testDB = conn
	testQueries = New(conn)

	os.Exit(m.Run())
}

// beginTx creates a new transaction and returns a Queries instance that uses it
func beginTx(t *testing.T) (*Queries, func()) {
	t.Helper()

	// Start a transaction
	tx, err := testDB.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new Queries instance with the transaction
	q := New(tx)

	// Return the Queries instance and a cleanup function
	return q, func() {
		// Always rollback the transaction in cleanup
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			t.Error(err)
		}
	}
}
