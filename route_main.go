package main

import (
	"ChitChat/data"
	"fmt"
	"net/http"
)

// GET /err?msg=
// Display error messages
func err(writer http.ResponseWriter, request *http.Request) {
	urlVals := request.URL.Query()
	// Validate if the user has logged in and has a valid session
	_, err := validateSession(writer, request)
	if err != nil {
		// anonymous user
		generateHTML(writer, urlVals.Get("msg"), "layout", "public.navbar", "error", "flash")
	} else {
		// logged-in user
		generateHTML(writer, urlVals.Get("msg"), "layout", "private.navbar", "error", "flash")
	}

}
func index(writer http.ResponseWriter, request *http.Request) {
	info("Access to chitchat.")
	threads, err := data.Threads()
	if err != nil {
		displayErrMessage(writer, request, "Cannot get threads")
	} else {
		// Validate if the user has logged in and has a valid session
		sess, err := validateSession(writer, request)
		info("The validated session is %+v", sess)
		if err != nil {
			warning(err, "Invalid session")
			md := data.MergedData{
				Threads: threads,
				Agent:   data.User{Name: "Anonymous User"},
				Mesg:    fmt.Sprintf("%s", "Please login to view contents."),
			}
			// anonymous user
			generateHTML(writer, md, "layout", "public.navbar", "flash", "index")
		} else {
			// logged-in user
			loginUser, _ := sess.User()
			info("Login user is %+v", loginUser)
			md := data.MergedData{Threads: threads, Agent: loginUser}
			generateHTML(writer, md, "layout", "private.navbar", "flash", "index")
		}
	}
}
