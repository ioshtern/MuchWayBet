package models

type User struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
	Role     string  `json:"role"`
}
