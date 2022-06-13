package main

import (
	"database/sql"
	"github.com/kerok-kristoffer/backendStub/util"
	"log"

	"github.com/kerok-kristoffer/backendStub/api"
	db "github.com/kerok-kristoffer/backendStub/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatalln("Cannot connect to db", err)
	}

	account := db.NewUserAccount(conn)
	server := api.NewServer(account)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatalln("cannot start server:", err)
	}
}
