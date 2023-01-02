package data

import "time"

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

// CreatePost creates a new post and add it to a thread
func (u *User) CreatePost(t Thread, body string) (post Post, err error) {
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, thread_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), body, u.Id, t.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}

// Helper functions

// Author function gets the user who wrote the post
func (p *Post) Author() (u User) {
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", p.UserId).Scan(&u.Id, &u.Uuid, &u.Name, &u.Email, &u.CreatedAt)
	return
}

// Posts function returns the posts wrote by the user
func (u *User) Posts() (posts []Post, err error) {
	// Get all posts wrote by the user
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts WHERE user_id = $1", u.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

// Posts function returns the posts in a thread
func (t *Thread) Posts() (posts []Post, err error) {
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts WHERE thread_id = $1", t.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	return
}

// CreatedAtDate formats the date field `CreatedAt` to display nicely on the browser
func (p *Post) CreatedAtDate() string {
	return p.CreatedAt.Format(time.ANSIC)
}
