package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	DB    *sql.DB
	DBErr error
)

func Setup() {
	var db = DB

	dsn := "file:config/rss.db?cache=shared"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		DBErr = err
		fmt.Println("db err: ", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS rss (name text, link text, last text);
    CREATE TABLE IF NOT EXISTS banned_word (value text);
	DROP TABLE IF EXISTS messages_send;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	DB = db
}

// GetDB helps you to get a connection
func GetDB() *sql.DB {
	return DB
}
