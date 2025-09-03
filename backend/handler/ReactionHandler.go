package handler

import (
	"database/sql"
	"encoding/json"
	"html"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		config.ResponseJSON(w, config.ErrorMethodnotAll.Code, map[string]any{
			"message": "mthod not allowd ",
			"status":  config.ErrorMethodnotAll.Code,
		})
		return
	}
	_, session := helpers.SessionChecked(w, r)

	var reaction config.Reactions

	err := json.NewDecoder(r.Body).Decode(&reaction)
	if err != nil {
		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": err,
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}
	if reaction.PostID == 0 || reaction.Like == "" || (reaction.Like != "1" && reaction.Like != "-1") {
		config.ResponseJSON(w, config.ErrorBadReq.Code, map[string]any{
			"message": "bad request",
			"status":  config.ErrorBadReq.Code,
		})
		return
	}

	// get user ID
	stmt2 := "SELECT id FROM users WHERE session = ?"
	var userid int
	err = config.Db.QueryRow(stmt2, session).Scan(&userid)
	if err != nil {
		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": "error in database",
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}

	// get user reaction (if exists)
	stmt := `SELECT value FROM likes WHERE postID = ? AND userID = ?`
	row := config.Db.QueryRow(stmt, reaction.PostID, userid)

	var reactionValue string
	err = row.Scan(&reactionValue)
reactionValue =  html.EscapeString(reactionValue)
	if err == sql.ErrNoRows {
		// insert
		stmt := `INSERT INTO likes (postID, userID, value) VALUES (?, ?, ?)`
		_, err := config.Db.Exec(stmt, reaction.PostID, userid, reaction.Like)
		if err != nil {
			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "error in database 1",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}
	} else if reactionValue == reaction.Like {
		// delete
		stmt := `DELETE FROM likes WHERE postID = ? AND userID = ?`
		_, err := config.Db.Exec(stmt, reaction.PostID, userid)
		if err != nil {
			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "error in database 2",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}
	} else {
		// update
		stmt := `UPDATE likes SET value = ? WHERE postID = ? AND userID = ?`
		_, err := config.Db.Exec(stmt, reaction.Like, reaction.PostID, userid)
		if err != nil {
			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "error in database 3",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}
	}

	stmtUserReaction := `
					SELECT value FROM likes WHERE postID = ? AND userID = ?
				`

	config.Db.QueryRow(stmtUserReaction, reaction.PostID, userid).Scan(&reaction.UserReactionPosts)

	stmtLikes := `SELECT COUNT(*)  FROM likes WHERE postID = ? AND value = '1'`
	config.Db.QueryRow(stmtLikes, reaction.PostID).Scan(&reaction.TotalLike)

	stmtDislikes := `SELECT COUNT(*) FROM likes WHERE postID = ? AND value = '-1'`
	config.Db.QueryRow(stmtDislikes, reaction.PostID).Scan(&reaction.TotalDislikes)
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "reaction updated successfully",
		"status":  http.StatusOK,
		"data":    reaction,
	})
}
