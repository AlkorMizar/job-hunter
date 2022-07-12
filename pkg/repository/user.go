package repository

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
