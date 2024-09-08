package github

import "time"

type User struct {
	Login string `json:"login"`
}

type Notification struct {
	Id         string
	Reason     string
	Repository string
	Subject    *NotificationSubject
	Unread     bool
	UpdatedAt  time.Time
}

type NotificationSubject struct {
	Title            string
	URL              string
	LatestCommentURL string
	Type             string
}
