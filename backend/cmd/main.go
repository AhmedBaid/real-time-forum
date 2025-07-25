package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"real_time/backend/config"
	"real_time/backend/router"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error
	config.Db, err = sql.Open("sqlite3", "./backend/database/db.db")
	if err != nil {
		fmt.Println("kayn error", err)
		return
	}
	defer config.Db.Close()
	query, err := os.ReadFile("./backend/database/query.sql")
	if err != nil {
		fmt.Println("error in readfile", err)
		return
	}
	_, err = config.Db.Exec(string(query))
	if err != nil {
		fmt.Println("error execute", err)
		return
	}
	fmt.Println("server listening on http://localhost:8080/")
	router.Router()
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
		return
	}
}
