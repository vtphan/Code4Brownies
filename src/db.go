//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

//-----------------------------------------------------------------
var database, _ = sql.Open("sqlite3", "./c4b.db")
var InsertBroadCastSQL, _ = database.Prepare("insert into broadcast (bid, content, language, date) values (?, ?, ?, ?)")
var InsertUserSQL, _ = database.Prepare("insert into user (uid) values (?)")
var InsertSubmissionSQL, _ = database.Prepare("insert into submission (sid, uid, bid, points, duration, description, language, date, content) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
var InsertPollSQL, _ = database.Prepare("insert into poll (uid, is_correct, points, date) values (?, ?, ?, ?)")
var UpdatePointsSQL, _ = database.Prepare("update submission set points=? where sid=?")

//-----------------------------------------------------------------
func exec_sql(s string) {
	sql_stmt, err := database.Prepare(s)
	if err != nil {
		panic(err)
	}
	sql_stmt.Exec()
}

//-----------------------------------------------------------------

func init_sqldb() {
	exec_sql("create table if not exists user (id integer primary key, uid text unique)")
	exec_sql("create table if not exists broadcast (id integer primary key, bid text unique, content blob, language text, date timestamp)")
	exec_sql("create table if not exists submission (id integer primary key, sid text unique, uid text, bid text, points integer, duration float, description text, language text, date timestamp, content blob)")
	exec_sql("create table if not exists poll (id integer primary key, uid text, is_correct integer, points integer, date timestamp)")
}

//-----------------------------------------------------------------
func RegisterStudent(uid string) {
	SEM.Lock()
	defer SEM.Unlock()
	if _, ok := Boards[uid]; ok {
		fmt.Println(uid + " is already registered.")
		return
	}
	Boards[uid] = &Board{
		Content:      Boards["__all__"].Content,
		Description:  Boards["__all__"].Description,
		StartingTime: time.Now(),
		Changed:      false,
		Ext:          Boards["__all__"].Ext,
		Bid:          Boards["__all__"].Bid,
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
		Boards[uid] = &Board{"", "", time.Now(), false, "", ""}
	}
	Boards["__all__"] = &Board{"", "", time.Now(), false, "", ""}
}
