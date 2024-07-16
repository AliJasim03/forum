package db

import (
	"database/sql"
	_ "github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

func OpenConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "db/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	return initDB(db)
}

func CloseConnection(db *sql.DB) {
	_ = db.Close()
}

func initDB(db *sql.DB) *sql.DB {
	sqlFile, err := os.ReadFile("db/init.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
