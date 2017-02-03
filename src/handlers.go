//
// Author: Vinhthuy Phan, 2015
//
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"html/template"
)

var Whiteboard string
var WhiteboardExt string
var Problems = make(map[string]time.Time)
var POLL_MODE = false
var POLL_RESULT = make(map[string]int)

type Data struct {
	SERVER string
}

//-----------------------------------------------------------------
// Student's handlers
//-----------------------------------------------------------------

//-----------------------------------------------------------------
// Query poll results
//-----------------------------------------------------------------
func query_pollHandler(w http.ResponseWriter, r *http.Request) {
	if POLL_MODE {
		js, err := json.Marshal(POLL_RESULT)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
		   w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// users query to know their current points
//-----------------------------------------------------------------
func my_pointsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("uid")
	total := 0
	for _, s := range ProcessedSubs {
		if user == s.Uid {
			total += s.Points
		}
	}
	mesg := fmt.Sprintf("%d points for %s\n", total, user)
	fmt.Fprintf(w, mesg)
}

//-----------------------------------------------------------------
// users submit their codes
//-----------------------------------------------------------------
func submit_postHandler(w http.ResponseWriter, r *http.Request) {
	uid, body, ext := r.FormValue("uid"), r.FormValue("body"), r.FormValue("ext")
	if POLL_MODE {
		// fmt.Println(uid, body)
		lines := strings.Split(body, "\n")
		if len(lines) > 0 && len(strings.Trim(lines[0], " ")) > 0 {
			POLL_RESULT[strings.Trim(lines[0], " ")]++
			fmt.Fprintf(w, "poll submitted.")
		} else {
			fmt.Fprintf(w, "your poll was not submitted.")
		}
	} else {
		AddSubmission(uid, body, ext)
		fmt.Println(uid, "submitted.")
		fmt.Fprintf(w, uid+" submitted succesfully.")
		// PrintState()
	}
}

//-----------------------------------------------------------------
// students receive broadcast
//-----------------------------------------------------------------
func receive_broadcastHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(map[string]string{"whiteboard": Whiteboard, "ext": WhiteboardExt})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------
// Instructor's handlers
//-----------------------------------------------------------------

func authorize(w http.ResponseWriter, r *http.Request) error {
	if r.Host != "localhost:4030" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized access")
	}
	return nil
}

//-----------------------------------------------------------------
// Collect poll answers from students
//-----------------------------------------------------------------
func start_pollHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		POLL_MODE = !POLL_MODE
		if !POLL_MODE {
			POLL_RESULT = make(map[string]int)
		}
		fmt.Fprint(w, POLL_MODE)
	} else {
		fmt.Fprint(w, "Unauthorized")
	}
}


//-----------------------------------------------------------------
// View poll results
//-----------------------------------------------------------------
func view_pollHandler(w http.ResponseWriter, r *http.Request) {
	if POLL_MODE {
		// tmpl, err := template.ParseFiles("poll.html")
		t := template.New("poll template")
		t, err := t.Parse(POLL_TEMPLATE)
		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			t.Execute(w, &Data{SERVER})
		} else {
			fmt.Println(err)
		}
		// fmt.Fprintf(w, "OK")
	} else {
		fmt.Fprintf(w, "There is no on-going poll.")
	}
}

//-----------------------------------------------------------------
// instructor broadcast contents to students
//-----------------------------------------------------------------
func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		Whiteboard = r.FormValue("content")
		WhiteboardExt = r.FormValue("ext")
		problem_id := get_problem_id(Whiteboard)
		Problems[problem_id] = time.Now()
		fmt.Fprintf(w, "Content is saved to whiteboard.")
	}
}

//-----------------------------------------------------------------
// return points of all users
//-----------------------------------------------------------------
func pointsHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		subs := loadDB()
		for k, v := range ProcessedSubs {
			subs[k] = v
		}
		js, err := json.Marshal(subs)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// give one brownie point to a user
//-----------------------------------------------------------------
func give_pointHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		sub := GetSubmission(r.FormValue("sid"))
		if sub != nil {
			stage := r.FormValue("stage")
			if stage == "1" {
				total := 0
				for _, s := range ProcessedSubs {
					if s.Uid == sub.Uid {
						total += s.Points
					}
				}
				fmt.Fprintf(w, fmt.Sprintf("%s (%d)", sub.Uid, total))
			} else if stage == "2" {
				sub.Points++
				fmt.Fprintf(w, "Point awarded to "+sub.Uid)
			}
		}
		// PrintState()
	}
}

//-----------------------------------------------------------------
// return all current NewSubs
//-----------------------------------------------------------------
func peekHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		js, err := json.Marshal(NewSubs)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves a new submission
//-----------------------------------------------------------------
func get_postHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		e, err := strconv.Atoi(r.FormValue("post"))
		if err != nil {
			fmt.Println(err.Error)
		} else {
			js, err := json.Marshal(ProcessSubmission(e))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves all new submissions
//-----------------------------------------------------------------
func get_postsHandler(w http.ResponseWriter, r *http.Request) {
	if authorize(w, r) == nil {
		js, err := json.Marshal(NewSubs)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			for len(NewSubs) > 0 {
				ProcessSubmission(0)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
		// PrintState()
	}
}

//-----------------------------------------------------------------

var POLL_TEMPLATE = `
<!DOCTYPE HTML>
<html>
<head>
	<script type="text/javascript">
	window.onload = function () {
		var updateInterval = 3000;
		var maxUpdateTime = 300000;
		var totalUpdateTime = 0;
		var dps = [];
		var poll_view_url = "http://{{.SERVER}}/query_poll";
		var get_data = function() {
			$.getJSON(poll_view_url, function( data ) {
				$.each( data, function( key, val ) {
					var new_key = true;
					for (var i=0; i< dps.length; i++){
						if (key == dps[i].label) {
							dps[i].y = val;
							new_key = false;
						}
					}
					if (new_key == true) {
						dps.push({"label": key, "y": val});
					}
				});
				totalUpdateTime += updateInterval;
				if (totalUpdateTime > maxUpdateTime) {
					clearInterval(updateInterval);
				}
			});
		}
		get_data();
		var chart = new CanvasJS.Chart("chartContainer",{
			theme: "theme2",
			axisX:{
			},
			axisY: {
				interval: 1,
				gridThickness: 0,
				title: ""
			},
			legend:{
				verticalAlign: "top",
				horizontalAlign: "centre",
				fontSize: 18
			},
			data : [{
				type: "bar",
				showInLegend: true,
				legendMarkerType: "none",
				legendText: "Poll result",
				indexLabel: "{y}",
				dataPoints: dps
			}]
		});
		chart.render();
		var updateChart = function () {
			get_data();
			chart.render();
			// console.log("updated", dps)
		};
		setInterval(updateChart, updateInterval);
	}
	</script>
	<script src="http://canvasjs.com/assets/script/canvasjs.min.js"></script>
  	<script src="http://code.jquery.com/jquery-3.1.1.min.js"></script>
</head>
<body>
	<div id="chartContainer" style="height:600px; width:100%;">
	</div>
</body>
</html>
`
