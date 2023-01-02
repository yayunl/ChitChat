package data

import (
	"fmt"
	"time"
)

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

// CreateThread creates a thread by a user
func (u *User) CreateThread(topic string) (t Thread, err error) {
	// Insert a new thread into the database
	statement := "INSERT INTO threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, topic, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	err = stmt.QueryRow(createUUID(), topic, u.Id, time.Now()).Scan(&t.Id, &t.Uuid, &t.Topic, &t.UserId, &t.CreatedAt)
	return
}

// Threads function gets all threads in the database
func Threads() (threads []Thread, err error) {
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
			return
		}
		threads = append(threads, conv)
	}
	rows.Close()
	return
}

// Helper functions

// CreatedAtDate formats the CreatedAt date to display nicely on the screen
func (t *Thread) CreatedAtDate() string {
	return fmt.Sprintf(t.CreatedAt.Format(time.ANSIC))
}

// NumOfPosts returns the number of posts in the thread
func (t *Thread) NumOfPosts() (c int) {
	rows, err := Db.Query("SELECT count(*) FROM posts WHERE thread_id = $1", t.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&c); err != nil {
			return
		}
	}
	defer rows.Close()
	return
}

// Author function gets the user who created the thread
func (t *Thread) Author() (a User) {
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", t.UserId).Scan(
		&a.Id, &a.Uuid, &a.Name, &a.Email, &a.CreatedAt)
	return
}

// ThreadByUUID get the thread given uuid
func ThreadByUUID(tuuid string) (t Thread, err error) {
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", tuuid).Scan(
		&t.Id, &t.Uuid, &t.Topic, &t.UserId, &t.CreatedAt)
	return
}
