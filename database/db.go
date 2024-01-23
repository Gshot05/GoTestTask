package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	return db, nil
}

func InsertPerson(person *Person) error {
	_, err := db.Exec(`
        INSERT INTO people (first_name, last_name, patronymic, gender, age, country)
        VALUES ($1, $2, $3, $4, $5, $6);
    `, person.FirstName, person.LastName, person.Patronymic, person.Gender, person.Age, person.Country)

	if err != nil {
		log.Println("Error inserting person into the database:", err)
	}

	return err
}

func DeletePersonByID(id int) error {
	_, err := db.Exec(`
		DELETE FROM people WHERE id = $1;
	`, id)

	if err != nil {
		log.Println("Error deleting person from the database:", err)
	}

	return err
}

func SelectPeople(limit int) ([]Person, error) {
	rows, err := db.Query(`
		SELECT first_name, last_name, patronymic, gender, age, country FROM people
		LIMIT $1;
	`, limit)
	if err != nil {
		log.Println("Error selecting people from the database:", err)
		return nil, err
	}
	defer rows.Close()

	people := make([]Person, 0)
	for rows.Next() {
		var person Person
		err := rows.Scan(
			&person.FirstName,
			&person.LastName,
			&person.Patronymic,
			&person.Gender,
			&person.Age,
			&person.Country,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		people = append(people, person)
	}

	return people, nil
}

func UpdatePerson(id int, updatedPerson *Person) error {
	_, err := db.Exec(`
        UPDATE people
        SET first_name=$1, last_name=$2, patronymic=$3, gender=$4, age=$5, country=$6
        WHERE id=$7;
    `, updatedPerson.FirstName, updatedPerson.LastName, updatedPerson.Patronymic, updatedPerson.Gender, updatedPerson.Age, updatedPerson.Country, id)

	if err != nil {
		log.Println("Error updating person in the database:", err)
	}

	return err
}
