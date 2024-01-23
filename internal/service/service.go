package service

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/database"
	"net/http"
	"strconv"
)

type Person database.Person

func getAPIResponse(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func enrichAge(person *Person) error {
	ageURL := fmt.Sprintf("https://api.agify.io/?name=%s", person.FirstName)
	var ageData map[string]interface{}
	err := getAPIResponse(ageURL, &ageData)
	if err != nil {
		log.Println("Error enriching age data:", err)
		return err
	}

	person.Age = int(ageData["age"].(float64))
	return nil
}

func enrichGender(person *Person) error {
	genderURL := fmt.Sprintf("https://api.genderize.io/?name=%s", person.FirstName)
	var genderData map[string]interface{}
	err := getAPIResponse(genderURL, &genderData)
	if err != nil {
		log.Println("Error enriching gender data:", err)
		return err
	}

	person.Gender = genderData["gender"].(string)
	return nil
}

func enrichNationality(person *Person) error {
	nationalizeURL := fmt.Sprintf("https://api.nationalize.io/?name=%s", person.FirstName)
	var nationalizeData map[string]interface{}
	err := getAPIResponse(nationalizeURL, &nationalizeData)
	if err != nil {
		log.Println("Error enriching nationality data:", err)
		return err
	}

	countries := nationalizeData["country"].([]interface{})
	if len(countries) > 0 {
		countryData := countries[0].(map[string]interface{})
		person.Country = countryData["country_id"].(string)
	}

	return nil
}

func EnrichPersonData(person *Person) error {
	err := enrichAge(person)
	if err != nil {
		return err
	}

	err = enrichGender(person)
	if err != nil {
		return err
	}

	err = enrichNationality(person)
	if err != nil {
		return err
	}

	return nil
}

func DeletePersonByID(id int) error {
	return database.DeletePersonByID(id)
}

func DeletePersonByIDFromRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("Invalid ID parameter:", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = DeletePersonByID(id)
		if err != nil {
			log.Println("Error deleting person:", err)
			http.Error(w, "Error deleting person", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func GetPeopleWithLimit(limit int) ([]Person, error) {
	peopleFromDB, err := database.SelectPeople(limit)
	if err != nil {
		log.Println("Error selecting people from database:", err)
		return nil, err
	}

	var people []Person
	for _, p := range peopleFromDB {
		people = append(people, Person(p))
	}

	return people, nil
}

func UpdatePerson(id int, updatedPerson *Person) error {
	err := database.UpdatePerson(id, (*database.Person)(updatedPerson))
	if err != nil {
		log.Println("Error updating person:", err)
		return err
	}

	return nil
}
