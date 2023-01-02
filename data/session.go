package data

import (
	"fmt"
	"time"
)

// Session and it's associated CRUD functions
type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// CreateSession creates a new session for a valid user that is logon
func (u *User) CreateSession() (s Session, err error) {
	// Create a session for a logon user
	statement := "INSERT INTO sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// Use QueryRow to retrieve a row and scan the returned data into the Session struct
	err = stmt.QueryRow(createUUID(), u.Email, u.Id, time.Now()).Scan(&s.Id, &s.Uuid, &s.Email, &s.UserId, &s.CreatedAt)
	return
}

// Session gets the session for an existing user
func (u *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", u.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Delete deletes a session
func (s *Session) Delete() (err error) {
	// Delete a session given uuid
	statement := "DELETE from sessions WHERE uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.Uuid)
	return
}

func DeleteAllSessions() (err error) {
	// Delete all sessions
	_, err = Db.Exec("DELETE FROM sessions")
	return
}

// Helper functions

// Valid function is to check if the session struct of the method receiver exists in the table `sessions` in the database
func (s *Session) Valid() (valid bool, err error) {
	fmt.Printf("UUID is %s\n", s.Uuid)
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", s.Uuid).Scan(
		&s.Id, &s.Uuid, &s.Email, &s.UserId, &s.CreatedAt)
	if err != nil || s.Id == 0 {
		valid = false
	} else {
		valid = true
	}
	return
}
