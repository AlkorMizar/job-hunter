package repo

import "time"

type User struct {
	ID          int       `json:"idUser" db:"idUser"`
	Login       string    `json:"login" db:"login"`
	FullName    string    `json:"fullName" db:"fullName"`
	Email       string    `json:"email" db:"email"`
	Password    []byte    `json:"password" db:"password"`
	DateCreated time.Time `json:"datecreated" db:"dateCreated"`
	LastCheck   time.Time `json:"lastcheck" db:"lastCheck"`
	Roles       map[string]struct{}
}

type Role struct {
	ID   int    `json:"idRole" db:"idRole"`
	Name string `json:"name" db:"name"`
}
