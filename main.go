//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"runtime"
)

var PORT = "4030"
var PASSCODE string

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
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&PASSCODE, "passcode", "password", "passcode to be used by the instructor to connect to the server.")
	flag.Parse()

	loadRecords()
	prepareCleanup()
	informIPAddress()

	http.HandleFunc("/submit_post", submit_postHandler)
	http.HandleFunc("/my_points", my_pointsHandler)
	http.HandleFunc("/points", pointsHandler)
	http.HandleFunc("/give_point", give_pointHandler)
	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/get_post", get_postHandler)
	err := http.ListenAndServe("0.0.0.0:"+PORT, nil)
	if err != nil {
		panic(err.Error() + "\n")
	}
}
