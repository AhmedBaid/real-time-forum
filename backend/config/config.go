package config

import (
	"database/sql"
	"net/http"
	"time"
)

type Users struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Categories struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	PostID int    `json:"post_id"`
}

type Comments struct {
	PostID               int       `json:"post_id"`
	Id                   int       `json:"Id"`
	Username             string    `json:"Username"`
	Comment              string    `json:"Comment"`
	Time                 time.Time `json:"time"`
	TimeFormattedComment string
	TotalLikes           int `json:"TotalLikes"`
	TotalDislikes        int `json:"TotalDislikes"`
	UserReactionComment  int `json:"UserReactionComment"`
}

type Catgs struct {
	Catgs []string `json:"catgs"`
}

type Posts struct {
	Id                int          `json:"id"`
	Username          string       `json:"username"`
	Title             string       `json:"title"`
	Description       string       `json:"description"`
	Time              time.Time    `json:"time"`
	TotalLikes        int          `json:"totalLikes"`
	TotalDislikes     int          `json:"totalDislikes"`
	Comments          []Comments   `json:"comments"`
	Categories        []Categories `json:"categories"`
	TotalComments     int          `json:"totalComments"`
	UserReactionPosts int          `json:"userReactionPosts"`
}

type ErrorPage struct {
	Code         int
	ErrorMessage string
}

var (
	Db *sql.DB

	ErrorBadReq = ErrorPage{
		Code:         http.StatusBadRequest,
		ErrorMessage: "Oops! It looks like there was an issue with your request. Please check your input and try again.",
	}

	ErrorNotFound = ErrorPage{
		Code:         http.StatusNotFound,
		ErrorMessage: "Uh-oh! The page you're looking for doesn't exist. It might have been moved or deleted.",
	}

	ErrorMethodnotAll = ErrorPage{
		Code:         http.StatusMethodNotAllowed,
		ErrorMessage: "The request method is not supported for this resource. Please check and try again with a valid method.",
	}

	ErrorInternalServerErr = ErrorPage{
		Code:         http.StatusInternalServerError,
		ErrorMessage: "Something went wrong on our end. We're working on fixing it—please try again later!",
	}
	ErrorToManyRequests = ErrorPage{
		Code:         http.StatusTooManyRequests,
		ErrorMessage: "Rate limit exceeded",
	}
	ErrorUnauthorized = ErrorPage{
		Code:         http.StatusUnauthorized,
		ErrorMessage: "Unauthorized access. Please log in to continue.",
	}
)
