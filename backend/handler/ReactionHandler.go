package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

	stmt2 := "SELECT id  FROM users WHERE session = ?"
	var userid int
	errr := config.Db.QueryRow(stmt2, session).Scan(&userid)

	if errr != nil {
		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": "error in database",
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}

	stmt := `select value from likes where postID = ? and  userID = ?`
	row := config.Db.QueryRow(stmt, reaction.PostID, userid)
	var reactionValue string
	errrr := row.Scan(&reactionValue)
	if errrr == sql.ErrNoRows {
		// make the like
		stmt := `insert into likes (postID, userID, value) values(?, ?, ?)`
		_, err := config.Db.Exec(stmt, reaction.PostID, userid, reaction.Like)
		if err != nil {
			fmt.Println(err)
			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "error in database 1 ",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}
		fmt.Println(reaction, "1")
		config.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "liked   successful",
			"status":  http.StatusOK,
			"data":    reaction,
		})

	} else {
		if reactionValue == reaction.Like {
			// delete the like
			stmt := `delete from likes where postID = ? and userID = ?`
			_, err := config.Db.Exec(stmt, reaction.PostID, userid)
			if err != nil {
				fmt.Println(err)
				config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
					"message": "error in database 2",
					"status":  config.ErrorInternalServerErr.Code,
				})
				return
			}
			fmt.Println(reaction, "2")

			config.ResponseJSON(w, http.StatusOK, map[string]any{
				"message": "liked   successful",
				"status":  http.StatusOK,
				"data":    reaction,
			})

			return
		} else {
			// update the like
			stmt := `update likes set value = ? where postID = ? and userID = ?`
			_, err := config.Db.Exec(stmt, reaction.Like, reaction.PostID, userid)
			if err != nil {
				config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
					"message": "error in database 3",
					"status":  config.ErrorInternalServerErr.Code,
				})
				return
			}
			fmt.Println(reaction, "3")

			config.ResponseJSON(w, http.StatusOK, map[string]any{
				"message": "liked   successful",
				"status":  http.StatusOK,
				"data":    reaction,
			})
			return
		}
	}
}
