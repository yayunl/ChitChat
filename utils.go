package main

import (
	"ChitChat/data"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var logger *log.Logger

func init() {
	logFile, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(logFile, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
}

// Helper functions
// displayErrMessage function is to redirect error messages to the error message page
func displayErrMessage(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// validateSession validates if the session struct created with using the data from the client/browser's cookie exists in the database.
// In other words, the function is to check if the user has logged in and has a valid session.
func validateSession(writer http.ResponseWriter, request *http.Request) (session data.Session, err error) {
	cookie, err := request.Cookie("_chitchat_session")
	if err == nil {
		info("Cookie is ", cookie.Value)
		// User is logged in and has a session
		session = data.Session{Uuid: cookie.Value}
		// Validate if the user's session is valid, meaning if the session exists in the db
		ok, _ := session.Valid()
		if !ok {
			err = errors.New("invalid session as it does not exist in the db")
			danger(err)
		}
	} else {
		danger(err, "Cannot get the cookie")
	}

	return
}

// parse HTML templates
// pass in a list of file names, and get a template
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}

// logging
// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func warning(args ...interface{}) {
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}
