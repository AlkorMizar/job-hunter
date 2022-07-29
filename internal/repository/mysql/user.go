package mysql

import (
	"fmt"

	"github.com/AlkorMizar/job-hunter/internal/model/repo"
	"go.uber.org/zap"
)

func (r *Repository) CreateUser(user *repo.User) (err error) {
	query := "INSERT INTO user (login, email, password, full_name) values (:login,:email,:password,:full_name)"
	res, err := r.db.NamedExec(query, user)

	if err != nil {
		r.log.Debug("In CreateUser error during insert", zap.Error(err))
		return fmt.Errorf("in CreateUser failed to exec query %s with", query)
	}

	num, err := res.RowsAffected()

	if err != nil {
		r.log.Debug("In CreateUser failed RowAffected", zap.Error(err))
		return fmt.Errorf("in CreateUser RowsAffected failed ")
	}

	if num != 1 {
		return fmt.Errorf("no data were inserted")
	}

	return nil
}

func (r *Repository) GetUserWithEamil(email string) (user repo.User, err error) {
	query := "SELECT * FROM user WHERE email=?"

	err = r.db.Get(&user, query, email)

	if err != nil {
		return user, fmt.Errorf("in GetUserWithEmail can't get user with error %w", err)
	}

	return user, err
}

func (r *Repository) GetRoles(user *repo.User) (roles map[string]struct{}, err error) {
	roles = make(map[string]struct{})
	rolesArr := []repo.Role{}
	query := "SELECT role.name from role JOIN user_has_role ON fk_user_id=? AND fk_role_id=role_id; "
	err = r.db.Select(&rolesArr, query, user.ID)

	if err != nil {
		return nil, fmt.Errorf("in GetRoles failed to select with error: %w", err)
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
		return fmt.Errorf("in SetRoles failed begin transaction")
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
			err = fmt.Errorf("in SetRoles failed commit transaction")
		}
	}()

	_, err = tx.Exec(`DELETE FROM user_has_role WHERE fk_user_id=?`, user.ID)
	if err != nil {
		r.log.Debug("in SetRoles error during delete", zap.Error(err))
		return fmt.Errorf("failed exec DELETE query")
	}

	for k := range user.Roles {
		res, err := tx.Exec(`
		insert into user_has_role (fk_user_id, fk_role_id)
		select ?, role_id from role
		where name = ?;`, user.ID, k)
		if err != nil {
			r.log.Debug("in SetRoles error during inser", zap.Error(err))
			return fmt.Errorf("failed exec INSERT query")
		}

		num, err := res.RowsAffected()

		if err != nil {
			r.log.Debug("In SetRoles failed RowAffected", zap.Error(err))
			return fmt.Errorf("RowsAffected failed")
		}

		if num != 1 {
			return fmt.Errorf("no data were inserted")
		}
	}

	return nil
}
