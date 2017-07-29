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
// INSTRUCTOR's HANDLERS
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
	var des string
	if r.FormValue("sids") == "__all__" {
		for _, board := range Boards {
			board.Content = r.FormValue("content")
			des = strings.SplitN(board.Content, "\n", 2)[0]
			if des != board.Description { // a new exercise/problem
				board.Description = des
				board.StartingTime = time.Now()
			}
		}
	} else {
		sids := strings.Split(r.FormValue("sids"), ",")
		for i := 0; i < len(sids); i++ {
			sid := string(sids[i])
			sub, ok := ProcessedSubs[sid]
			if ok {
				Boards[sub.Uid].Content = r.FormValue("content")
				des = strings.SplitN(Boards[sub.Uid].Content, "\n", 2)[0]
				if des != Boards[sub.Uid].Description { // a new exercise/problem
					Boards[sub.Uid].Description = des
					Boards[sub.Uid].StartingTime = time.Now()
				}
			} else {
				fmt.Fprintf(w, "sid "+sid+" is not found.")
				return
			}
		}
	}
	fmt.Fprintf(w, "Content copied.")
}

//-----------------------------------------------------------------
// instructor sends test data
//-----------------------------------------------------------------
// func test_dataHandler(w http.ResponseWriter, r *http.Request) {
// 	ProblemTestData = r.FormValue("content")
// 	fmt.Fprintf(w, "Test data received.")
// }

//-----------------------------------------------------------------
// instructor sends a signal to clear whiteboard
//-----------------------------------------------------------------
// func clear_boardHandler(w http.ResponseWriter, r *http.Request) {
// 	ProblemStartingTime = time.Now()
// 	ProblemDescription = "none"
// 	ProblemID = "none"
// 	Whiteboard = ""
// 	fmt.Fprintf(w, "Whiteboard is clear.")
// }

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
