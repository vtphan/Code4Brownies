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
func check_boardHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	board, ok := Boards[uid]
	if !ok {
		fmt.Println("check_board: unknown user", uid)
	} else {
		t := template.New("check board template")
		t, err := t.Parse(CHECK_BOARD_TEMPLATE)
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
func queue_lengthHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("queue length template")
	t, err := t.Parse(VIEW_SUBMISSION_QUEUE_TEMPLATE)
	if err == nil {
		data := struct{ Count int }{len(NewSubs)}
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
