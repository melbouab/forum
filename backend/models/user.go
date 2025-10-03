// package models - User struct (corrected)
package models

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// Removed unused ConfirmPassword field
}
