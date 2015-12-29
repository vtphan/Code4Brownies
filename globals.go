package main 

//-----------------------------------------------------------------
// GLOBALS
//-----------------------------------------------------------------

var Posts = PostQueue{}                         // posts of currently active users
var Points = &Point{data: make(map[string]int)} // points of currently active users
var AllUsers = make(map[string]*User)           // maps uids to users
