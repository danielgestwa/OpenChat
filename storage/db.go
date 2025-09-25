package storage

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func connect() *sql.DB {

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb

	var login string = os.Getenv("DB_USER")
	var password string = os.Getenv("DB_PASSWORD")
	var url string = os.Getenv("DB_ADDR")
	var dbName string = os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", login+":"+password+"@tcp("+url+")/"+dbName)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	return db

}

func exit(db *sql.DB) {
	db.Close()
}
