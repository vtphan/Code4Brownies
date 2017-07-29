//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

//-----------------------------------------------------------------
func RegisterStudent(uid string) {
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
	record := []string{uid, "0", "0", "Register", "Register",
		time.Now().Format("Mon Jan 2 15:04:05 MST 2006")}
	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
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
func loadDB() map[string]*Submission {
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
	reader := csv.NewReader(userFile)
	entries := make(map[string]*Submission)

	// Skip header
	reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		uid := record[0]
		points, err := strconv.Atoi(record[1])
		duration, err := strconv.Atoi(record[2])
		sid := record[3]
		des := record[4]
		timestamp := record[5]
		if sid == "Register" {
			Boards[uid] = &Board{"", "", time.Now(), false}
		} else {
			s := &Submission{
				Uid:       uid,
				Points:    points,
				Duration:  duration,
				Sid:       sid,
				Pdes:      des,
				Timestamp: timestamp,
			}
			entries[sid] = s
		}
	}
	return entries
}
