//
// Author: Vinhthuy Phan, 2015 - 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

//-----------------------------------------------------------------
func view_public_boardHandler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.FormValue("i"))
	if err != nil {
		i = 0
	}
	if i >= len(PublicBoard) {
		i = 0
	}
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	temp := template.New("public board template")
	t, _ := temp.Funcs(funcMap).Parse(PUBLIC_BOARD_TEMPLATE)
	if i >= 0 && i < len(PublicBoard) {
		idx := make([]string, 0)
		for j := 0; j < len(PublicBoard); j++ {
			if i == j {
				idx = append(idx, "active")
			} else {
				idx = append(idx, "")
			}
		}
		x := ""
		if r.Host == "localhost:4030" {
			x = "x"
		}
		data := struct {
			Content string
			Idx     []string
			X       string
			AltText string
		}{Content: PublicBoard[i].Content, Idx: idx, X: x, AltText: ""}
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, &data)
	} else {
		data := struct {
			Content string
			Idx     []string
			X       string
			AltText string
		}{Content: "", Idx: []string{}, X: "", AltText: "Reload"}
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, &data)
	}
}

//-----------------------------------------------------------------
// Students tracking their white boards
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
// Instructor and TAs tracking student submissions
//-----------------------------------------------------------------
func track_submissionsHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("track submission template")
	t, err := t.Parse(TRACK_SUBMISSIONS_TEMPLATE)
	if err == nil {
		var count1, count2 string
		count1 = fmt.Sprintf("%d", len(NewSubs))
		if r.FormValue("view") == "ta" {
			count2 = fmt.Sprintf("%d", len(TABoardIn))
		} else {
			count2 = fmt.Sprintf("%d", len(TABoardOut))
		}
		data := struct {
			Count1 string
			Count2 string
		}{count1, count2}
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, data)
	} else {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------
func view_questionsHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("questions template")
	t, err := temp.Parse(QUESTIONS_TEMPLATE)
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
