package domain

type User struct {
	Username string  `bson:"username"`
	Password string  `bson:"password"`
	Email    string  `bson:"email"`
	Balance  float64 `bson:"balance"`
	Role     string  `bson:"role"`
}
