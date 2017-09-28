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
func AddSubmission(uid, bid, body, ext string) {
	SEM.Lock()
	defer SEM.Unlock()
	board, ok := Boards[uid]
	if ok {
		dur := int(time.Since(board.StartingTime).Seconds())
		des := strings.SplitN(body, "\n", 2)[0]
		if des != board.Description {
			des = ""
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
			Duration:  dur,
			Pdes:      des,
			Timestamp: timestamp,
		}
		AllSubs[sid] = sub
		NewSubs = append(NewSubs, sub)
		if len(NewSubs) == 1 {
			fmt.Print("\x07")
		}
		InsertSubmissionSQL.Exec(sid, uid, bid, 0, dur, des, ext, time.Now(), body)
	}
}

// ------------------------------------------------------------------
// Remove a submission from NewSubs
func RemoveSubmission(i int) *Submission {
	if i < 0 || len(NewSubs) == 0 || i > len(NewSubs) {
		return &Submission{}
	} else {
		SEM.Lock()
		defer SEM.Unlock()
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
