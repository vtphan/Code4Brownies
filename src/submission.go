//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"sync"
	"fmt"
	"strings"
	"math/rand"
	"time"
)

//-----------------------------------------------------------------
// ProcessedSubs of students' submissions.
// Submitted asynchronously, submissions must be synchronized.
//-----------------------------------------------------------------

type Submission struct {
	Sid  string   // submission id
	Uid  string   // user id
	Pid  string   // problem id. Example:  # :: dynamic programming (scafolding)
	Body string
	Ext string
	Points int
	Duration int  // in seconds
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var sem sync.Mutex
var NewSubs = make([]*Submission, 0)
var ProcessedSubs = make(map[string]*Submission)

func get_problem_id(program string) string {
	things := strings.SplitN(program, "\n", 2)
	if len(things) > 0 && len(things[0]) > 2 {
		return strings.Replace(strings.Trim(things[0][2:], " "), ",", "", -1)
	}
	return "none"
}

// ------------------------------------------------------------------
func GetSubmission(sid string) *Submission {
	if _, ok := ProcessedSubs[sid]; ok {
		return ProcessedSubs[sid]
	}
	return nil
}

// ------------------------------------------------------------------
func AddSubmission(uid, body, ext string) {
	sem.Lock()
	defer sem.Unlock()
	pid := get_problem_id(body)
	duration := 0
	if _, ok := Problems[pid]; !ok {
		pid = "undefined"
	} else {
		duration = int(time.Since(Problems[pid]).Seconds())
	}
	NewSubs = append(NewSubs, &Submission{RandStringRunes(10),uid,pid,body,ext,0,duration})
	if len(NewSubs) == 1 {
		fmt.Print("\x07")
	}
}

// ------------------------------------------------------------------
// Remove from NewSubs and add to ProcessedSubs
func ProcessSubmission(i int) *Submission {
	if i < 0 || len(NewSubs) == 0 || i > len(NewSubs) {
		return &Submission{}
	} else {
		sem.Lock()
		defer sem.Unlock()
		s := NewSubs[i]
		NewSubs = append(NewSubs[:i], NewSubs[i+1:]...)
		ProcessedSubs[s.Sid] = s
		return s
	}
}

// ------------------------------------------------------------------
func PrintState() {
	fmt.Println("------\n\tNewSubs:")
	for _, s := range(NewSubs) {
		fmt.Printf("Sid: %s\nUid: %s\nPid: %s\nExt: %s\nBody length: %d\nPoints: %d\nDuration: %d\n\n",
			s.Sid, s.Uid, s.Pid, s.Ext, len(s.Body), s.Points, s.Duration )
	}
	fmt.Println("\n\tProcessedSubs:")
	for _, s := range(ProcessedSubs) {
		fmt.Printf("Sid: %s\nUid: %s\nPid: %s\nExt: %s\nBody length: %d\nPoints: %d\nDuration: %d\n\n",
			s.Sid, s.Uid, s.Pid, s.Ext, len(s.Body), s.Points, s.Duration )
	}
	fmt.Println("------")
}


