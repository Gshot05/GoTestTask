package handler

import (
	"encoding/json"
	"log"
	"myapp/database"
	"myapp/internal/service"
	"net/http"
	"strconv"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodDelete:
		HandleDelete(w, r)
	case http.MethodGet:
		HandleGetPeople(w, r)
	case http.MethodPut:
		HandleUpdatePerson(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var person database.Person

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&person)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = service.EnrichPersonData((*service.Person)(&person))
	if err != nil {
		log.Println("Error enriching data:", err)
		http.Error(w, "Error enriching data", http.StatusInternalServerError)
		return
	}

	err = database.InsertPerson(&person)
	if err != nil {
		log.Println("Error inserting data into database:", err)
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
			log.Println("Invalid ID parameter:", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = service.DeletePersonByID(id)
		if err != nil {
			log.Println("Error deleting data from database:", err)
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
			log.Println("Invalid limit parameter:", err)
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}

		people, err := service.GetPeopleWithLimit(limit)
		if err != nil {
			log.Println("Error selecting data from database:", err)
			http.Error(w, "Error selecting data from database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(people)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func HandleUpdatePerson(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		var updatedPerson database.Person

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&updatedPerson)
		if err != nil {
			log.Println("Error decoding request payload:", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("Invalid ID parameter:", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = service.UpdatePerson(id, ((*service.Person)(&updatedPerson)))
		if err != nil {
			log.Println("Error updating data in the database:", err)
			http.Error(w, "Error updating data in the database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
