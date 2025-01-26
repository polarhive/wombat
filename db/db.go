package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// creates tables if necessary.
func InitializeSQLiteDB(dbFile string) {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	createTables()
}

// createTables creates the necessary tables in the database.
func createTables() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY
	);
	CREATE TABLE IF NOT EXISTS links (
		source TEXT,
		target TEXT,
		FOREIGN KEY (source) REFERENCES nodes(id),
		FOREIGN KEY (target) REFERENCES nodes(id)
	);
	`
	db.Exec(createTableSQL)
}

// ensure no duplicates.
func SaveNode(id string) {
	if id == "" {
		return
	}

	stmt, err := db.Prepare("INSERT OR IGNORE INTO nodes(id) VALUES(?)")
	if err != nil {
		log.Println("Error preparing SaveNode statement:", err)
		return
	}
	defer stmt.Close()
	stmt.Exec(id)
}

// saves a link between two nodes into the database.
func SaveLink(source, target string) {
	if source == "" || target == "" {
		return
	}

	stmt, err := db.Prepare("INSERT INTO links(source, target) VALUES(?, ?)")
	if err != nil {
		log.Println("Error preparing SaveLink statement:", err)
		return
	}
	defer stmt.Close()
	stmt.Exec(source, target)
}

func CloseDB() {
	db.Close()
}
