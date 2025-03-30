package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/definitely-unique-username/simple_bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..", ".env")

	if err != nil {
		log.Fatal("failed to read config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
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
