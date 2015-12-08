//
// Live Coding (server)
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var PORT = "4030"

//-----------------------------------------------------------------
// POSTS
//-----------------------------------------------------------------

type Post struct {
	Uid  string
	Body string
}

type PostQueue struct {
	queue []*Post
	sem   sync.Mutex
}

func (P *PostQueue) Add(uid, body string) {
	P.sem.Lock()
	P.queue = append(P.queue, &Post{uid, body})
	P.sem.Unlock()
}

func (P *PostQueue) Remove(i int) *Post {
	if i < 0 || len(P.queue) == 0 || i > len(P.queue) {
		return &Post{}
	} else {
		P.sem.Lock()
		post := P.queue[i]
		P.queue = append(P.queue[:i], P.queue[i+1:]...)
		P.sem.Unlock()
		return post
	}
}

func (P *PostQueue) Get(i int) *Post {
	P.sem.Lock()
	defer P.sem.Unlock()
	return P.queue[i]
}

//-----------------------------------------------------------------
// POINTS
//-----------------------------------------------------------------

type Point struct {
	data map[string]int // maps uids to brownie points
	sem  sync.Mutex
}

func (P *Point) addOne(usr string) {
	P.sem.Lock()
	_, ok := P.data[usr]
	if !ok {
		P.data[usr] = 0
	}
	P.data[usr] += 1
	P.sem.Unlock()
}

func (P *Point) get(usr string) int {
	P.sem.Lock()
	defer P.sem.Unlock()

	_, ok := P.data[usr]
	if !ok {
		P.data[usr] = 0
	}
	return P.data[usr]
}

//-----------------------------------------------------------------
// USERS
//-----------------------------------------------------------------

type User struct {
	name   string
	points int
}

//-----------------------------------------------------------------
// load records in csv file into global AllUsers
//-----------------------------------------------------------------
func loadRecords(csvFile string) {
	if csvFile == "" {
		flag.PrintDefaults()
		log.Fatal("Must give file containg user records.")
	}
	userFile, err := os.Open(csvFile)
	defer userFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(userFile)
	var points int

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		points, err = strconv.Atoi(record[2])
		user := &User{record[1], points}
		AllUsers[record[0]] = user
	}

	for k, v := range AllUsers {
		fmt.Println(k, v)
	}
}

//-----------------------------------------------------------------
// HTTP HANDLERS
//-----------------------------------------------------------------

func my_pointsHandler(w http.ResponseWriter, r *http.Request) {
	_, user := r.FormValue("login"), r.FormValue("username")
	points, ok := Points.data[user]
	if !ok {
		points = 0
	}
	fmt.Fprintf(w, user+" has "+strconv.Itoa(points)+" points.")
}

//-----------------------------------------------------------------
// users submit their codes
//-----------------------------------------------------------------
func submit_postHandler(w http.ResponseWriter, r *http.Request) {
	_, user, body := r.FormValue("login"), r.FormValue("uid"), r.FormValue("body")
	if _, ok := AllUsers[user]; ok {
		Posts.Add(user, body)
		fmt.Println(user, "submitted.")
		fmt.Fprintf(w, "1")
	} else {
		fmt.Println(user, "non existent.")
		fmt.Fprintf(w, "0")
	}
}

//-----------------------------------------------------------------
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		Points.addOne(r.FormValue("uid"))
		fmt.Println("+1", r.FormValue("uid"))
		fmt.Fprintf(w, "Point awarded to "+r.FormValue("uid"))
	}
}

//-----------------------------------------------------------------
// return all current posts
//-----------------------------------------------------------------
func postsHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		js, err := json.Marshal(Posts.queue)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves code
//-----------------------------------------------------------------
func get_postHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		e, err := strconv.Atoi(r.FormValue("post"))
		if err != nil {
			fmt.Println(err.Error)
		} else {
			js, err := json.Marshal(Posts.Remove(e))
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
// var Points = &Point{}                 // points of currently active users
// var AllUsers = make(map[string]*User) // maps uids to users
//-----------------------------------------------------------------

func prepareCleanup(userFile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		t := time.Now()
		fmt.Println("\n" + t.Format("Mon Jan 2 15:04:05 MST 2006"))

		fmt.Println("Updating points for all users")
		for uid, p := range Points.data {
			fmt.Println(uid, "\t", p)
			AllUsers[uid].points += p
		}

		outFile, err := os.Create(userFile)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()

		fmt.Println("Writing data to", userFile)
		w := csv.NewWriter(outFile)
		for uid, user := range AllUsers {
			record := []string{uid, user.name, strconv.Itoa(user.points)}
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
// GLOBALS
//-----------------------------------------------------------------

var Posts = PostQueue{}                         // posts of currently active users
var Points = &Point{data: make(map[string]int)} // points of currently active users
var AllUsers = make(map[string]*User)           // maps uids to users

//-----------------------------------------------------------------
// MAIN
//-----------------------------------------------------------------

var PASSCODE string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var userFilename string
	flag.StringVar(&PASSCODE, "passcode", "password", "passcode to be used by the instructor to connect to the server.")
	flag.StringVar(&userFilename, "users", "", "csv-formatted filename containing usernames,real names.")
	flag.Parse()

	loadRecords(userFilename)
	prepareCleanup(userFilename)
	informIPAddress()

	// Register handlers and start serving app
	http.HandleFunc("/submit_post", submit_postHandler)
	http.HandleFunc("/my_points", my_pointsHandler)
	http.HandleFunc("/give_point", give_pointHandler)
	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/get_post", get_postHandler)
	err := http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
