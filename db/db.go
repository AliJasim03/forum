package db

import (
    "database/sql"
    "log"
    _ "golang.org/x/crypto/bcrypt"
    _ "github.com/gofrs/uuid"
    _ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
    db, err := sql.Open("sqlite3", "./forum.db")
    if err != nil {
        log.Fatal(err)
    }
    return db
}

func CloseDB(db *sql.DB) {
    db.Close()
}