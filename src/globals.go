//
// Author: Vinhthuy Phan, 2015 - 2018
//
package main

import (
	"math/rand"
	"sync"
	"time"
)

const VERSION = "0.42"

var ADDR = ""
var PORT = "4030"
var USER_DB string
var SERVER = ""
var TA_DB string
var TA_INFO = make(map[string]string)

type Board struct {
	Content      string
	HelpContent  string
	Description  string
	StartingTime time.Time
	Ext          string
	Bid          string // id of current broadcast
}

var Boards = make(map[string][]*Board)

type BroadcastData struct {
	Content      string `json:"content"`
	Sids         string `json:"sids"`
	Ext          string `json:"ext"`
	Help_content string `json:"help_content"`
	Hints        int    `json:"hints"`
	Mode         int    `json:"mode"`
	// Original_sid string `json:"original_sid"`
}

var Questions []string

type Submission struct {
	Sid       string // submission id
	Bid       string // broadcast id
	Uid       string // user id
	Body      string
	Ext       string
	Points    int
	Pdes      string // problem description
	Timestamp string
	// Duration  int    // in seconds
}

// ------------------------------------------------------------------
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// ------------------------------------------------------------------

var SUBS_SEM sync.Mutex
var BOARDS_SEM sync.Mutex
var NewSubs = make([]*Submission, 0)
var AllSubs = make(map[string]*Submission)

type QuestionsData struct {
	SERVER    string
	Questions []string
}

var QUESTIONS_TEMPLATE = `
<html>
	<head>
  		<title>Questions</title>
  		<script src="http://code.jquery.com/jquery-3.1.1.min.js"></script>
	    <script type="text/javascript">
			var updateInterval = 5000;		// 5 sec update interval
			var maxUpdateTime =  1800000;   // no longer update after 30 min.
			var totalUpdateTime = 0;
			function getQuestions() {
				var url = "http://localhost:4030/get_questions";
				$.getJSON(url, function( questions ) {
					$("#content").html("");
					$.each(questions, function(key, value) {
						// console.log(key, value);
						$("#content").append("<li>" + value + "</li>");
					});
				});
			}
			$(document).ready(function(){
				getQuestions();
				handle = setInterval(getQuestions, updateInterval);
			});
	    </script>
	</head>
	<body>
	<h1>Questions</h1>
	<ol>
		<div id="content"></div>
	</ol>
	</body>
</html>
`

type PollData struct {
	Description string
}

var POLL_ON = false
var POLL_DESCRIPTION = ""
var POLL_RESULT = make(map[string]string)
var POLL_COUNT = make(map[string]int)
var POLL_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
  	<script src="http://code.jquery.com/jquery-3.1.1.min.js"></script>
    <script type="text/javascript">
		var updateInterval = 3000;
		var maxUpdateTime = 300000;
		var totalUpdateTime = 0;
      google.charts.load('current', {'packages':['corechart']});
      google.charts.setOnLoadCallback(drawChart);
      function drawChart() {
      	// console.log(totalUpdateTime);
			var poll_view_url = "http://localhost:4030/query_poll";
			$.getJSON(poll_view_url, function( result ) {
				var data = new google.visualization.DataTable();
				data.addColumn('string', 'Answer');
				data.addColumn('number', 'Votes');
				$.each( result, function( key, val ) {
					data.addRow([key,val]);
				});
		      var options = {
		        	title: 'Poll Results',
		     		chartArea: {width: 600, height: 400}, width: 1000, height: 600,
		      };
	         var chart = new google.visualization.BarChart(document.getElementById('chart_div'));
	         chart.draw(data, options);
	      });
			totalUpdateTime += updateInterval;
			if (totalUpdateTime > maxUpdateTime) {
				clearInterval(handle);
			}
      }
	   $(document).ready(function(){
			handle = setInterval(drawChart, updateInterval);
      });
    </script>
  </head>

  <body>
  	<div style="padding:20px 0 0 200px"><pre id="description">{{.Description}}</pre></div>
    <div id="chart_div"></div>
  </body>
</html>
`

var TRACK_SUBMISSIONS_TEMPLATE = `
<html>
	<head>
  		<title>Track Submissions</title>
		<meta http-equiv="refresh" content="10" />
	<style>
		div {
		    font-family: monospace;
		    font-size: 150%;
		    color: red;
		    padding-top:0.5em;
		    padding-left:0.5em;
		}
	</style>
	</head>
	<body>
	<div>
	{{.Message}}
	</div>
	</body>
</html>
`

var TRACK_BOARD_TEMPLATE = `
<html>
	<head>
  		<title>Track Virtual Board</title>
		<meta http-equiv="refresh" content="10" />
	</head>
	<style>
		pre {
		    font-family: monospace;
		    font-size: 100%;
		    color: red;
		    padding-top:1em;
		    padding-left:0.5em;
		}
	</style>
	<body>
	<pre>{{.Message}}</pre>
	</body>
</html>
`

const YOU_GOT_CODE = `
Y
O
U

G
O
T

C
O
D
E
`
