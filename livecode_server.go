//
// Live Coding (server)
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
	"flag"
	"encoding/csv"
	"log"
	"io"
	// "path/filepath"
	// "strings"
)

var PORT = "4030"

//-----------------------------------------------------------------
// ENTRIES
//-----------------------------------------------------------------

type Entry struct {
	User string
	Body string
	N    int
}

type EntryList struct {
	list []*Entry
	m    sync.Mutex
}

func NewEntryList() *EntryList {
	return &EntryList{}
}

func (E *EntryList) Len() int {
	return len(E.list)
}

func (E *EntryList) Add(user, body string) {
	E.m.Lock()
	E.list = append(E.list, &Entry{user, body, 0})
	E.m.Unlock()
}

func (E *EntryList) Remove(i int) *Entry {
	if i < 0 || len(E.list) == 0 || i > len(E.list) {
		return &Entry{}
	} else {
		E.m.Lock()
		entry := E.list[i]
		E.list = append(E.list[:i], E.list[i+1:]...)
		E.m.Unlock()
		return entry
	}
}

func (E *EntryList) Get(i int) *Entry {
	return E.list[i]
}

func (E *EntryList) Show() {
	for i := 0; i < len(E.list); i++ {
		fmt.Println(i, E.list[i])
	}
}

//-----------------------------------------------------------------
// USERS
//-----------------------------------------------------------------

type UserType struct {
	Points map[string]int
}

func NewUsers() *UserType {
	U := &UserType{}
	U.Points = make(map[string]int)
	return U
}

func (U *UserType) OnePoint(usr string) {
	score, ok := U.Points[usr]
	if !ok {
		U.Points[usr] = 1
	} else {
		U.Points[usr] = score + 1
	}
}

func (U *UserType) GetPoints(usr string) int {
	_, ok := U.Points[usr]
	if !ok {
		U.Points[usr] = 0
	}
	return U.Points[usr]
}

func (U *UserType) Show() {
	for key, value := range U.Points {
		fmt.Println(key, "\t", value)
	}
}

//-----------------------------------------------------------------
// HTTP HANDLERS
//-----------------------------------------------------------------

// Clients share their codes by POSTing to server_address/share
func pointsHandler(w http.ResponseWriter, r *http.Request) {
	_, user := r.FormValue("login"), r.FormValue("username")
	points, ok := Users.Points[user]
	if !ok {
		points = 0
	}
	fmt.Fprintf(w, user+" has "+strconv.Itoa(points)+" points.")
}

// Clients share their codes by POSTing to server_address/share
func shareHandler(w http.ResponseWriter, r *http.Request) {
	_, user, body := r.FormValue("login"), r.FormValue("username"), r.FormValue("body")
	if _, ok := ExistingUsers.Points[user]; ok {
		Entries.Add(user, body)
		fmt.Println(user, "submitted.")
		fmt.Fprintf(w, "1")
	} else {
		fmt.Println(user, "non existent.")
		fmt.Fprintf(w, "0")
	}
}

// give one brownie point to the user that is just recently dequed
func brownieHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Users.OnePoint(r.FormValue("user"))
		Users.Show()
		fmt.Fprintf(w, "Point awarded to "+r.FormValue("user"))
	}
}

func entriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		js, err := json.Marshal(Entries.list)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

// Instructor retrieves code, one by one, by invoking server_address/deque
func request_entryHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		e, err := strconv.Atoi(r.FormValue("entry"))
		if err != nil {
			fmt.Println(err.Error)
		} else {
			js, err := json.Marshal(Entries.Remove(e))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
		}
	}
}

//-----------------------------------------------------------------
func prepareCleanup(userFile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		t := time.Now()
		fmt.Println("\n" + t.Format("Mon Jan 2 15:04:05 MST 2006") + "\nTOTAL POINTS")
		Users.Show()
		fmt.Println("\nUpdating",userFile,"\n")
		for user, point := range(Users.Points) {
			ExistingUsers.Points[user] += point
		}
	   outFile, err := os.Create(userFile)
	   if err != nil {
	   	panic(err)
	   }
	   defer outFile.Close()
		w := csv.NewWriter(outFile)
		for user, point := range(ExistingUsers.Points) {
			record := []string{user,strconv.Itoa(point)}
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
func informIPAddress() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err.Error() + "\n")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Server address " + ipnet.IP.String() + ":" + PORT)
			}
		}
	}
}

//-----------------------------------------------------------------
func loadExistingRecord(filename string) {
	userFile, err := os.Open(filename)
	defer userFile.Close()

	var points int
   if err != nil {
   	log.Fatal(err)
   }
	r := csv.NewReader(userFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		points, err = strconv.Atoi(record[1])
		ExistingUsers.Points[record[0]] = points
	}
	for k, v := range(ExistingUsers.Points) {
		fmt.Println(k,v)
	}
}

//-----------------------------------------------------------------
// MAIN
//-----------------------------------------------------------------

var PASSCODE string
var Entries = NewEntryList()
var Users = NewUsers()
var ExistingUsers = NewUsers()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var userFilename string
	flag.StringVar(&PASSCODE, "passcode", "password", "passcode to be used by the instructor to connect to the server.")
	flag.StringVar(&userFilename, "users", "", "csv-formatted filename containing usernames,real names.")
	flag.Parse()

	loadExistingRecord(userFilename)
	prepareCleanup(userFilename)
	informIPAddress()

	// Register handlers and start serving app
	http.HandleFunc("/share", shareHandler)
	http.HandleFunc("/points", pointsHandler)
	http.HandleFunc("/brownie", brownieHandler)
	http.HandleFunc("/entries", entriesHandler)
	http.HandleFunc("/request_entry", request_entryHandler)
	err := http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
