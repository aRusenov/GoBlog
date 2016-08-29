package models

const (
	USER = 0
	ADMIN = 1
)

type User struct {
	ID             uint
	Name, Username string
	HashedPassword []byte
	Posts          []Post
	Comments       []Comment
	Role		uint
}

