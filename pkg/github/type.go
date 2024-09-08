package github

import "time"

type User struct {
	Login string `json:"login"`
}

type Notification struct {
	Id         string
	Reason     string
	Repository string
	Subject    string
	Unread     bool
	UpdatedAt  time.Time
}
