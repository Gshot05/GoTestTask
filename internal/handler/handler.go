package handler

import (
	"encoding/json"
	"myapp/database"
	"myapp/internal/service"
	"net/http"
	"strconv"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handlePost(w, r)
	} else if r.Method == http.MethodDelete {
		HandleDelete(w, r)
	} else if r.Method == http.MethodGet {
		HandleGetPeople(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var person database.Person

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&person)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = service.EnrichPersonData((*service.Person)(&person))
	if err != nil {
		http.Error(w, "Error enriching data", http.StatusInternalServerError)
		return
	}

	err = database.InsertPerson(&person)
	if err != nil {
		http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = service.DeletePersonByID(id)
		if err != nil {
			http.Error(w, "Error deleting data from database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func HandleGetPeople(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}

		people, err := service.GetPeopleWithLimit(limit)
		if err != nil {
			http.Error(w, "Error selecting data from database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(people)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
