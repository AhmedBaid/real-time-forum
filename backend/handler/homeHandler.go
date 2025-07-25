package handler

import (
	"html/template"
	"net/http"

	// "real_time/backend/helpers"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("frontend/main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = temp.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// booll, session := helpers.SessionChecked(w, r)
	// if booll {
		
	// }
}
