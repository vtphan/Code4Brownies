# SQL example with golang

package main

import (
	"fmt"
	"strconv"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var database, _ = sql.Open("sqlite3", "./nraboy.db")
var CreateStmt, _ = database.Prepare("create table if not exists people (id integer primary key, firstname text, lastname text)")
var InsertStmt, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")

func main() {
	CreateStmt.Exec()
	InsertStmt.Exec("Jason", "Bourne")
	rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}
}
