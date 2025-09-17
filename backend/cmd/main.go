// main.go
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"real_time/backend/config"
	"real_time/backend/router"

	_	"github.com/mattn/go-sqlite3"
)

func main() {
	var err error
	config.Db, err = sql.Open("sqlite3", "./backend/database/db.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer config.Db.Close()

	query, err := os.ReadFile("./backend/database/query.sql")
	if err != nil {
		fmt.Println("Error reading query file:", err)
		return
	}

	_, err = config.Db.Exec(string(query))
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}

	router.Router()

	fmt.Println("Server listening on http://localhost:8080/")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
