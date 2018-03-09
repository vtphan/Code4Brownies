//
// Author: Vinhthuy Phan, 2015 - 2018
//
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	// "strconv"
	"strings"
	"time"
)

//-----------------------------------------------------------------
func testHandler(w http.ResponseWriter, r *http.Request) {
	var m []BroadcastData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(m); i++ {
		fmt.Println(i, m[i].Sids, m[i].Content)
	}
	fmt.Fprintf(w, "Ok")
}

//-----------------------------------------------------------------
// INSTRUCTOR's HANDLERS
//-----------------------------------------------------------------

//-----------------------------------------------------------------
func share_with_taHandler(w http.ResponseWriter, r *http.Request) {
	TABoard_SEM.Lock()
	defer TABoard_SEM.Unlock()
	content, ext := r.FormValue("content"), r.FormValue("ext")
	TABoardIn = append(TABoardIn, &Code{Content: content, Ext: ext})
	fmt.Fprintf(w, "Content shared with TA.")
}

//-----------------------------------------------------------------
func get_from_taHandler(w http.ResponseWriter, r *http.Request) {
	TABoard_SEM.Lock()
	defer TABoard_SEM.Unlock()
	js, err := json.Marshal(TABoardOut)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

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
// instructor broadcasts materials to students
//-----------------------------------------------------------------
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	BOARDS_SEM.Lock()
	defer BOARDS_SEM.Unlock()

	// Get the json data
	var data []BroadcastData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		fmt.Fprintf(w, "No content copied to white boards.")
		return
	}

	// Determine broadcast ids
	bids := make([]string, 0)
	for i := 0; i < len(data); i++ {
		bid := "wb_" + RandStringRunes(6)
		_, err = InsertBroadcastSQL.Exec(
			bid,
			data[i].Content,
			data[i].Ext,
			time.Now(),
			data[i].Hints,
			"Instructor",
		)
		if err != nil {
			fmt.Println("Error inserting into broadcast table.", err)
			return
		}
		bids = append(bids, bid)
	}

	// Determine which boards to insert content
	selected_uid := make([]string, 0)
	if data[0].Sids == "__all__" {
		for uid, _ := range Boards {
			selected_uid = append(selected_uid, uid)
		}
	} else {
		sids := strings.Split(data[0].Sids, ",")
		for i := 0; i < len(sids); i++ {
			sid := string(sids[i])
			sub, ok := AllSubs[sid]
			if ok {
				selected_uid = append(selected_uid, sub.Uid)
			}
		}
	}

	// Insert content into white boards
	var des string
	mode := data[0].Mode
	rand_idx := make([]int, 0)
	if mode == 2 {
		i := 0
		for j := 0; j < len(selected_uid); j++ {
			rand_idx = append(rand_idx, i)
			i = (i + 1) % len(data)
		}
		rand.Shuffle(len(rand_idx), func(i, j int) {
			rand_idx[i], rand_idx[j] = rand_idx[j], rand_idx[i]
		})
	}
	for j := 0; j < len(selected_uid); j++ {
		uid := selected_uid[j]
		if mode < 2 {
			for i := 0; i < len(data); i++ {
				des = strings.SplitN(data[i].Content, "\n", 2)[0]
				b := &Board{
					Content:      data[i].Content,
					HelpContent:  data[i].Help_content,
					Ext:          data[i].Ext,
					Bid:          bids[i],
					Description:  des,
					StartingTime: time.Now(),
				}
				Boards[uid] = append(Boards[uid], b)
			}
		} else {
			i := rand_idx[j]
			des = strings.SplitN(data[i].Content, "\n", 2)[0]
			b := &Board{
				Content:      data[i].Content,
				HelpContent:  data[i].Help_content,
				Ext:          data[i].Ext,
				Bid:          bids[i],
				Description:  des,
				StartingTime: time.Now(),
			}
			Boards[uid] = append(Boards[uid], b)
		}
	}
	fmt.Fprintf(w, "Content copied to white boards.")
}
