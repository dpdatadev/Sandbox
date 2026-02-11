package main

// TODO, refactor into our common DAL (reuseable CRUD code)

import (
	"database/sql"
	"log"
	"runtime/debug"

	_ "github.com/mattn/go-sqlite3"
)

// Ensure DB is created
func CreateDB() {
	db, err := sql.Open("sqlite3", "./test1.db")
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}
	defer db.Close()
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS dispatch_log (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        uuid TEXT,
		notes TEXT
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}
	log.Println("Table 'dispatch_log' created successfully")

}

func SeedDB(uuid string, notes string) {
	db, err := sql.Open("sqlite3", "./test1.db")
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO dispatch_log(uuid, notes) VALUES(?, ?)", uuid, notes)
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}
	log.Println("New log inserted successfully")
}

func QueryDBTest() {
	db, err := sql.Open("sqlite3", "./test1.db")
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, uuid, notes from dispatch_log")
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var uuid string
		var notes string
		err = rows.Scan(&id, &uuid, &notes)
		if err != nil {
			log.Fatal(err)
			debug.PrintStack()
		}
		log.Printf("Log: %d, %s, %s", id, uuid, notes)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
		debug.PrintStack()
	}
}
