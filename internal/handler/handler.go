package handler

import (
	"encoding/json"
	"myapp/database"
	"myapp/internal/service"
	"net/http"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var person database.Person

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&person)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Обогащение данных
		err = service.EnrichPersonData((*service.Person)(&person))
		if err != nil {
			http.Error(w, "Error enriching data", http.StatusInternalServerError)
			return
		}

		// Вставка обогащенных данных в базу данных
		err = database.InsertPerson(&person)
		if err != nil {
			http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(person)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
