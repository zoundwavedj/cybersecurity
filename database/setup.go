package database

import (
	"database/sql"

	// Required for sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// Db global var
var Db *sql.DB

// Setup function for initiating db connection and preparing schema
func Setup() {
	var err error

	Db, err = sql.Open("sqlite3", "file:local.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	statement, err := Db.Prepare("CREATE TABLE IF NOT EXISTS superuser (id TEXT PRIMARY KEY, username TEXT UNIQUE, password TEXT)")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	statement.Close()

	statement, err = Db.Prepare("CREATE TABLE IF NOT EXISTS authentication (userId TEXT, token TEXT, expiresAt INTEGER, FOREIGN KEY (userId) REFERENCES superuser(id))")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	statement.Close()

	statement, err = Db.Prepare("CREATE TABLE IF NOT EXISTS user (id TEXT PRIMARY KEY, name TEXT, dob TEXT, email TEXT, ssn TEXT)")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	statement.Close()
}
