package postgres

import (
	"database/sql"
	"fmt"

	"github.com/AlkorMizar/job-hunter/internal/model/repo"
	"github.com/AlkorMizar/job-hunter/internal/util"
	"github.com/aws/smithy-go/time"
	"go.uber.org/zap"
)

type UserPostgres struct {
	repo.User
	DateCreatedStr string `json:"date_created" db:"date_created"`
	LastCheckStr   string `json:"last_check" db:"last_check"`
}

func (up UserPostgres) ConvertToUser() (user repo.User, err error) {
	user = up.User
	user.DateCreated, err = time.ParseDateTime(up.DateCreatedStr)
	if err != nil {
		return repo.User{}, fmt.Errorf("during user convertion: invald user %d date, %w", user.ID, err)
	}
	user.LastCheck, err = time.ParseDateTime(up.LastCheckStr)
	if err != nil {
		return repo.User{}, fmt.Errorf("during user convertion: invald user %d date, %w", user.ID, err)
	}
	return
}

func (r *Repository) CreateUser(user *repo.User) (err error) {
	defer util.Wrap(&err, "in CreateUser")

	var id int

	query := "INSERT INTO \"user\" (login, email, password, full_name) values ($1,$2,$3,$4) RETURNING user_id"

	err = r.db.QueryRow(query, user.Login, user.Email, user.Password, user.FullName).Scan(&id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("can't insert user")
	}

	if err != nil {
		r.log.Debug("error during insert new user", zap.Error(err))
		return fmt.Errorf("failed to exec query %s", query)
	}

	return nil
}

func (r *Repository) GetUserWithEamil(email string) (user repo.User, err error) {
	defer util.Wrap(&err, "in GetUserWithEamil")

	query := "SELECT * FROM \"user\" WHERE email=$1"

	var postgrUser UserPostgres

	err = r.db.Get(&postgrUser, query, email)
	if err != nil {
		r.log.Debug("error during selecting user by email", zap.Error(err))

		err = fmt.Errorf("can't find user")
		return user, err
	}

	user, err = postgrUser.ConvertToUser()
	if err != nil {
		r.log.Debug("invalid date fields in user", zap.Error(err))

		err = fmt.Errorf("invalid user id db")
		return user, err
	}
	return user, err
}

func (r *Repository) GetRoles(user *repo.User) (roles map[string]struct{}, err error) {
	defer util.Wrap(&err, "in GetRoles")

	roles = make(map[string]struct{})

	rolesArr := []repo.Role{}

	query := "SELECT role.name from role JOIN user_has_role ON fk_user_id=$1 AND fk_role_id=role_id;"

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

	del := `DELETE FROM user_has_role WHERE fk_user_id=?`

	_, err = tx.Exec(del, user.ID)
	if err != nil {
		return fmt.Errorf("failed exec query %s with error: %w", del, err)
	}

	var id int

	insert := `INSERT INTO user_has_role (fk_user_id, fk_role_id)
			   SELECT $1, role_id from role
			   WHERE name = $2
			   RETURNING fk_user_id;`
	for roleName := range user.Roles {
		row := tx.QueryRow(insert, user.ID, roleName)
		err = row.Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to exec query %s with error: %w", insert, err)
		}
	}

	defer func() {
		e := tx.Commit()
		if e != nil && err == nil {
			err = fmt.Errorf("failed commit transaction with error: %w", e)
		}
	}()

	return nil
}
