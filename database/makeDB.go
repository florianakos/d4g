package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	os.Remove("enigma.db")
	db, err := sql.Open("sqlite3", "./enigma.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createEnigmaTable := 	`CREATE TABLE enigma ( id TEXT PRIMARY KEY,
																							first_name TEXT NOT NULL,
																							last_name TEXT NOT NULL,
																							email TEXT UNIQUE NOT NULL,
																							gender TEXT NOT NULL,
																							address TEXT NOT NULL,
																							city TEXT NOT NULL,
																							phone TEXT NOT NULL,
																							image TEXT,
																							openings TEXT,
																							specialty TEXT);`

	_, err = db.Exec(createEnigmaTable)
	if err != nil {
			log.Println(err)
			return
	}
	fmt.Printf("Table %s created\n", createEnigmaTable)
}
