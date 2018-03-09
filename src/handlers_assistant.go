//
// Author: Vinhthuy Phan, 2015 - 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	// "strings"
	// "time"
)

//-----------------------------------------------------------------
func ta_share_with_teacherHandler(w http.ResponseWriter, r *http.Request) {
	TABoard_SEM.Lock()
	defer TABoard_SEM.Unlock()
	content, ext := r.FormValue("content"), r.FormValue("ext")
	TABoardOut = append(TABoardOut, &Code{Content: content, Ext: ext})
	fmt.Fprintf(w, "Content shared with instructor.")
}

//-----------------------------------------------------------------
func ta_get_from_teacherHandler(w http.ResponseWriter, r *http.Request) {
	TABoard_SEM.Lock()
	defer TABoard_SEM.Unlock()
	js, err := json.Marshal(TABoardIn)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
