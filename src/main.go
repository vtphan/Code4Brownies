//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

//-----------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err.Error() + "\n")
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			return ipnet.IP.String() + ":" + PORT
		}
	}
	return ""
}

//-----------------------------------------------------------------
// Make sure teacher runs server on his own laptop.
//-----------------------------------------------------------------
func Authorize(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "localhost:4030" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
		} else {
			fn(w, r)
		}
	}
}

//-----------------------------------------------------------------
// Register automatically if a student is not yet registered.
//-----------------------------------------------------------------
func AutoRegister(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := Boards[r.FormValue("uid")]; !ok {
			RegisterStudent(r.FormValue("uid"))
		}
		fn(w, r)
	}
}

//-----------------------------------------------------------------
func main() {
	SERVER = informIPAddress()
	fmt.Println("*********************************************")
	fmt.Printf("*   Code4Brownies (%s)\n", VERSION)
	fmt.Printf("*   Server address: %s\n", SERVER)
	fmt.Println("*********************************************\n")
	rand.Seed(time.Now().UnixNano())
	USER_DB = filepath.Join(".", "c4b.db")
	flag.StringVar(&USER_DB, "db", USER_DB, "user database.")
	flag.Parse()

	// student handlers
	http.HandleFunc("/share", AutoRegister(shareHandler))
	http.HandleFunc("/my_points", AutoRegister(my_pointsHandler))
	http.HandleFunc("/receive_broadcast", AutoRegister(receive_broadcastHandler))
	http.HandleFunc("/check_broadcast", AutoRegister(check_broadcastHandler))

	// teacher handlers
	http.HandleFunc("/query_poll", Authorize(query_pollHandler))
	http.HandleFunc("/view_poll", Authorize(view_pollHandler))
	http.HandleFunc("/answer_poll", Authorize(answer_pollHandler))
	http.HandleFunc("/give_points", Authorize(give_pointsHandler))
	http.HandleFunc("/peek", Authorize(peekHandler))
	http.HandleFunc("/broadcast", Authorize(broadcastHandler))
	http.HandleFunc("/get_post", Authorize(get_postHandler))
	http.HandleFunc("/get_posts", Authorize(get_postsHandler))

	init_sqldb()
	loadWhiteboards()

	// Start serving app
	err := http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
