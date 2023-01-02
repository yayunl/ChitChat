package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
)

type MergedData struct {
	Threads []Thread
	Agent   User
	Mesg    string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres password=password dbname=chitchat sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return
}

// Create a unique string for a new user
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// Encrypt user password with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

//type Post struct {
//	Id        int
//	Content   string
//	Author    string `sql:"not null"`
//	Comments  []Comment
//	CreatedAt time.Time
//}
//type Comment struct {
//	Id        int
//	Content   string
//	Author    string `sql:"not null"`
//	PostId    int    `sql:"index"`
//	CreatedAt time.Time
//}

//var Db *gorm.DB

//func init() {
//	//connStr := "postgres://postgres:postgrespw@localhost:49153/myDB?sslmode=disable"
//	dsn := "host=localhost user=postgres password=postgrespw dbname=myDB port=49153 sslmode=disable "
//	var err error
//	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		panic(err)
//	}
//	Db.AutoMigrate(&Post{}, &Comment{})
//}
