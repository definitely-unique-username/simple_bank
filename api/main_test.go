package api

import (
	"os"
	"testing"
	"time"

	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/definitely-unique-username/simple_bank/util"
	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SymmetricalKey:      util.RandString(32),
		AccessTokenDuration: time.Minute,
	}

	return NewServer(&config, store)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
