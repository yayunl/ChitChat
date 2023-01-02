package main

import (
	"ChitChat/data"
	_ "ChitChat/data"
	"net/http"
)

// GET /login
func login(w http.ResponseWriter, r *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(w, nil)
}

// GET /signup
// signup function gets the register page
func signup(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// Create a new user account if it does not exit in the db
func registerAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}

	// Check if the user exists in the db
	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	// Create the user
	if err := user.Create(); err != nil {
		danger(err, "Cannot create this user.")
		// Redirect to login page
		http.Redirect(w, r, "/login", 302)
	}
	// Automatically login in the user and create the session
	authenticate(w, r)
}

// POST /auth
// authenticate the user given the email and password by creating a session associated with the user
func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	// Get the user given the email from the database
	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		danger(err, "The user does not exist in the db. Cannot create a session for it.")
		// redirect
		http.Redirect(w, r, "/login", 302)
	}
	// Create a session for the user if the user passes validation
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		info("Creating session for the user", user)
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create a session")
		}
		// Create a cookie to save the session data
		cookie := http.Cookie{
			Name:     "_chitchat_session",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		// Set the cookie and redirect to the home page
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
		info("Cookie is created")
	} else {
		danger("The user fails the encryption validation")
		// The user fails validation
		http.Redirect(w, r, "/login", 302)
	}
}

// GET /logout
// logout logs out a logon user by removing the associated session in the cookie in the request
func logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("_chitchat_session")
	// Remove the cookie from the db
	if err != http.ErrNoCookie {
		warning(err, "Failed to get the cookie")
		s := data.Session{Uuid: cookie.Value}
		s.Delete()
	}
	info("Logged out.")
	http.Redirect(w, r, "/login", 302)
}
