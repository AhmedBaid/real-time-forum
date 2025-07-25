package config

import (
	"database/sql"
	"time"
)

var Db *sql.DB

type Users struct {
	Username string
	Email    string
	Password string
}
type Categories struct {
	Name   string
	Id     int
	PostID int
}
type Posts struct {
	Id                int
	Username          string
	Title             string
	Description       string
	Time              time.Time
	TimeFormatted     string
	TotalLikes        int
	TotalDislikes     int
	Comments          []Comments
	Categories        []Categories
	TotalComments     int
	UserReactionPosts int
}
type Catgs struct {
	Catgs []string
}
type Comments struct {
	PostID               int
	Id                   int
	Username             string
	Comment              string
	Time                 time.Time
	TimeFormattedComment string
}
