package domain

type User struct {
	ID       int64   `bson:"id"`
	Username string  `bson:"username"`
	Password string  `bson:"password"`
	Email    string  `bson:"email"`
	Balance  float64 `bson:"balance"`
	Role     string  `bson:"role"`
}
