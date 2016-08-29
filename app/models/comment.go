package models

import "time"

type Comment struct {
	Id uint
	Content string
	CreatedAt time.Time
	CreatedByID uint
	PostID uint
}
