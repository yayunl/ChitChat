package main

import (
	"ChitChat/data"
	"fmt"
	"net/http"
)

// GET /thread/new
// newThreadPage gets the form for creating a new thread
func newThreadPage(w http.ResponseWriter, r *http.Request) {
	// Validate the user is login and has a valid session
	sess, err := validateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		loginUser, _ := sess.User()
		info("Login user is %+v", loginUser)
		md := data.MergedData{Agent: loginUser}
		generateHTML(w, md, "layout", "private.navbar", "flash", "new.Thread")
	}
}

// POST /thread/create
// createThread creates a new thread
func createThread(w http.ResponseWriter, r *http.Request) {
	// Validate the user is login and has a valid session
	sess, err := validateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// Check if the parse form works
		err = r.ParseForm()
		if err != nil {
			danger(err, "Cannot parse post form for creating a new thread")
		}
		// Get the user of the session
		loginUser, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user of the session")
		}
		// Extract the topic from the posted form
		topic := r.PostFormValue("topic")
		if _, err := loginUser.CreateThread(topic); err != nil {
			danger(err, "Cannot create a new thread by the user")
		}
		http.Redirect(w, r, "/", 302)
	}
}

// GET /thread/read
// showThread shows details of a thread with the uuid of the thread in the request's URL
func showThread(w http.ResponseWriter, r *http.Request) {
	// Get the uuid of the thread from the URL
	urlVals := r.URL.Query()
	uuid := urlVals.Get("id")
	info("Thread uuid is ", uuid)
	// Get the thread with uuid
	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		danger(err)
		displayErrMessage(w, r, "Cannot show thread")
	} else {
		// Validate if the user has logged in and has valid session
		sess, err := validateSession(w, r)

		if err != nil {
			md := data.MergedData{
				Threads: []data.Thread{thread},
				Agent:   data.User{Name: "Anonymous User"},
			}
			generateHTML(w, md, "layout", "public.navbar", "login", "flash")
			//http.Redirect(w, r, "/login", 302)
		} else {
			// Get the user of the session
			loginUser, err := sess.User()
			if err != nil {
				danger(err, "Cannot find the user of the session")
			} else {
				md := data.MergedData{
					Threads: []data.Thread{thread},
					Agent:   loginUser,
				}
				generateHTML(w, md, "layout", "private.navbar", "private.thread", "flash")
			}
		}

	}
}

// POST /thread/post
// createPost creates a new post in the thread
func createPost(w http.ResponseWriter, r *http.Request) {
	// Validate the user is login and has a valid session
	sess, err := validateSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		// Check if the parse form works
		err = r.ParseForm()
		if err != nil {
			danger(err, "Cannot parse post form for creating a new thread")
		}
		// Get the user of the session
		loginUser, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user of the session")
		}
		// Extract the data from the form for the post
		postBody := r.PostFormValue("body")
		threadUUID := r.PostFormValue("uuid")
		thread, err := data.ThreadByUUID(threadUUID)
		if err != nil {
			displayErrMessage(w, r, "Cannot create the post because it failed to read the thread")
		}
		if _, err := loginUser.CreatePost(thread, postBody); err != nil {
			danger(err, "Cannot create the post")
		}
		// Redirect to the page that displays the details of the thread
		http.Redirect(w, r, fmt.Sprintf("/thread/show?id=%s", threadUUID), 302)
	}

}
