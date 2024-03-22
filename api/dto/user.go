package dto

type User struct {
	Username string `json:"username" bingo:"required"`
	Password string `json:"password" bingo:"required"`
}
