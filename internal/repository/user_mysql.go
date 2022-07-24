package repository

import (
	"fmt"

	"github.com/AlkorMizar/job-hunter/internal/model/repo"
	"github.com/AlkorMizar/job-hunter/internal/util"
)

func (r *Repository) CreateUser(user *repo.User) (err error) {
	defer util.Wrap(&err, "in CreateUser")

	query := "INSERT INTO user (login, email, password, fullName) values (:login,:email,:password,:fullName)"
	res, err := r.db.NamedExec(query, user)

	if err != nil {
		return fmt.Errorf("failed to exec query %s with error: %w", query, err)
	}

	num, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("RowsAffected failed with error %w", err)
	}

	if num != 1 {
		return fmt.Errorf("no data were inserted")
	}

	return nil
}

func (r *Repository) GetUserWithEamil(email string) (user repo.User, err error) {
	defer util.Wrap(&err, "in GetUserWithEamil")

	query := "SELECT * FROM user WHERE email=?"

	err = r.db.Get(&user, query, email)

	return user, err
}

func (r *Repository) GetRoles(user *repo.User) (roles map[string]struct{}, err error) {
	defer util.Wrap(&err, "in GetRoles")

	roles = make(map[string]struct{})
	rolesArr := []repo.Role{}
	query := "SELECT role.name from role JOIN user_has_role ON User_idUser=? AND Role_idRole=idRole; "
	err = r.db.Select(&rolesArr, query, user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to select with error: %w", err)
	}

	for _, v := range rolesArr {
		roles[v.Name] = struct{}{}
	}

	return roles, err
}

func (r *Repository) SetRoles(user *repo.User) (err error) {
	defer util.Wrap(&err, "in SetRoles")

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed r.db.Begin with error: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM user_has_role WHERE User_idUser=?`, user.ID)
	if err != nil {
		return fmt.Errorf("failed exec DELETE query with error: %w", err)
	}

	for k := range user.Roles {
		res, err := tx.Exec(`
		insert into user_has_role (User_idUser, Role_idRole)
		select ?, idRole from role
		where name = ?;`, user.ID, k)
		if err != nil {
			return fmt.Errorf("failed exec INSERT query with error: %w", err)
		}

		num, err := res.RowsAffected()

		if err != nil {
			return fmt.Errorf("RowsAffected failed with error %w", err)
		}

		if num != 1 {
			return fmt.Errorf("no data were inserted")
		}
	}

	defer func() {
		e := tx.Commit()
		if e != nil && err == nil {
			err = fmt.Errorf("failed commit with error: %w", e)
		}
	}()

	return nil
}
