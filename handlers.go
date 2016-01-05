package main 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"errors"
)

//-----------------------------------------------------------------

func verifyPasscode(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("passcode") != PASSCODE {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized access")
	}
	return nil
}

//-----------------------------------------------------------------
// users query to know their current points
//-----------------------------------------------------------------
func my_pointsHandler(w http.ResponseWriter, r *http.Request) {
	_, user := r.FormValue("login"), r.FormValue("uid")
	points, ok := Points.data[user]
	if !ok {
		points = 0
	}
	fmt.Fprintf(w, user+" has "+strconv.Itoa(points)+" brownies.")
}

//-----------------------------------------------------------------
// users register their uids and names.
//-----------------------------------------------------------------
func registerHandler(w http.ResponseWriter, r *http.Request) {
	uid, name := r.FormValue("uid"), r.FormValue("name")
	if _, ok := AllUsers[uid]; ok {
		fmt.Fprintf(w, uid + " already exists.")
	} else {
		RegisteredUsers[uid] = name
		fmt.Fprintf(w, "Waiting for approval.")
	}
}

//-----------------------------------------------------------------
// users submit their codes
//-----------------------------------------------------------------
func submit_postHandler(w http.ResponseWriter, r *http.Request) {
	_, user, body := r.FormValue("login"), r.FormValue("uid"), r.FormValue("body")
	if _, ok := AllUsers[user]; ok {
		Posts.Add(user, body)
		fmt.Println(user, "submitted.")
		fmt.Fprintf(w, "1")
	} else {
		fmt.Println(user, "non existent.")
		fmt.Fprintf(w, "0")
	}
}


//-----------------------------------------------------------------
// approve all currently registered users.
//-----------------------------------------------------------------

func approveHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		approved := make(map[string]string)
		err := json.Unmarshal([]byte(r.FormValue("approved")), &approved)
		if err != nil {
			fmt.Println(err)
		}
		for uid, name := range(approved) {
			delete(RegisteredUsers, uid)
			_, ok := AllUsers[uid] 
			if !ok {
				AllUsers[uid] = &User{name, 0}
				fmt.Println("\tApprove",uid,name)
			} else {
				fmt.Println("\t", uid, name, "already exists.")
			}
		}
	}
}


//-----------------------------------------------------------------
// return all currently registered users, waiting for approval.
//-----------------------------------------------------------------

func registered_usersHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		js, err := json.Marshal(RegisteredUsers)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// return points of currently awarded users
//-----------------------------------------------------------------
func pointsHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		js, err := json.Marshal(Points.data)
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
	if verifyPasscode(w, r) == nil {
		Points.addOne(r.FormValue("uid"))
		fmt.Println("+1", r.FormValue("uid"))
		fmt.Fprintf(w, "Point awarded to "+r.FormValue("uid"))
	}
}

//-----------------------------------------------------------------
// return all current posts
//-----------------------------------------------------------------
func postsHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		js, err := json.Marshal(Posts.queue)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	}
}

//-----------------------------------------------------------------
// Instructor retrieves code
//-----------------------------------------------------------------
func get_postHandler(w http.ResponseWriter, r *http.Request) {
	if verifyPasscode(w, r) == nil {
		e, err := strconv.Atoi(r.FormValue("post"))
		if err != nil {
			fmt.Println(err.Error)
		} else {
			js, err := json.Marshal(Posts.Remove(e))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
		}
	}
}
