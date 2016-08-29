package view

import "time"

type PostViewModel struct {
	Id uint
	Title, Content string
	Username string
	CreatedAt time.Time
}
