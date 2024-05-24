package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // Omitting password from JSON responses
}

type ResultResponse struct {
    RESULTS []User `json:"results"`
}
