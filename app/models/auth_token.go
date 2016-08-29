package models

import "time"

type AuthToken struct {
	Id uint
	Selector string
	Token string
	UserID uint
	Expires time.Time
}
