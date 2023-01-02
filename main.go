package main

import (
	"log"
	"net/http"
)

func runApp() {

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	// Home page
	mux.HandleFunc("/", index)
	// error page
	mux.HandleFunc("/err", err)

	// User related routes: login, logout, registration and authentication
	mux.HandleFunc("/login", login)       // get login page
	mux.HandleFunc("/auth", authenticate) // Internal facing route which is used to create a session for the user to be logged in

	mux.HandleFunc("/logout", logout)            // logout a logon user
	mux.HandleFunc("/signup", signup)            // get signup page
	mux.HandleFunc("/register", registerAccount) // POST req. Internal facing route which is used to create a new user

	// Threads related routes: read, create and delete
	mux.HandleFunc("/thread/new", newThreadPage)      // Display the form for creating a thread
	mux.HandleFunc("/thread/create", createThread)    // POST req. Internal facing route that is used to create a new thread
	mux.HandleFunc("/thread/show", showThread)        // Show the details of a thread
	mux.HandleFunc("/thread/post/create", createPost) // POST req. Internal facing route for adding a post to its thread.

	// Run the server
	log.Fatal(http.ListenAndServe("localhost:8001", mux))
}

func main() {
	runApp()
}
