// main.go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "real_time/backend/config"
    "real_time/backend/handler"
    "real_time/backend/helpers"
    "real_time/backend/router"

    _ "modernc.org/sqlite"
)

func main() {
    var err error
    config.Db, err = sql.Open("sqlite", "./backend/database/db.db")
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

    go handler.HandleBroadcast(config.Db)

    http.HandleFunc("/api/current-user", func(w http.ResponseWriter, r *http.Request) {
        _, session := helpers.SessionChecked(w, r)
        var userID int
        var username string
        err := config.Db.QueryRow("SELECT id,  username FROM users WHERE session = ?", session).Scan(&userID, &username)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{"userId": userID, "username" :  username})
    })

    http.HandleFunc("/ws", handler.WsHandler)
    http.HandleFunc("/messages", handler.GetMessagesHandler)
    http.HandleFunc("/mark-read/", handler.MarkReadHandler)

    router.Router()

    fmt.Println("Server listening on http://localhost:8080/")
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
}