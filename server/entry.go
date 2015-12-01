package main

import (
	"fmt"
	"sync"
)

type Entry struct {
	User string
	Body string
	N int
}

var CurEntry *Entry

type EntryList struct {
	list []*Entry
	m sync.Mutex
}

func NewEntryList() *EntryList {
	return &EntryList{}
}

func (E *EntryList) Len() int {
	return len(E.list)
}

func (E *EntryList) Add(user, body string) {
	E.m.Lock()
	E.list = append(E.list, &Entry{user,body,0})
	E.m.Unlock()
}

func (E *EntryList) Deque()  *Entry {
	if len(E.list) == 0 {
		return &Entry{}
	} else {
		E.m.Lock()
		CurEntry = E.list[0]
		E.list = E.list[1:]
		CurEntry.N = len(E.list)
		E.m.Unlock()
		return CurEntry
	}
}

func (E *EntryList) Show() {
	fmt.Println(">>>")
	for i:=0; i<len(E.list); i++ {
		fmt.Println(i, E.list[i])
	}
	fmt.Println("<<<")
}

