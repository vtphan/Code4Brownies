//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//-----------------------------------------------------------------
// STUDENT's HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// return brownie points a user has received.
//-----------------------------------------------------------------
func my_pointsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("uid")
	all_points, today_points := 0, 0
	var date time.Time
	today := time.Now().Day()
	rows, _ := database.Query("select points, date from submission where uid=?", user)
	defer rows.Close()
	for rows.Next() {
		p := 0
		rows.Scan(&p, &date)
		if date.Day() == today {
			today_points += p
		}
		all_points += p
	}
	rows2, _ := database.Query("select points, date from poll where uid=?", user)
	defer rows2.Close()
	for rows2.Next() {
		p := 0
		rows2.Scan(&p, &date)
		if date.Day() == today {
			today_points += p
		}
		all_points += p
	}
	str := "%s\nToday: %d brownie points.\nAll-time: %d brownie points.\n"
	mesg := fmt.Sprintf(str, user, today_points, all_points)
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
// students share their codes
//-----------------------------------------------------------------
func shareHandler(w http.ResponseWriter, r *http.Request) {
	uid, body, ext := r.FormValue("uid"), r.FormValue("body"), r.FormValue("ext")
	mode := r.FormValue("mode")
	bid := r.FormValue("bid")
	hints_used, err := strconv.Atoi(r.FormValue("hints_used"))
	if err != nil {
		fmt.Fprint(w, "Error converting number of hints used.")
		return
	}

	if mode == "code" {
		AddSubmission(uid, bid, body, ext, hints_used)
		fmt.Println(uid, "submitted.")
		fmt.Fprintf(w, uid+", thank you for sharing.")
	} else if mode == "poll" {
		body = strings.ToLower(body)
		prev_answer, ok := POLL_RESULT[uid]
		if ok {
			POLL_COUNT[prev_answer]--
		}
		POLL_RESULT[uid] = body
		POLL_COUNT[body]++
		fmt.Fprintf(w, uid+", thank you for voting.")
	} else if mode == "ask" {
		Questions = append(Questions, body)
		fmt.Fprint(w, "Your question will be addressed soon.")
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
		// js, err = json.Marshal(map[string]string{
		// 	"content":      board.Content,
		// 	"help_content": board.HelpContent,
		// 	"ext":          board.Ext,
		// 	"bid":          board.Bid,
		// })
		js, err = json.Marshal(board)
		Boards[uid] = []*Board{}
	}
	if err != nil {
		fmt.Println(err.Error())
		js, err = json.Marshal([]*Board{})
		// js, err = json.Marshal(map[string]string{"content": "", "ext": ""})
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
