package main

import (
	"fmt"
	"sync"
)

type UserType struct {
	Points map[string]int
	m sync.Mutex
}

func NewUsers() *UserType {
	U := &UserType{}
	U.Points = make(map[string]int)
	return U
}

func (U *UserType) OnePoint(usr string) {
	score, ok := U.Points[usr]
	U.m.Lock()
	if !ok {
		U.Points[usr] = 1
	} else {
		U.Points[usr] = score + 1
	}
	U.m.Unlock()
}

func (U *UserType) GetPoints(usr string) int {
	_, ok := U.Points[usr]
	if !ok {
		U.Points[usr] = 0
	}
	return U.Points[usr]
}

func (U *UserType) Show() {
	for key,value := range U.Points {
		fmt.Println(key,"\t",value)
	}
}
