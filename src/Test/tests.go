package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//--------------------------------------------------------------
const URL = "http://141.225.10.144:4030/"
const N = 100

type handler func(chan<- int)

//--------------------------------------------------------------
func submit(done chan<- int) {
	handler_url := URL + "submit_post"
	id := strconv.Itoa(rand.Intn(100000))
	uid := "Tester" + id
	v := url.Values{}
	v.Set("uid", uid)
	v.Add("body", "print('hello world')")
	v.Add("ext", "py")
	resp, err := http.PostForm(handler_url, v)
	if err != nil {
		panic(err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(uid, ":", resp.Status, string(body))
	}
	done <- 1
}

//--------------------------------------------------------------
func whiteboard(done chan<- int) {
	handler_url := URL + "receive_broadcast"
	resp, err := http.Get(handler_url)
	if err != nil {
		panic(err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Status, string(body))
	}
	done <- 1
}

//--------------------------------------------------------------
func test_handler(h handler) {
	done := make(chan int, N+1)
	for i := 0; i < N; i++ {
		go h(done)
	}
	for i := 0; i < N; i++ {
		<-done
	}
	fmt.Println("Finish testing", N, "requests/posts.")
}

//--------------------------------------------------------------

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	test_handler(submit)
	test_handler(whiteboard)
}
