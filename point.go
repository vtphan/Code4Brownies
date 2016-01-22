//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"sync"
)

//-----------------------------------------------------------------
// POINTS
//-----------------------------------------------------------------

type Point struct {
	data map[string]int // maps uids to brownie points
	sem  sync.Mutex
}

var Points = &Point{data: make(map[string]int)} // points of currently active users

func (P *Point) addOne(uid string) {
	P.sem.Lock()
	if _, ok := P.data[uid]; !ok {
		P.data[uid] = 0
	}
	P.data[uid] += 1

	if _, ok := AllUsers[uid]; !ok {
		AllUsers[uid] = &User{0}
	}
	AllUsers[uid].points += 1
	P.sem.Unlock()
}

func (P *Point) get(uid string) int {
	P.sem.Lock()
	defer P.sem.Unlock()

	_, ok := P.data[uid]
	if !ok {
		P.data[uid] = 0
	}
	return P.data[uid]
}
