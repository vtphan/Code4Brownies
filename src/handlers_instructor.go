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
// Clear whiteboards
//-----------------------------------------------------------------
func clear_whiteboardsHandler(w http.ResponseWriter, r *http.Request) {
	for uid, _ := range Boards {
		Boards[uid] = make([]*Board, 0)
	}
	fmt.Fprintf(w, "Whiteboards cleared.")
}

//-----------------------------------------------------------------
// Clear questions
//-----------------------------------------------------------------
func clear_questionsHandler(w http.ResponseWriter, r *http.Request) {
	Questions = []string{}
	fmt.Fprintf(w, "Questions cleared.")
}

//-----------------------------------------------------------------
// Query poll results
//-----------------------------------------------------------------
func query_pollHandler(w http.ResponseWriter, r *http.Request) {
	counts := make(map[string]int)
	for k, v := range POLL_COUNT {
		if v > 0 {
			counts[k] = v
		}
	}
	js, err := json.Marshal(counts)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
// Start a new poll
//-----------------------------------------------------------------
func start_pollHandler(w http.ResponseWriter, r *http.Request) {
	POLL_DESCRIPTION = r.FormValue("description")
	if POLL_DESCRIPTION == "" {
		fmt.Fprint(w, "Empty")
	} else {
		fmt.Fprint(w, "Ok")
		POLL_ON = true
		POLL_RESULT = make(map[string]string)
		POLL_COUNT = make(map[string]int)
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
		t.Execute(w, &PollData{Description: POLL_DESCRIPTION})
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
	POLL_ON = false
	fmt.Fprintf(w, "Complete poll.")
}

//-----------------------------------------------------------------
// instructor hands out quiz questions
//-----------------------------------------------------------------
func send_quiz_questionHandler(w http.ResponseWriter, r *http.Request) {
	bid := "qz_" + RandStringRunes(6)
	question, answer := r.FormValue("question"), r.FormValue("answer")
	content := question + "\n\nANSWER: "

	_, err := InsertQuizSQL.Exec(bid, question, answer, time.Now())

	if err != nil {
		fmt.Println("Error inserting into quiz table.", err)
	} else {
		BOARDS_SEM.Lock()
		defer BOARDS_SEM.Unlock()

		for uid, _ := range Boards {
			b := &Board{
				Content:      content,
				HelpContent:  answer,
				Ext:          "txt",
				Bid:          bid,
				Description:  "Quiz " + bid,
				StartingTime: time.Now(),
			}
			Boards[uid] = append(Boards[uid], b)
		}
	}
}

//-----------------------------------------------------------------
// instructor broadcast contents to students
//-----------------------------------------------------------------
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	var des string
	bid := "wb_" + RandStringRunes(6)
	content, ext := r.FormValue("content"), r.FormValue("ext")
	help_content := r.FormValue("help_content")
	hints, err := strconv.Atoi(r.FormValue("hints"))
	if err != nil {
		fmt.Fprint(w, "Error converting number of hints.")
		return
	}

	_, err = InsertBroadCastSQL.Exec(bid, content, ext, time.Now(), hints)
	if err != nil {
		fmt.Println("Error inserting into broadcast table.", err)
	}

	BOARDS_SEM.Lock()
	defer BOARDS_SEM.Unlock()

	if r.FormValue("sids") == "__all__" {
		// for _, board := range Boards {
		for uid, _ := range Boards {
			des = strings.SplitN(content, "\n", 2)[0]
			b := &Board{
				Content:      content,
				HelpContent:  help_content,
				Ext:          ext,
				Bid:          bid,
				Description:  des,
				StartingTime: time.Now(),
			}
			Boards[uid] = append(Boards[uid], b)
		}
	} else {
		sids := strings.Split(r.FormValue("sids"), ",")
		for i := 0; i < len(sids); i++ {
			sid := string(sids[i])
			sub, ok := AllSubs[sid]
			if ok {
				des = strings.SplitN(content, "\n", 2)[0]
				b := &Board{
					Content:      content,
					HelpContent:  help_content,
					Ext:          ext,
					Bid:          bid,
					Description:  des,
					StartingTime: time.Now(),
				}
				Boards[sub.Uid] = append(Boards[sub.Uid], b)
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
	SUBS_SEM.Lock()
	defer SUBS_SEM.Unlock()
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
	SUBS_SEM.Lock()
	defer SUBS_SEM.Unlock()

	js, err := json.Marshal(NewSubs)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		NewSubs = make([]*Submission, 0)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
