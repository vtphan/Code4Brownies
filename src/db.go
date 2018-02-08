//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	// "time"
)

//-----------------------------------------------------------------
var database *sql.DB
var InsertBroadCastSQL *sql.Stmt
var InsertUserSQL *sql.Stmt
var InsertSubmissionSQL *sql.Stmt
var InsertPollSQL *sql.Stmt
var UpdatePointsSQL *sql.Stmt
var InsertAttendanceSQL *sql.Stmt

//-----------------------------------------------------------------
func init_db() {
	var err error
	prepare := func(s string) *sql.Stmt {
		stmt, err := database.Prepare(s)
		if err != nil {
			panic(err)
		}
		return stmt
	}

	database, err = sql.Open("sqlite3", USER_DB)
	if err != nil {
		panic(err)
	}

	create_tables()

	InsertBroadCastSQL = prepare("insert into broadcast (bid, content, language, date, hints) values (?, ?, ?, ?, ?)")
	InsertUserSQL = prepare("insert into user (uid) values (?)")
	InsertSubmissionSQL = prepare("insert into submission (sid, uid, bid, points, description, language, date, content, hints_used) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	InsertPollSQL = prepare("insert into poll (uid, is_correct, points, date) values (?, ?, ?, ?)")
	UpdatePointsSQL = prepare("update submission set points=? where sid=?")
	InsertAttendanceSQL = prepare("insert into attendance (uid, date) values (?, ?)")
}

//-----------------------------------------------------------------

func create_tables() {
	execSQL := func(s string) {
		sql_stmt, err := database.Prepare(s)
		if err != nil {
			panic(err)
		}
		sql_stmt.Exec()
	}
	execSQL("create table if not exists user (id integer primary key, uid text unique)")
	execSQL("create table if not exists broadcast (id integer primary key, bid text unique, content blob, language text, date timestamp, hints integer)")
	execSQL("create table if not exists submission (id integer primary key, sid text unique, uid text, bid text, points integer, description text, language text, date timestamp, content blob, hints_used integer)")
	execSQL("create table if not exists poll (id integer primary key, uid text, is_correct integer, points integer, date timestamp)")
	execSQL("create table if not exists attendance (id integer primary key, uid text, date timestamp)")
}

//-----------------------------------------------------------------
func RegisterStudent(uid string) {
	SEM.Lock()
	defer SEM.Unlock()
	if _, ok := Boards[uid]; ok {
		fmt.Println(uid + " is already registered.")
		return
	}
	Boards[uid] = make([]*Board, 0)
	for i := 0; i < len(Boards["__all__"]); i++ {
		b := &Board{
			Content:      Boards["__all__"][i].Content,
			HelpContent:  Boards["__all__"][i].HelpContent,
			Description:  Boards["__all__"][i].Description,
			StartingTime: Boards["__all__"][i].StartingTime,
			Ext:          Boards["__all__"][i].Ext,
			Bid:          Boards["__all__"][i].Bid,
		}
		Boards[uid] = append(Boards[uid], b)
	}

	_, err := InsertUserSQL.Exec(uid)
	if err != nil {
		fmt.Println("Error inserting into user table.", err)
	} else {
		fmt.Println("New user", uid)
	}
}

//-----------------------------------------------------------------
func loadWhiteboards() {
	rows, _ := database.Query("select uid from user")
	defer rows.Close()
	var uid string
	for rows.Next() {
		rows.Scan(&uid)
		Boards[uid] = make([]*Board, 0)
	}
	Boards["__all__"] = make([]*Board, 0)
}
