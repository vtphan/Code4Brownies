//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var Whiteboard string
var WhiteboardExt string
var Problems = make(map[string]time.Time)

//-----------------------------------------------------------------
// Student's handlers
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// users query to know their current points
//-----------------------------------------------------------------
func my_pointsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("uid")
	total := 0
	for _, s := range ProcessedSubs {
		if user == s.Uid {
			total += s.Points
		}
	}
	mesg := fmt.Sprintf("%d points for %s\n", total, user)
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
// users submit their codes
//-----------------------------------------------------------------
func submit_postHandler(w http.ResponseWriter, r *http.Request) {
	uid, body, ext := r.FormValue("uid"), r.FormValue("body"), r.FormValue("ext")
	AddSubmission(uid, body, ext)
	fmt.Println(uid, "submitted.")
	fmt.Fprintf(w, uid+" submitted succesfully.")
	// PrintState()
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
// Instructor's handlers
//-----------------------------------------------------------------

func authorize(w http.ResponseWriter, r *http.Request) error {
	if r.Host != "localhost:4030" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized access")
	}
	return nil
}

//-----------------------------------------------------------------
// instructor broadcast contents to students
//-----------------------------------------------------------------
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		Whiteboard = r.FormValue("content")
		WhiteboardExt = r.FormValue("ext")
		problem_id := get_problem_id(Whiteboard)
		Problems[problem_id] = time.Now()
		fmt.Fprintf(w, "Content is saved to whiteboard.")
	}
}

//-----------------------------------------------------------------
// return points of currently awarded users
//-----------------------------------------------------------------
func pointsHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
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
}

//-----------------------------------------------------------------
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		s := GetSubmission(r.FormValue("sid"))
		if s != nil {
			s.Points++
			fmt.Fprintf(w, "Point awarded to "+s.Uid)
		} else {
			fmt.Fprintf(w, "No submission is associated with this file.")
		}
		// PrintState()
	}
}

//-----------------------------------------------------------------
// return all current NewSubs
//-----------------------------------------------------------------
func peekHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		js, err := json.Marshal(NewSubs)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves code
//-----------------------------------------------------------------
func get_postHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
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
}

//-----------------------------------------------------------------
// Instructor retrieves all codes
//-----------------------------------------------------------------
func get_postsHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
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
		// PrintState()
	}
}
