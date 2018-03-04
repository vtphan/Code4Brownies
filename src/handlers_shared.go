//
// Author: Vinhthuy Phan, 2015 - 2018
//
// Handlers for both instructor and TAs. Although the authorization is
// done differently (main.go), the operations are identical.
//
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//-----------------------------------------------------------------
// Instructor/TAs give feedback to a student
//-----------------------------------------------------------------
func feedbackHandler(w http.ResponseWriter, r *http.Request, author string) {
	BOARDS_SEM.Lock()
	defer BOARDS_SEM.Unlock()
	content, ext, sid := r.FormValue("content"), r.FormValue("ext"), r.FormValue("sid")
	sub, ok := AllSubs[sid]
	if ok {
		_, err := InsertFeedbackSQL.Exec(author, sub.Uid, content, sid, time.Now())
		if err != nil {
			fmt.Println("Error inserting feedback.", err)
			fmt.Fprintf(w, "Error inserting feedback.")
		} else {
			bid := ""
			SelectBidFromSidSQL.QueryRow(sid).Scan(&bid)
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
		}
		// fmt.Printf("%s gave feedback\n", author)
		fmt.Fprintf(w, "Content feed %s's virtual board.", sub.Uid)
	} else {
		fmt.Fprintf(w, "sid %s is not found.", sid)
	}
}

//-----------------------------------------------------------------
// Instructor/TA retrieve all new submissions
//-----------------------------------------------------------------
func get_postsHandler(w http.ResponseWriter, r *http.Request, author string) {
	SUBS_SEM.Lock()
	defer SUBS_SEM.Unlock()

	js, err := json.Marshal(NewSubs)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		// NewSubs = make([]*Submission, 0)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
