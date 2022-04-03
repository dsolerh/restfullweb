package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Book is a placeholder for books
type Book struct {
	id     int
	name   string
	author string
}

func main() {
	db, err := sql.Open("mysql", "test:pass@tcp(localhost:3306)/test?parseTime=true")
	if err != nil {
		log.Println(err)
	}
	// Create table
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY, 
		isbn INTEGER,
		author VARCHAR(64),
		name VARCHAR(64) NULL
	)`)

	if err != nil {
		log.Println("Error in create table", err)
	} else {
		log.Println("Successfully created table books!")
	}
	statement.Exec()
	dbOperations(db)
}

func dbOperations(db *sql.DB) {
	// Create
	statement, _ := db.Prepare("INSERT INTO books (name, author, isbn) VALUES (?,?,?)")
	r, err := statement.Exec("A Tale of Two Cities", "Charles Dickens", 140430547)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(r)
	}
	log.Println("Inserted the book into database!")

	// Read all
	rows, _ := db.Query("SELECT id, name, author FROM books")
	var tempBook Book
	for rows.Next() {
		fmt.Printf("\"query databasa\": %v\n", "query database")
		rows.Scan(&tempBook.id, &tempBook.name, &tempBook.author)
		log.Printf("ID:%d, Book:%s, Author:%s\n", tempBook.id, tempBook.name, tempBook.author)
	}

	// Update
	statement, _ = db.Prepare("update books set name=? where id=?")
	statement.Exec("The Tale of Two Cities", 1)
	log.Println("Successfully updated the book in database!")

	//Delete
	// statement, _ = db.Prepare("delete from books where id=?")
	// statement.Exec(1)
	// log.Println("Successfully deleted the book in database!")

}
