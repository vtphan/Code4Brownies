//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//-----------------------------------------------------------------
// USERS
//-----------------------------------------------------------------

type User struct {
	points int
}

const WRITE_TO_DB_INTERVAL = 6

var USER_DB string
var AllUsers = make(map[string]*User) // maps uids to users
var RegisteredUsers = make(map[string]string)
var DB_LOADED = false

//-----------------------------------------------------------------
func writeToUserDB() {
	if DB_LOADED == true {
		t := time.Now()

		outFile, err := os.Create(USER_DB)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()

		fmt.Println(t.Format("Mon Jan 2 15:04:05 MST 2006: write data to ") + USER_DB)
		w := csv.NewWriter(outFile)
		for uid, user := range AllUsers {
			record := []string{uid, strconv.Itoa(user.points)}
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			panic(err)
		}
	}
}

//-----------------------------------------------------------------

func prepareCleanup() {
	ticker := time.NewTicker(WRITE_TO_DB_INTERVAL * time.Second)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				writeToUserDB()
			case <-quit:
				fmt.Println("Preparing to stop server...")
				writeToUserDB()
				ticker.Stop()
				os.Exit(1)
			}
		}
	}()
}

//-----------------------------------------------------------------
// load records in csv file into global AllUsers
//-----------------------------------------------------------------
func loadRecords() {
	var userFile *os.File
	var err error

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

	// Create backup file
	backupFile, err := os.Create(USER_DB + ".bak")
	if err != nil {
		panic(err)
	}
	defer backupFile.Close()

	// Read into global AllUsers record
	reader := csv.NewReader(userFile)
	writer := csv.NewWriter(backupFile)
	var points int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		points, err = strconv.Atoi(record[1])
		user := &User{points}
		AllUsers[record[0]] = user
		if err := writer.Write(record); err != nil {
			log.Fatalln("error backing up record to csv:", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		panic(err)
	}
	fmt.Println("Duplicate db to", USER_DB+".bak")
	DB_LOADED = true
}
