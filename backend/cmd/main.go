package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"real_time/backend/config"
	"real_time/backend/router"

	_ "modernc.org/sqlite"
)

func main() {
	var err error
	config.Db, err = sql.Open("sqlite", "./backend/database/db.db")
	if err != nil {
		fmt.Println("kayn error", err)
		return
	}

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

	config.Db.Exec(`INSERT INTO categories (name, icon) VALUES('Sport', '<i class="fa-solid fa-medal"></i>'),('Music', '<i class="fa-solid fa-music"></i>'),('Movies', '<i class="fa-solid fa-film"></i>'),('Science', '<i class="fa-solid fa-flask"></i>'),('Gym', '<i class="fa-solid fa-dumbbell"></i>'),('Tecknology', '<i class="fa-solid fa-microchip"></i>'),('Culture', '<i class="fa-solid fa-person-walking"></i>'),('Politics', '<i class="fa-solid fa-landmark"></i>');`)


	fmt.Println("server listening on http://localhost:8080/")
	router.Router()
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
		return
	}
}
