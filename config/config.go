package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Dbconn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "my_note"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
