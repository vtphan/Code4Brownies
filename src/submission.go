//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"fmt"
	"strings"
	"time"
)

// ------------------------------------------------------------------
func AddSubmission(uid, bid, body, ext string, hints_used int) {
	SUBS_SEM.Lock()
	defer SUBS_SEM.Unlock()
	_, ok := Boards[uid]
	if ok {
		des := ""
		if strings.HasPrefix(body, "#") || strings.HasPrefix(body, "//") {
			des = strings.SplitN(body, "\n", 2)[0]
		}
		sid := RandStringRunes(6)
		timestamp := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
		sub := &Submission{
			Sid:       sid,
			Bid:       bid,
			Uid:       uid,
			Body:      body,
			Ext:       ext,
			Points:    0,
			Pdes:      des,
			Timestamp: timestamp,
			// Duration:  dur,
		}
		AllSubs[sid] = sub
		NewSubs = append(NewSubs, sub)
		if len(NewSubs) == 1 {
			fmt.Print("\x07")
		}
		InsertSubmissionSQL.Exec(sid, uid, bid, 0, des, ext, time.Now(), body, hints_used)
	}
}

// ------------------------------------------------------------------
// Remove a submission from NewSubs
// ------------------------------------------------------------------
func RemoveSubmissionBySID(sid string) {
	SUBS_SEM.Lock()
	defer SUBS_SEM.Unlock()
	for i := 0; i < len(NewSubs); i++ {
		if NewSubs[i].Sid == sid {
			NewSubs = append(NewSubs[:i], NewSubs[i+1:]...)
			return
		}
	}
}

// ------------------------------------------------------------------
// Remove a submission from NewSubs
// ------------------------------------------------------------------
func RemoveSubmission(i int) *Submission {
	if i < 0 || len(NewSubs) == 0 || i > len(NewSubs) {
		return &Submission{}
	} else {
		SUBS_SEM.Lock()
		defer SUBS_SEM.Unlock()
		s := NewSubs[i]
		NewSubs = append(NewSubs[:i], NewSubs[i+1:]...)
		return s
	}
}

// ------------------------------------------------------------------
func ProcessPollResult(uid string, is_correct int) {
	brownies := 1
	if is_correct == 1 {
		brownies = 2
	}
	InsertPollSQL.Exec(uid, is_correct, brownies, time.Now())
}

// ------------------------------------------------------------------
