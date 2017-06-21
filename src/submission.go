//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"fmt"
	"math/rand"
	"time"
)

//-----------------------------------------------------------------
// ProcessedSubs of students' submissions.
// Submitted asynchronously, submissions must be synchronized.
//-----------------------------------------------------------------

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
	SEM.Lock()
	defer SEM.Unlock()
	dur := int(time.Since(ProblemStartingTime).Seconds())
	pid := ProblemID
	des := ProblemDescription
	timestamp := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
	NewSubs = append(NewSubs, &Submission{RandStringRunes(10), uid, pid, body, ext, 0, dur, des, timestamp})
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
		SEM.Lock()
		defer SEM.Unlock()
		s := NewSubs[i]
		NewSubs = append(NewSubs[:i], NewSubs[i+1:]...)
		ProcessedSubs[s.Sid] = s
		return s
	}
}

// ------------------------------------------------------------------
func PrintState() {
	fmt.Println("------\n\tNewSubs:")
	for _, s := range NewSubs {
		fmt.Printf("Sid: %s\nUid: %s\nPid: %s\nExt: %s\nBody length: %d\nPoints: %d\nDuration: %d\n\n",
			s.Sid, s.Uid, s.Pid, s.Ext, len(s.Body), s.Points, s.Duration)
	}
	fmt.Println("\n\tProcessedSubs:")
	for _, s := range ProcessedSubs {
		fmt.Printf("Sid: %s\nUid: %s\nPid: %s\nExt: %s\nBody length: %d\nPoints: %d\nDuration: %d\n\n",
			s.Sid, s.Uid, s.Pid, s.Ext, len(s.Body), s.Points, s.Duration)
	}
	fmt.Println("------")
}
