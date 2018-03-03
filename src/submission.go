//
// Author: Vinhthuy Phan, 2015 - 2018
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
		InsertSubmissionSQL.Exec(sid, uid, bid, 0, des, ext, time.Now(), body, hints_used, nil)
	}
}

// ------------------------------------------------------------------
// Remove a submission from NewSubs
// ------------------------------------------------------------------
func RemoveSubmissionBySID(sid string) bool {
	SUBS_SEM.Lock()
	defer SUBS_SEM.Unlock()
	for i := 0; i < len(NewSubs); i++ {
		if NewSubs[i].Sid == sid {
			_, err := UpdateCompletionTimeSQL.Exec(time.Now(), NewSubs[i].Sid)
			if err != nil {
				fmt.Println("Failed to update completion time")
			}
			NewSubs = append(NewSubs[:i], NewSubs[i+1:]...)
			return true
		}
	}
	return false
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
		_, err := UpdateCompletionTimeSQL.Exec(time.Now(), s.Sid)
		if err != nil {
			fmt.Println("Failed to update completion time")
		}
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
