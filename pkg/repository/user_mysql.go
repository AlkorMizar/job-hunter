package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

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

type UserManagMysql struct {
	db sqlx.DB
}

func NewUserManagMsql(db *sqlx.DB) *UserManagMysql {
	return &UserManagMysql{
		db: *db,
	}
}

func (r *UserManagMysql) CreateUser(user *User) error {
	query := "INSERT INTO user (login, email, password, fullName) values (:login,:email,:password,:fullName)"
	res, err := r.db.NamedExec(query, user)

	if err != nil {
		return err
	}

	num, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if num != 1 {
		return fmt.Errorf("couldn't insert")
	}

	return nil
}

func (r *UserManagMysql) GetUser(email string) (User, error) {
	var user User

	query := "SELECT * FROM user WHERE email=?"

	err := r.db.Get(&user, query, email)

	return user, err
}

func (r *UserManagMysql) GetRoles(user *User) (map[string]struct{}, error) {
	roles := make(map[string]struct{})
	rolesArr := []Role{}
	query := "SELECT role.name from role JOIN user_has_role ON User_idUser=? AND Role_idRole=idRole; "
	err := r.db.Select(&rolesArr, query, user.ID)

	if err != nil {
		return nil, err
	}

	for _, v := range rolesArr {
		roles[v.Name] = struct{}{}
	}

	return roles, err
}

func (r *UserManagMysql) SetRoles(user *User) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM user_has_role WHERE User_idUser=?`, user.ID)
	if err != nil {
		return err
	}

	for k := range user.Roles {
		res, e := tx.Exec(`
		insert into user_has_role (User_idUser, Role_idRole)
		select ?, idRole from role
		where name = ?;`, user.ID, k)
		if e != nil {
			return e
		}

		num, e := res.RowsAffected()

		if e != nil {
			return e
		}

		if num != 1 {
			return fmt.Errorf("couldn't insert")
		}
	}

	defer func() {
		e := tx.Commit()
		if e != nil {
			err = e
		}
	}()

	return nil
}
