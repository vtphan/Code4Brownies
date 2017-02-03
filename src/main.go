//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var ADDR = ""
var PORT = "4030"
var USER_DB string
var SERVER = ""
//-----------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err.Error() + "\n")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			// fmt.Println("Server address http://" + ipnet.IP.String() + ":" + PORT)
			return ipnet.IP.String() + ":" + PORT
		}
	}
	return ""
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
		record := []string{sub.Uid, sub.Pid, strconv.Itoa(sub.Points), strconv.Itoa(sub.Duration), sub.Sid}
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
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		uid := record[0]
		pid := record[1]
		points, err := strconv.Atoi(record[2])
		duration, err := strconv.Atoi(record[3])
		sid := record[4]
		s := &Submission{Uid: uid, Pid: pid, Points: points, Duration: duration, Sid: sid}
		entries[sid] = s
	}
	return entries
}

//-----------------------------------------------------------------

func prepareCleanup() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-quit:
				fmt.Println("Preparing to stop server...")
				writeDB()
				os.Exit(1)
			}
		}
	}()
}

//-----------------------------------------------------------------
func main() {
	SERVER = informIPAddress()
	fmt.Println("Server address:", "http://" + SERVER)

	rand.Seed(time.Now().UnixNano())
	USER_DB = filepath.Join(".", "C4B_DB.csv")
	flag.StringVar(&USER_DB, "db", USER_DB, "user database in csv format, which consists of UID,POINTS.")
	flag.Parse()
	prepareCleanup()

	// student handlers
	http.HandleFunc("/submit_post", submit_postHandler) // rename this
	http.HandleFunc("/my_points", my_pointsHandler)
	http.HandleFunc("/receive_broadcast", receive_broadcastHandler)
	http.HandleFunc("/query_poll", query_pollHandler)

	// teacher handlers
	http.HandleFunc("/points", pointsHandler)
	http.HandleFunc("/give_point", give_pointHandler)
	http.HandleFunc("/peek", peekHandler)
	http.HandleFunc("/broadcast", broadcastHandler)
	http.HandleFunc("/get_post", get_postHandler)
	http.HandleFunc("/get_posts", get_postsHandler)
	http.HandleFunc("/start_poll", start_pollHandler)
	http.HandleFunc("/poll", view_pollHandler)
	err := http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
