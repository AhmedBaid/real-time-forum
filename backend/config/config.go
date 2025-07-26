package config

import (
	"database/sql"
	"net/http"
	"time"
)
type Users struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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
	Id                   int       `json:"id"`
	Username             string    `json:"username"`
	Comment              string    `json:"comment"`
	Time                 time.Time `json:"time"`
	TimeFormattedComment string    `json:"time_formatted_comment"`
}

type Catgs struct {
	Catgs []string `json:"catgs"`
}

type Posts struct {
	Id                int         `json:"id"`
	Username          string      `json:"username"`
	Title             string      `json:"title"`
	Description       string      `json:"description"`
	Time              time.Time   `json:"time"`
	TimeFormatted     string      `json:"time_formatted"`
	TotalLikes        int         `json:"total_likes"`
	TotalDislikes     int         `json:"total_dislikes"`
	Comments          []Comments  `json:"comments"`
	Categories        []Categories`json:"categories"`
	TotalComments     int         `json:"total_comments"`
	UserReactionPosts int         `json:"user_reaction_posts"`
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
