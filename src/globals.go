//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main
import (
	"time"
	"sync"
)
var ADDR = ""
var PORT = "4030"
var USER_DB string
var SERVER = ""

var Whiteboard string
var WhiteboardExt string

var ProblemStartingTime time.Time
var ProblemDescription string
var ProblemID string

type Submission struct {
	Sid      string // submission id
	Uid      string // user id
	Pid      string // problem id
	Body     string
	Ext      string
	Points   int
	Duration int // in seconds
	Pdes 		string // problem description
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var SEM sync.Mutex
var NewSubs = make([]*Submission, 0)
var ProcessedSubs = make(map[string]*Submission)

type TemplateData struct {
	SERVER string
}
var POLL_MODE = false
var POLL_RESULT = make(map[string]int)
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