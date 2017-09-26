//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

//-----------------------------------------------------------------
var CreateUserTable = "create table if not exists user (id integer primary key, uid text unique, points integer)"
var CreateBroadcastTable = "create table if not exists broadcast (id integer primary key, bid text unique, content blob, date timestamp)"
var CreateSubmissionTable = "create table if not exists submission (id integer primary key, sid text unique, uid text, bid text, points integer, duration float, description text, language text, date timestamp, content blob)"
var CreatePollTable = "create table if not exists poll (id integer primary key, uid text, points integer, date timestamp)"
var database, _ = sql.Open("sqlite3", "./c4b.db")
var InsertBroadCastSQL, _ = database.Prepare("insert into broadcast (bid, content, date) values (?, ?, ?)")
var InsertUserSQL, _ = database.Prepare("insert into user (uid, points) values (?, ?)")
var InsertSubmissionSQL, _ = database.Prepare("insert into submission (sid, uid, bid, points, duration, description, language, date, content) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
var InsertPollSQL, _ = database.Prepare("insert into poll (uid, points, date) values (?, ?, ?)")
var UpdatePointsSQL, _ = database.Prepare("update submission set points=? where sid=?")

//-----------------------------------------------------------------
func init_sqldb() {
	sql_stmt, _ := database.Prepare(CreateUserTable)
	sql_stmt.Exec()
	sql_stmt, _ = database.Prepare(CreateBroadcastTable)
	sql_stmt.Exec()
	sql_stmt, _ = database.Prepare(CreateSubmissionTable)
	sql_stmt.Exec()
	sql_stmt, _ = database.Prepare(CreatePollTable)
	sql_stmt.Exec()
}

//-----------------------------------------------------------------
func show_users() {
	fmt.Println("begin user table")
	rows, _ := database.Query("SELECT id, uid, points FROM user")
	var id int
	var uid string
	var points int
	for rows.Next() {
		rows.Scan(&id, &uid, &points)
		fmt.Println(id, ";", uid, ";", points)
	}
	fmt.Println("end user table")
}

//-----------------------------------------------------------------
func show_broadcasts() {
	fmt.Println("begin broadcast table")
	rows, _ := database.Query("SELECT id, bid, content, date FROM broadcast")
	var id int
	var bid string
	var t time.Time
	var content []byte
	for rows.Next() {
		rows.Scan(&id, &bid, &content, &t)
		fmt.Println(id, ";", bid, ";", t)
		// fmt.Println(id, ";", bid, ";", t, ";", string(content))
	}
	fmt.Println("end broadcast table")
}

//-----------------------------------------------------------------
func show_submissions() {
	fmt.Println("begin submission table")
	rows, _ := database.Query("SELECT id, sid, uid, bid, points, description FROM submission")
	var id, points int
	var bid string
	var sid string
	var uid string
	var description string
	for rows.Next() {
		rows.Scan(&id, &sid, &uid, &bid, &points, &description)
		fmt.Println(id, ";", sid, ";", uid, ";", bid, ";", points, ";", description)
	}
	fmt.Println("end submission table")
}

//-----------------------------------------------------------------
func RegisterStudent(uid string) {
	SEM.Lock()
	defer SEM.Unlock()
	if _, ok := Boards[uid]; ok {
		fmt.Println(uid + " is already registered.")
		return
	}
	var err error
	var outFile *os.File
	if _, err = os.Stat(USER_DB); err == nil {
		outFile, err = os.OpenFile(USER_DB, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	} else {
		outFile, err = os.Create(USER_DB)
	}
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	t := time.Now()
	fmt.Println(t.Format("Mon Jan 2 15:04:05 MST 2006: write data to ") + USER_DB)
	w := csv.NewWriter(outFile)
	record := []string{uid, "0", "0", "Register", "", "Register",
		time.Now().Format("Mon Jan 2 15:04:05 MST 2006")}
	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
	}

	// Add a board for this new student
	Boards[uid] = &Board{"", "", time.Now(), false, "", ""}
	// TODO
	// Initialize with content of Boards["*"]

	_, err = InsertUserSQL.Exec(uid, 0)
	if err != nil {
		fmt.Println("Error inserting into user table.", err)
	} else {
		fmt.Println(uid, "is inserted.")
	}

}

//-----------------------------------------------------------------
func writeDB() {
	var err error
	var outFile *os.File
	if _, err = os.Stat(USER_DB); err == nil {
		outFile, err = os.OpenFile(USER_DB, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	} else {
		outFile, err = os.Create(USER_DB)
	}
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	t := time.Now()
	fmt.Println(t.Format("Mon Jan 2 15:04:05 MST 2006: write data to ") + USER_DB)
	w := csv.NewWriter(outFile)
	for _, sub := range ProcessedSubs {
		record := []string{
			sub.Uid,
			strconv.Itoa(sub.Points),
			strconv.Itoa(sub.Duration),
			sub.Sid,
			sub.Bid,
			sub.Pdes,
			sub.Timestamp,
		}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
	}
}

//-----------------------------------------------------------------
func initDB() {
	outFile, err := os.OpenFile(USER_DB, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer outFile.Close()
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(outFile)
	_, err = w.WriteString("uid,points,duration,sid,bid,des,timestamp\n")
	if err != nil {
		panic(err)
	}
	w.Flush()
}

//-----------------------------------------------------------------
func loadDB() (bool, map[string]*Submission) {
	var userFile *os.File
	var err error
	entries := make(map[string]*Submission)

	if _, err = os.Stat(USER_DB); os.IsNotExist(err) {
		userFile, err = os.Create(USER_DB)
		if err != nil {
			panic(err)
		}
	} else {
		userFile, err = os.Open(USER_DB)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer userFile.Close()
	reader := csv.NewReader(userFile)

	// Skip header
	record, err := reader.Read()
	empty_file := false
	if len(record) == 0 {
		empty_file = true
	}
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		uid := record[0]
		points, _ := strconv.Atoi(record[1])
		duration, _ := strconv.Atoi(record[2])
		sid := record[3]
		bid := record[4]
		des := record[5]
		timestamp := record[6]
		if sid == "Register" {
			Boards[uid] = &Board{"", "", time.Now(), false, "", ""}
		} else {
			s := &Submission{
				Sid:       sid,
				Bid:       bid,
				Uid:       uid,
				Points:    points,
				Duration:  duration,
				Pdes:      des,
				Timestamp: timestamp,
			}
			entries[sid] = s
		}
	}
	return empty_file, entries
}
