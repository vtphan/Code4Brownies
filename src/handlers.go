//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)


//-----------------------------------------------------------------
// STUDENT's HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// Query poll results
//-----------------------------------------------------------------
func query_pollHandler(w http.ResponseWriter, r *http.Request) {
	if POLL_MODE {
		js, err := json.Marshal(POLL_RESULT)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(js)
		}
	}
}

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
// users submit their codes
//-----------------------------------------------------------------
func submit_postHandler(w http.ResponseWriter, r *http.Request) {
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
// students receive broadcast
//-----------------------------------------------------------------
func receive_broadcastHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(map[string]string{"whiteboard": Whiteboard, "ext": WhiteboardExt})
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
			t.Execute(w, &Data{SERVER})
		} else {
			fmt.Println(err)
		}
		// fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "There is no on-going poll.")
	}
}


//-----------------------------------------------------------------
// INSTRUCTOR's HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// Collect poll answers from students
//-----------------------------------------------------------------
func start_pollHandler(w http.ResponseWriter, r *http.Request) {
	POLL_MODE = !POLL_MODE
	if !POLL_MODE {
		POLL_RESULT = make(map[string]int)
	}
	fmt.Fprint(w, POLL_MODE)
}

//-----------------------------------------------------------------
// instructor broadcast contents to students
//-----------------------------------------------------------------
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	Whiteboard = r.FormValue("content")
	WhiteboardExt = r.FormValue("ext")
	fmt.Fprintf(w, "Content is saved to whiteboard.")
}

//-----------------------------------------------------------------
// instructor broadcast contents to students
//-----------------------------------------------------------------
func new_problemHandler(w http.ResponseWriter, r *http.Request) {
	ProblemStartingTime = time.Now()
	ProblemDescription = r.FormValue("description")
	ProblemID = RandStringRunes(8)
	fmt.Fprintf(w, "Starting a new problem. Clock restarted.")
}

//-----------------------------------------------------------------
// return points of all users
//-----------------------------------------------------------------
func pointsHandler(w http.ResponseWriter, r *http.Request) {
	subs := loadDB()
	for k, v := range ProcessedSubs {
		subs[k] = v
	}
	js, err := json.Marshal(subs)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointsHandler(w http.ResponseWriter, r *http.Request) {
	sub := GetSubmission(r.FormValue("sid"))
	if sub != nil {
		points, err := strconv.Atoi(r.FormValue("points"))
		if err != nil {
			fmt.Println(err)
		} else {
			sub.Points = points
			mesg := fmt.Sprintf("%s: %d points.\n", sub.Uid, points)
			fmt.Fprintf(w, mesg)
		}
	}
}

//-----------------------------------------------------------------
// return all current NewSubs
//-----------------------------------------------------------------
func peekHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(NewSubs)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
// Instructor retrieves a new submission
//-----------------------------------------------------------------
func get_postHandler(w http.ResponseWriter, r *http.Request) {
	e, err := strconv.Atoi(r.FormValue("post"))
	if err != nil {
		fmt.Println(err.Error)
	} else {
		js, err := json.Marshal(ProcessSubmission(e))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves all new submissions
//-----------------------------------------------------------------
func get_postsHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(NewSubs)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for len(NewSubs) > 0 {
			ProcessSubmission(0)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}


