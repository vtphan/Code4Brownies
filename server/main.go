package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"os"
	"os/signal"
	"syscall"
)

var PORT = "4030"
var Entries = NewEntryList()
var Users = NewUsers()

/* Handlers */

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "There are "+strconv.Itoa(Entries.Len())+" entries.")
}

// Clients share their codes by POSTing to server_address/share
func shareHandler(w http.ResponseWriter, r *http.Request) {
	_, user, body := r.FormValue("login"), r.FormValue("username"), r.FormValue("body")
	Entries.Add(user, body)
	fmt.Println("+", user, ":", Entries.Len(), "entries in queue.")
	fmt.Fprintf(w, "Got it!")
}

// Instructor retrieves code, one by one, by invoking server_address/deque
func dequeHandler(w http.ResponseWriter, r *http.Request) {
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

func currentEntryHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(CurEntry)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// give one brownie point to the user that is just recently dequed
func brownieHandler(w http.ResponseWriter, r *http.Request) {
	if CurEntry != nil {
		Users.OnePoint(CurEntry.User)
		Users.Show()
		fmt.Fprintf(w, "One point awarded to " + CurEntry.User)
	} else {
		fmt.Fprintf(w, "No entry has been dequed.")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("TOTAL POINTS")
		Users.Show()
		os.Exit(1)
	}()


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
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/share", shareHandler)
	http.HandleFunc("/deque", dequeHandler)
	http.HandleFunc("/currentEntry", currentEntryHandler)
	http.HandleFunc("/brownie", brownieHandler)
	err = http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
