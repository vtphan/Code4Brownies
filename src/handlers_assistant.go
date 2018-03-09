//
// Author: Vinhthuy Phan, 2015 - 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

//-----------------------------------------------------------------
// TA gives brownie points to a user
//-----------------------------------------------------------------
func ta_give_pointsHandler(w http.ResponseWriter, r *http.Request) {
	if sub, ok := AllSubs[r.FormValue("sid")]; ok {
		points, err := strconv.Atoi(r.FormValue("points"))
		if err != nil {
			fmt.Fprint(w, "Failed")
		} else {
			success := RemoveSubmissionBySID(r.FormValue("sid"))
			if success == false {
				// if instructor graded this submission, ignore TA.
				fmt.Fprintf(w, "This submission is already graded.")
			} else {
				sub.Points = points
				_, err = UpdatePointsSQL.Exec(sub.Points, r.FormValue("sid"))
				if err != nil {
					fmt.Fprint(w, "Failed")
				} else {
					mesg := fmt.Sprintf("%s: %d points.\n", sub.Uid, sub.Points)
					fmt.Fprintf(w, mesg)
				}
			}
		}
	}

}
