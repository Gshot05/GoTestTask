package database

import (
	"fmt"
)

func Migrate() {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS people (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(255),
            last_name VARCHAR(255),
            patronymic VARCHAR(255),
            gender VARCHAR(255),
            age INTEGER,
            country VARCHAR(255)
        );
    `)

	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	fmt.Println("Migration successful")
}
