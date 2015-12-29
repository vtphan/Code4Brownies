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
