//
// Author: Vinhthuy Phan, 2015 - 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

//-----------------------------------------------------------------
func track_boardHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	board, ok := Boards[uid]
	if !ok {
		fmt.Println("track_board: unknown user", uid)
	} else {
		t := template.New("check board template")
		t, err := t.Parse(TRACK_BOARD_TEMPLATE)
		if err == nil {
			data := struct{ Message string }{""}
			if len(board) > 0 {
				data = struct{ Message string }{YOU_GOT_CODE}
			}
			w.Header().Set("Content-Type", "text/html")
			t.Execute(w, data)
		} else {
			fmt.Println(err)
		}
	}
}

//-----------------------------------------------------------------
func track_submissionsHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("track submission template")
	t, err := t.Parse(TRACK_SUBMISSIONS_TEMPLATE)
	if err == nil {
		data := struct{ Message string }{""}
		if len(NewSubs) > 0 {
			data = struct{ Message string }{fmt.Sprintf("%d", len(NewSubs))}
		}
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, data)
	} else {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------
func view_questionsHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("questions template")
	t, err := t.Parse(QUESTIONS_TEMPLATE)
	if err == nil {
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, &QuestionsData{Questions: Questions})
	} else {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------
func get_questionsHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(Questions)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
