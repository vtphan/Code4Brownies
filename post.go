package main

import (
	"sync"
)

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

