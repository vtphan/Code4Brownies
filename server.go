//
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
	"syscall"
	"time"
	"errors"
)

var PORT = "4030"



//-----------------------------------------------------------------
// USERS
//-----------------------------------------------------------------

type User struct {
	name   string
	points int
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
	fmt.Fprintf(w, user+" has "+strconv.Itoa(points)+" brownies.")
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

func verifyPasscode(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized access")
	}
	return nil
}
//-----------------------------------------------------------------
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		Points.addOne(r.FormValue("uid"))
		fmt.Println("+1", r.FormValue("uid"))
		fmt.Fprintf(w, "Point awarded to "+r.FormValue("uid"))
	}
}

//-----------------------------------------------------------------
// return all current posts
//-----------------------------------------------------------------
func postsHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
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
	if verifyPasscode(w, r) == nil {
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

	// Create backup file
	backupFile, err := os.Create(csvFile + ".bak")
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
		points, err = strconv.Atoi(record[2])
		user := &User{record[1], points}
		AllUsers[record[0]] = user
		if err := writer.Write(record); err != nil {
			log.Fatalln("error backing up record to csv:", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		panic(err)
	}
	fmt.Println("Duplicated db to", csvFile+".bak")
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
// MAIN
//-----------------------------------------------------------------

var PASSCODE string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var userFilename string
	flag.StringVar(&PASSCODE, "passcode", "password", "passcode to be used by the instructor to connect to the server.")
	flag.StringVar(&userFilename, "db", "", "csv file with 3 fields: uid,name,points.")
	flag.Parse()

	loadRecords(userFilename)
	prepareCleanup(userFilename)
	informIPAddress()

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
