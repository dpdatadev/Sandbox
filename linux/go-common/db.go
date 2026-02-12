package main

// TODO, refactor into our common DAL (reuseable CRUD code)

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Dispatch struct {
	uuid  string
	notes string
}

func CreateNewDispatch(notesString string) *Dispatch {
	return &Dispatch{
		uuid.NewString(),
		notesString,
	}
}

// Ensure DB is created
func CreateDB() {
	db, err := sql.Open("sqlite3", "./test1.db")
	if err != nil {
		log.Fatal(err)
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
	}
	log.Println("Table 'dispatch_log' created successfully")

}

func SeedDB(uuid string, notes string) {
	db, err := sql.Open("sqlite3", "./test1.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dispatch := CreateNewDispatch("This is a Test Dispatch!")
	_, err = db.Exec("INSERT INTO dispatch_log(uuid, notes) VALUES(?, ?)", &dispatch.uuid, &dispatch.notes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("New log inserted successfully")
}

func QueryDBTest() {
	db, err := sql.Open("sqlite3", "./test1.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, uuid, notes from dispatch_log")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		var dispatch Dispatch
		err = rows.Scan(&id, &dispatch.uuid, &dispatch.notes)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Log: %d, %s, %s", id, dispatch.uuid, dispatch.notes)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
