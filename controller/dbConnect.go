package controller

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}