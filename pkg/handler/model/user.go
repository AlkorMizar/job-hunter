package model

type NewUser struct {
	Login     string `json:"login" binding:"required" minimum:"5" maximum:"40" validate:"required,min=3,max=40"`
	Email     string `json:"email" binding:"required" maximum:"255" validate:"required,email"`
	Password  string `json:"password" binding:"required"  minimum:"5" maximum:"40" validate:"required,eqfield=CPassword"`
	CPassword string `json:"cPassword" binding:"required"  minimum:"5" maximum:"40" validate:"required,min=5,max=40"`
}

type AuthInfo struct {
	Email    string `json:"email" binding:"required" maximum:"255" validate:"required,email"`
	Password string `json:"password" binding:"required"  minimum:"5" maximum:"40" validate:"required,email,min=5,max=40"`
}
