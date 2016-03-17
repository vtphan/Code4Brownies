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

var USER_DB string
var AllUsers = make(map[string]*User) // maps uids to users
var RegisteredUsers = make(map[string]string)

//-----------------------------------------------------------------

func prepareCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		t := time.Now()
		fmt.Println("\n" + t.Format("Mon Jan 2 15:04:05 MST 2006"))

		outFile, err := os.Create(USER_DB)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()

		fmt.Println("Writing data to", USER_DB)
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
		os.Exit(1)
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
	fmt.Println("Duplicated db to", USER_DB+".bak")
}
