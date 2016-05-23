//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	// "strconv"
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
	// user := r.FormValue("uid")
	// _, ok := AllUsers[user]
	// if !ok {
	// 	AllUsers[user] = &User{0}
	// }
	// record := AllUsers[user]
	// cur_points := Points.get(user)
	// mesg := fmt.Sprintf("Points for %s\nCurrent: %d\nTotal: %d\n", user, cur_points, record.points)
	// fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
// users submit their codes
//-----------------------------------------------------------------
func submit_postHandler(w http.ResponseWriter, r *http.Request) {
	uid, body, ext := r.FormValue("uid"), r.FormValue("body"), r.FormValue("ext")
	AddSubmission(uid, body, ext)
	fmt.Println(uid, "submitted.")
	fmt.Fprintf(w, uid+" submitted succesfully.")
	PrintState()
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

func verifyPasscode(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized access")
	}
	return nil
}

//-----------------------------------------------------------------
// instructor broadcast contents to students
//-----------------------------------------------------------------
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		Whiteboard = r.FormValue("content")
		WhiteboardExt = r.FormValue("ext")
		problem_id := get_problem_id(Whiteboard)
		Problems[problem_id] = time.Now()
		fmt.Println(Problems)
		fmt.Fprintf(w, "Content is saved to whiteboard.")
	}
}

//-----------------------------------------------------------------
// return points of currently awarded users
//-----------------------------------------------------------------
// func pointsHandler(w http.ResponseWriter, r *http.Request) {
// 	if verifyPasscode(w, r) == nil {
// 		js, err := json.Marshal(Points.data)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		} else {
// 			w.Header().Set("Content-Type", "application/json")
// 			w.Write(js)
// 		}
// 	}
// }

//-----------------------------------------------------------------
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		s := GetSubmission(r.FormValue("sid"))
		if s != nil {
			s.Points++
			fmt.Fprintf(w, "Point awarded to " + s.Uid)
		} else {
			fmt.Fprintf(w, "Not found.")
		}

		PrintState()
	}
}

//-----------------------------------------------------------------
// return all current NewSubs
//-----------------------------------------------------------------
// func peekHandler(w http.ResponseWriter, r *http.Request) {
// 	if verifyPasscode(w, r) == nil {
// 		js, err := json.Marshal(NewSubs.queue)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		} else {
// 			w.Header().Set("Content-Type", "application/json")
// 			w.Write(js)
// 		}
// 	}
// }

//-----------------------------------------------------------------
// Instructor retrieves code
//-----------------------------------------------------------------
// func get_postHandler(w http.ResponseWriter, r *http.Request) {
// 	if verifyPasscode(w, r) == nil {
// 		e, err := strconv.Atoi(r.FormValue("post"))
// 		if err != nil {
// 			fmt.Println(err.Error)
// 		} else {
// 			js, err := json.Marshal(NewSubs.Remove(e))
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			} else {
// 				w.Header().Set("Content-Type", "application/json")
// 				w.Write(js)
// 			}
// 		}
// 	}
// }

//-----------------------------------------------------------------
// Instructor retrieves all codes
//-----------------------------------------------------------------
func get_postsHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
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
		PrintState()
	}
}
