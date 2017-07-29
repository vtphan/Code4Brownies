//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//-----------------------------------------------------------------
// STUDENT's HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// users query to know their current points
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
	mesg := fmt.Sprintf("%s: %d entries, %d points.\n", user, entries, points)
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
// students register
//-----------------------------------------------------------------
func registerHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	SEM.Lock()
	defer SEM.Unlock()
	RegisterStudent(uid)
	fmt.Fprint(w, uid+" registered.")
}

//-----------------------------------------------------------------
// students share their codes
//-----------------------------------------------------------------
func shareHandler(w http.ResponseWriter, r *http.Request) {
	uid, body, ext := r.FormValue("uid"), r.FormValue("body"), r.FormValue("ext")

	if POLL_MODE {
		// fmt.Println(uid, body)
		lines := strings.Split(body, "\n")
		if len(lines) > 0 && len(strings.Trim(lines[0], " ")) > 0 {
			POLL_RESULT[strings.Trim(lines[0], " ")]++
			fmt.Fprintf(w, "poll submitted.")
		} else {
			fmt.Fprintf(w, "your poll was not submitted.")
		}
	} else {
		AddSubmission(uid, body, ext)
		fmt.Println(uid, "submitted.")
		fmt.Fprintf(w, uid+" submitted succesfully.")
		// PrintState()
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
		js, err = json.Marshal(map[string]string{"content": board.Content})
	}
	if err != nil {
		fmt.Println(err.Error())
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
