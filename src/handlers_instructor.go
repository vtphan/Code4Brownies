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
// INSTRUCTOR's HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// Query poll results
//-----------------------------------------------------------------
func query_pollHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(POLL_COUNT)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
// View poll results
//-----------------------------------------------------------------
func view_pollHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("poll template")
	t, err := t.Parse(POLL_TEMPLATE)
	if err == nil {
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, &TemplateData{SERVER})
	} else {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------
// Answer poll
//-----------------------------------------------------------------
func answer_pollHandler(w http.ResponseWriter, r *http.Request) {
	answer := strings.ToLower(r.FormValue("answer"))
	for k, v := range POLL_RESULT {
		if v == answer {
			ProcessPollResult(k, 1)
		} else {
			ProcessPollResult(k, 0)
		}
	}
	POLL_RESULT = make(map[string]string)
	POLL_COUNT = make(map[string]int)
	fmt.Fprintf(w, "Complete poll.")
}

//-----------------------------------------------------------------
// instructor broadcast contents to students
//-----------------------------------------------------------------
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	var des string
	bid := "wb_" + RandStringRunes(6)
	content, ext := r.FormValue("content"), r.FormValue("ext")
	_, err := InsertBroadCastSQL.Exec(bid, content, ext, time.Now())
	if err != nil {
		fmt.Println("Error inserting into broadcast table.", err)
	}
	if r.FormValue("sids") == "__all__" {
		for _, board := range Boards {
			board.Content = content
			board.Ext = ext
			board.Bid = bid
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
			sub, ok := AllSubs[sid]
			if ok {
				Boards[sub.Uid].Content = content
				Boards[sub.Uid].Ext = ext
				Boards[sub.Uid].Bid = bid
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
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointsHandler(w http.ResponseWriter, r *http.Request) {
	if sub, ok := AllSubs[r.FormValue("sid")]; ok {
		points, err := strconv.Atoi(r.FormValue("points"))
		if err != nil {
			fmt.Println(err)
		} else {
			sub.Points = points
			UpdatePointsSQL.Exec(points, r.FormValue("sid"))
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
		js, err := json.Marshal(RemoveSubmission(e))
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
		NewSubs = make([]*Submission, 0)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
