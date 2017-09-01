//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"math/rand"
	"sync"
	"time"
)

const VERSION = "0.27"

var ADDR = ""
var PORT = "4030"
var USER_DB string
var SERVER = ""

type Board struct {
	Content      string
	Description  string
	StartingTime time.Time
	Changed      bool
	Ext          string
	Bid          string // id of current broadcast
}

var Boards = make(map[string]*Board)

type Submission struct {
	Sid       string // submission id
	Bid       string // broadcast id
	Uid       string // user id
	Body      string
	Ext       string
	Points    int
	Duration  int    // in seconds
	Pdes      string // problem description
	Timestamp string
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

var SEM sync.Mutex
var NewSubs = make([]*Submission, 0)
var ProcessedSubs = make(map[string]*Submission)

type TemplateData struct {
	SERVER string
}

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
    <div id="chart_div"></div>
  </body>
</html>
`
