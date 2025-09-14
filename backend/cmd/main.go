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

	_, err = config.Db.Exec(`
        INSERT OR IGNORE INTO categories (name, icon) VALUES
        ('Sport', '<i class="fa-solid fa-medal"></i>'),
        ('Music', '<i class="fa-solid fa-music"></i>'),
        ('Movies', '<i class="fa-solid fa-film"></i>'),
        ('Science', '<i class="fa-solid fa-flask"></i>'),
        ('Gym', '<i class="fa-solid fa-dumbbell"></i>'),
        ('Tecknology', '<i class="fa-solid fa-microchip"></i>'),
        ('Culture', '<i class="fa-solid fa-person-walking"></i>'),
        ('Politics', '<i class="fa-solid fa-landmark"></i>');
    `)
	if err != nil {
		fmt.Println("Error inserting categories:", err)
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
