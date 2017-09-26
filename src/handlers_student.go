//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strings"
)

//-----------------------------------------------------------------
// STUDENT's HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// return brownie points a user has received.
//-----------------------------------------------------------------
func my_pointsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("uid")
	entries, points := 0, 0
	for _, s := range ProcessedSubs {
		if user == s.Uid {
			if s.Points > 0 {
				points += s.Points
				entries += 1
			}
		}
	}
	rows, _ := database.Query("select points from submission where uid=?", user)
	defer rows.Close()
	all_points := 0
	for rows.Next() {
		p := 0
		rows.Scan(&p)
		all_points += p
	}
	rows2, _ := database.Query("select points from poll where uid=?", user)
	defer rows2.Close()
	for rows2.Next() {
		p := 0
		rows2.Scan(&p)
		all_points += p
	}
	str := "%s\nToday: %d points, %d submissions.\nAll-time: %d points.\n"
	mesg := fmt.Sprintf(str, user, points, entries, all_points)
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
// students share their codes
//-----------------------------------------------------------------
func shareHandler(w http.ResponseWriter, r *http.Request) {
	uid, body, ext := r.FormValue("uid"), r.FormValue("body"), r.FormValue("ext")
	mode := r.FormValue("mode")
	bid := r.FormValue("bid")
	// PrintState()
	if mode == "code" {
		AddSubmission(uid, bid, body, ext)
		fmt.Println(uid, "submitted.")
		fmt.Fprintf(w, uid+", thank you for sharing.")
	} else if mode == "poll" {
		prev_answer, ok := POLL_RESULT[uid]
		if ok {
			POLL_COUNT[prev_answer]--
		}
		POLL_RESULT[uid] = body
		POLL_COUNT[body]++
		fmt.Fprintf(w, uid+", thank you for sharing.")
	} else {
		fmt.Fprint(w, "Unknown mode.")
	}
}

//-----------------------------------------------------------------
// student receives broadcast or content from his own board
//-----------------------------------------------------------------
func receive_broadcastHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	var js []byte
	var err error
	board, ok := Boards[uid]
	if ok {
		js, err = json.Marshal(map[string]string{
			"content": board.Content,
			"ext":     board.Ext,
			"bid":     board.Bid,
		})
	}
	if err != nil {
		fmt.Println(err.Error())
		js, err = json.Marshal(map[string]string{"content": "", "ext": ""})
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		board.Changed = false
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
// student checks to see if there is something new on his/her board
//-----------------------------------------------------------------
func check_broadcastHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	var err error
	board, ok := Boards[uid]
	if ok {
		fmt.Fprintf(w, "%t", board.Changed)
	} else {
		fmt.Println(err.Error())
	}
}
