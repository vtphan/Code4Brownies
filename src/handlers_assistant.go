//
// Author: Vinhthuy Phan, 2015 - 2018
//
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//-----------------------------------------------------------------
// Deprecated
//-----------------------------------------------------------------
func deprecated_feedbackHandler(w http.ResponseWriter, r *http.Request) {
	BOARDS_SEM.Lock()
	defer BOARDS_SEM.Unlock()
	content, ext, sid := r.FormValue("content"), r.FormValue("ext"), r.FormValue("sid")
	bid := ""
	err := SelectBidFromSidSQL.QueryRow(sid).Scan(&bid)
	if err != nil {
		fmt.Println("Error retrieving bid with", sid)
	}
	if bid == "" {
		bid = "wb_" + RandStringRunes(6)
		_, err = InsertBroadcastSQL.Exec(bid, content, ext, time.Now(), 0, "TA")
		if err != nil {
			fmt.Println("Error inserting into broadcast table.", err)
		}
	}
	sub, ok := AllSubs[sid]
	if ok {
		des := strings.SplitN(content, "\n", 2)[0]
		b := &Board{
			Content:      content,
			HelpContent:  "",
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
	fmt.Fprintf(w, "Content copied student's virtual board.")
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
