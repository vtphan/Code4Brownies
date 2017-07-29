//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
// students receive feedback
//-----------------------------------------------------------------
// func receive_feedbackHandler(w http.ResponseWriter, r *http.Request) {
// 	uid := r.FormValue("uid")
// 	content, ok := Feedback[uid]
// 	if !ok {
// 		content = ""
// 	}
// 	js, err := json.Marshal(map[string]string{"content": content})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	} else {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(js)
// 	}
// }

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
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
// PUBLIC HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// View poll results
//-----------------------------------------------------------------
func view_pollHandler(w http.ResponseWriter, r *http.Request) {
	if POLL_MODE {
		// tmpl, err := template.ParseFiles("poll.html")
		t := template.New("poll template")
		t, err := t.Parse(POLL_TEMPLATE)
		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			t.Execute(w, &TemplateData{SERVER})
		} else {
			fmt.Println(err)
		}
		// fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "There is no on-going poll.")
	}
}
