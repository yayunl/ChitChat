package data

import (
	"time"
)

// User struct and it's associated CRUD
type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Create a new user and save user info into the db //
func (u *User) Create() (err error) {
	statement := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// Query the newly inserted record and use the returned data to update the User struct
	err = stmt.QueryRow(createUUID(), u.Name, u.Email, Encrypt(u.Password), time.Now()).Scan(&u.Id, &u.Uuid, &u.CreatedAt)
	return
}

// Users retrieves the users from the db //
func Users() (users []User, err error) {
	// Get all users in the table `users` of the db
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	// Iterate over rows and append each user into an array of User struct
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

func UserByEmail(email string) (user User, err error) {
	// Get a user by email
	// Query a row
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).Scan(
		&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func UserByUUID(uuid string) (user User, err error) {
	// Get a user by uuid
	// Query a row
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).Scan(
		&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Update a user record in the db //
func (u *User) Update() (err error) {
	statement := "UPDATE users set name = $2, email = $3 WHERE id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Id, u.Name, u.Email)
	return
}

// Delete users from the db //
func (u *User) Delete() (err error) {
	// Delete a user given Id
	statement := "DELETE from users WHERE id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Id)
	return
}

func DeleteAllUsers() (err error) {
	// Delete all users in the db
	statement := "DELETE from users"
	_, err = Db.Exec(statement)
	return
}

// Helper functions //

// User retrieves a user from the table `users` with the field `UserId` in the session within the request sent by the browser/client
func (s *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", s.UserId).Scan(
		&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}
