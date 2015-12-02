package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"os"
	"os/signal"
	"syscall"
	"sync"
	"path/filepath"
	"time"
)

var PORT = "4030"

//-----------------------------------------------------------------
// ENTRIES
//-----------------------------------------------------------------

type Entry struct {
	User string
	Body string
	N int
}

type EntryList struct {
	list []*Entry
	m sync.Mutex
}

func NewEntryList() *EntryList {
	return &EntryList{}
}

func (E *EntryList) Len() int {
	return len(E.list)
}

func (E *EntryList) Add(user, body string) {
	E.m.Lock()
	E.list = append(E.list, &Entry{user,body,0})
	E.m.Unlock()
}

func (E *EntryList) Deque()  *Entry {
	if len(E.list) == 0 {
		return &Entry{}
	} else {
		E.m.Lock()
		CurEntry = E.list[0]
		E.list = E.list[1:]
		CurEntry.N = len(E.list)
		E.m.Unlock()
		return CurEntry
	}
}

func (E *EntryList) Show() {
	for i:=0; i<len(E.list); i++ {
		fmt.Println(i, E.list[i])
	}
}


//-----------------------------------------------------------------
// USERS
//-----------------------------------------------------------------

type UserType struct {
	Points map[string]int
	m sync.Mutex
}

func NewUsers() *UserType {
	U := &UserType{}
	U.Points = make(map[string]int)
	return U
}

func (U *UserType) OnePoint(usr string) {
	score, ok := U.Points[usr]
	U.m.Lock()
	if !ok {
		U.Points[usr] = 1
	} else {
		U.Points[usr] = score + 1
	}
	U.m.Unlock()
}

func (U *UserType) GetPoints(usr string) int {
	_, ok := U.Points[usr]
	if !ok {
		U.Points[usr] = 0
	}
	return U.Points[usr]
}

func (U *UserType) Show() {
	for key,value := range U.Points {
		fmt.Println(key,"\t",value)
	}
}


//-----------------------------------------------------------------
// HTTP HANDLERS
//-----------------------------------------------------------------

// Clients share their codes by POSTing to server_address/share
func shareHandler(w http.ResponseWriter, r *http.Request) {
	_, user, body := r.FormValue("login"), r.FormValue("username"), r.FormValue("body")
	Entries.Add(user, body)
	fmt.Println("+", user, ":", Entries.Len(), "entries in queue.")
	fmt.Fprintf(w, "Got it!")
}

// Instructor retrieves code, one by one, by invoking server_address/deque
func dequeHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		entry := Entries.Deque()
		js, err := json.Marshal(entry)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("-", entry.User, ":", Entries.Len(), "entries left.")
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

func currentEntryHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		js, err := json.Marshal(CurEntry)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

// give one brownie point to the user that is just recently dequed
func brownieHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		if CurEntry != nil {
			Users.OnePoint(CurEntry.User)
			Users.Show()
			fmt.Fprintf(w, "One point awarded to " + CurEntry.User)
		} else {
			fmt.Fprintf(w, "No entry has been dequed.")
		}
	}
}

//-----------------------------------------------------------------
// MAIN
//-----------------------------------------------------------------

var PASSCODE string
var CurEntry *Entry
var Entries = NewEntryList()
var Users = NewUsers()

func main() {
	// Get passcode
	_, gofile := filepath.Split(os.Args[0])
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run "+gofile+" passcode")
		os.Exit(1)
	}
	PASSCODE = os.Args[1]
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Prepare for cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		t := time.Now()
		fmt.Println("\n" + t.Format("Mon Jan 2 15:04:05 MST 2006") +"\nTOTAL POINTS")
		Users.Show()
		os.Exit(1)
	}()

	// Figure IP address
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

	// Register handlers and start serving app
	http.HandleFunc("/share", shareHandler)
	http.HandleFunc("/deque", dequeHandler)
	http.HandleFunc("/currentEntry", currentEntryHandler)
	http.HandleFunc("/brownie", brownieHandler)
	err = http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
