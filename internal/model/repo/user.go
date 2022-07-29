package repo

import "time"

type User struct {
	ID          int       `json:"user_id" db:"user_id"`
	Login       string    `json:"login" db:"login"`
	FullName    string    `json:"full_name" db:"full_name"`
	Email       string    `json:"email" db:"email"`
	Password    []byte    `json:"password" db:"password"`
	DateCreated time.Time `json:"date_created" db:"date_created"`
	LastCheck   time.Time `json:"last_check" db:"last_check"`
	Roles       map[string]struct{}
}

type Role struct {
	ID   int    `json:"roleID" db:"role_id"`
	Name string `json:"name" db:"name"`
}
