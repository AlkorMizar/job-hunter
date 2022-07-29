package postgres

import (
	"database/sql"
	"fmt"

	"github.com/AlkorMizar/job-hunter/internal/model/repo"
	"github.com/aws/smithy-go/time"
	"go.uber.org/zap"
)

type UserPostgres struct {
	repo.User
	DateCreatedStr string `json:"date_created" db:"date_created"`
	LastCheckStr   string `json:"last_check" db:"last_check"`
}

func (up *UserPostgres) ConvertToUser() (user repo.User, err error) {
	user = up.User

	user.DateCreated, err = time.ParseDateTime(up.DateCreatedStr)

	if err != nil {
		return repo.User{}, fmt.Errorf("during user conversion: invald user %d date, %w", user.ID, err)
	}

	user.LastCheck, err = time.ParseDateTime(up.LastCheckStr)
	if err != nil {
		return repo.User{}, fmt.Errorf("during user conversion: invald user %d date, %w", user.ID, err)
	}

	return
}

func (r *Repository) CreateUser(user *repo.User) (err error) {
	var id int

	query := "INSERT INTO \"user\" (login, email, password, full_name) values ($1,$2,$3,$4) RETURNING user_id"

	err = r.db.QueryRow(query, user.Login, user.Email, user.Password, user.FullName).Scan(&id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("in CreateUser can't insert user")
	}

	if err != nil {
		r.log.Debug("error during insert new user", zap.String("func", "CreateUser"), zap.Error(err))
		return fmt.Errorf("in CreateUser failed to exec query %s", query)
	}

	return nil
}

func (r *Repository) GetUserWithEamil(email string) (user repo.User, err error) {
	query := "SELECT * FROM \"user\" WHERE email=$1"

	var postgrUser UserPostgres

	err = r.db.Get(&postgrUser, query, email)
	if err != nil {
		r.log.Debug("error during selecting user by email", zap.String("func", "GetUserWithEamil"), zap.Error(err))
		return user, fmt.Errorf("in GetUserWithEamil can't find user")
	}

	user, err = postgrUser.ConvertToUser()
	if err != nil {
		r.log.Warn("invalid date fields in user", zap.String("func", "GetUserWithEamil"), zap.Error(err))
		return user, fmt.Errorf("in GetUserWithEamil invalid user id db")
	}

	return user, err
}

func (r *Repository) GetRoles(user *repo.User) (roles map[string]struct{}, err error) {
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
	tx, err := r.db.Begin()

	if err != nil {
		r.log.Debug("In SetRoles can't begin transaction", zap.Error(err))
		return fmt.Errorf("failed r.db.Begin with error: %w", err)
	}

	defer func() {
		if err != nil {
			errRoll := tx.Rollback()
			if err != nil {
				r.log.Debug("during transaction rollback", zap.String("func", "SetRoles"), zap.Error(errRoll))

				err = fmt.Errorf("in SerRoles error during rollback")
			}

			return
		}

		e := tx.Commit()
		if e != nil {
			r.log.Debug("In SetRoles can't commit transaction", zap.Error(err))
			err = fmt.Errorf("failed commit transaction with error: %w", e)
		}
	}()

	del := `DELETE FROM user_has_role WHERE fk_user_id=?`

	_, err = tx.Exec(del, user.ID)
	if err != nil {
		r.log.Debug("in SetRoles error during delete", zap.Error(err))
		return fmt.Errorf("failed exec DELETE query")
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
			r.log.Debug("in SetRoles error during inser", zap.Error(err))
			return fmt.Errorf("failed exec INSERT query")
		}
	}

	return nil
}
