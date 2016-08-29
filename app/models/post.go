package models

import "time"

type Post struct {
	Id uint
	ReadCount uint
	Title, Content, Description string
	CreatedByID uint
	CreatedAt time.Time
	Comments []Comment
}
