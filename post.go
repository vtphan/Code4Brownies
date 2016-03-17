//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"sync"
	"os/exec"
	"fmt"
)

//-----------------------------------------------------------------
// POSTS
//-----------------------------------------------------------------

type Post struct {
	Uid  string
	Body string
	Ext string
}

type PostQueue struct {
	queue []*Post
	sem   sync.Mutex
}

var Posts = PostQueue{}                         // posts of currently active users

func (P *PostQueue) Add(uid, body, ext string) {
	P.sem.Lock()
	P.queue = append(P.queue, &Post{uid, body, ext})
	if len(P.queue) == 1 {
		_, err := exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Output()
		if err != nil {
			fmt.Println(err)
		}
	}
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

func (P *PostQueue) Clear() {
	P.sem.Lock()
	P.queue = nil
	P.sem.Unlock()
}


func (P *PostQueue) Get(i int) *Post {
	P.sem.Lock()
	defer P.sem.Unlock()
	return P.queue[i]
}

