package database

type Person struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Patronymic string `json:"patronymic"`
	Gender     string `json:"gender"`
	Age        int    `json:"age"`
	Country    string `json:"country"`
}
