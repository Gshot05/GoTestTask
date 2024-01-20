package service

import (
	"encoding/json"
	"fmt"
	"myapp/database"
	"net/http"
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
