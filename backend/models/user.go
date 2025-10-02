package models

type User struct {
	Id               int    `json:"id"`
	FirstName        string `json:"fist_name"`
	LastName         string `json:"last_name"`
	UserName         string `json:"user_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	ConfiremPassword string `json:"verifiy_password"`
}
