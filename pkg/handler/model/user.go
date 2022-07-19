package model

type NewUser struct {
	Login     string   `json:"login" binding:"required" minimum:"5" maximum:"40" default:"test" validate:"required,min=3,max=40,alphanum"`
	Email     string   `json:"email" binding:"required" maximum:"255"  default:"test@test.com" validate:"required,email"`
	Roles     []string `json:"roles" binding:"required" validate:"required"`
	Password  string   `json:"password" binding:"required"  minimum:"5" maximum:"40" default:"test1" validate:"required,eqfield=CPassword"`
	CPassword string   `json:"cPassword" binding:"required"  minimum:"5" maximum:"40" default:"test1" validate:"required,min=5,max=40"`
}

type AuthInfo struct {
	Email    string `json:"email" binding:"required" maximum:"255" default:"test@test.com" validate:"required,email"`
	Password string `json:"password" binding:"required"  minimum:"5" maximum:"40" default:"test1" validate:"required,min=5,max=40"`
}

type JSONResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type User struct {
	Login    string   `json:"login" binding:"required" minimum:"5" maximum:"40" default:"test" validate:"required,min=3,max=40"`
	Email    string   `json:"email" binding:"required" maximum:"255"  default:"test@test.com" validate:"required,email"`
	Roles    []string `json:"roles" binding:"required" validate:"required"`
	FullName string   `json:"fullName" binding:"required" minimum:"5" maximum:"150" default:"test" validate:"required,min=3,max=150"`
}

type UpdateInfo struct {
	Login    string `json:"login" default:"test" validate:"required_without_all=Email FullName,omitempty,min=3,max=40,alphanum"`
	Email    string `json:"email" default:"test@test.com" validate:"required_without_all=Login FullName,omitempty,email"`
	FullName string `json:"fullName" default:"test" validate:"required_without_all=Login Email,omitempty,min=3,max=150"`
}

type Passwords struct {
	CurrPassword string `json:"CurrPasswoord" binding:"required"  minimum:"5" maximum:"40" default:"test1" validate:"required,min=5,max=40,excludesall= "`
	NewPassword  string `json:"NewPassword" binding:"required"  minimum:"5" maximum:"40" default:"test1" validate:"required,min=5,max=40,eqfield=CPassword,excludesall= "`
	CPassword    string `json:"CPassword" binding:"required"  minimum:"5" maximum:"40" default:"test1" validate:"required,min=5,max=40,excludesall= "`
}
