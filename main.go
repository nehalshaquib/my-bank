package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nehalshaquib/my-bank/api"
	db "github.com/nehalshaquib/my-bank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/myBank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Println("cannot connect to db:", err)
		return
	}

	if err = api.NewServer(db.NewStore(conn)).Start(); err != nil {
		log.Fatalln("error in starting server: ", err)
	}

}
