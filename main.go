package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/definitely-unique-username/simple_bank/api"
	db "github.com/definitely-unique-username/simple_bank/db/sqlc"
	"github.com/definitely-unique-username/simple_bank/util"
)

func main() {
	config, err := util.LoadConfig("./", ".env")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.SevrverAddress); err != nil {
		log.Fatal("Cannot start server", err)
	}
}
