//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	// "path/filepath"
	"math/rand"
	"time"
)

var ADDR = ""
var PORT = "4030"
var PASSCODE string

//-----------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err.Error() + "\n")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Server address http://" + ipnet.IP.String() + ":" + PORT)
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

//-----------------------------------------------------------------
func main() {
	rand.Seed(time.Now().UnixNano())
	os.Mkdir("db", 0777)
	// USER_DB = filepath.Join(".", "db", "C4B_DB.csv")

	ADDR = informIPAddress()
	if ADDR == "" {
		panic("Unable to connect to the network.")
	}

	flag.StringVar(&PASSCODE, "passcode", ADDR, "passcode to be used by the instructor to connect to the server.")
	// flag.StringVar(&USER_DB, "db", USER_DB, "user database in csv format, which consists of UID,POINTS.")
	flag.Parse()

	// loadRecords()
	// prepareCleanup()

	// student handlers
	http.HandleFunc("/submit_post", submit_postHandler)   // rename this
	http.HandleFunc("/my_points", my_pointsHandler)
	http.HandleFunc("/receive_broadcast", receive_broadcastHandler)

	// teacher handlers
	http.HandleFunc("/points", pointsHandler)
	http.HandleFunc("/give_point", give_pointHandler)
	http.HandleFunc("/peek", peekHandler)
	http.HandleFunc("/broadcast", broadcastHandler)
	http.HandleFunc("/get_post", get_postHandler)
	http.HandleFunc("/get_posts", get_postsHandler)
	err := http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
