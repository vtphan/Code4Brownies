//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

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
