package mysql

import (
	"fmt"
	"reflect"
	"strings"

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
		r.log.Debug("Error during select", zap.String("func", "mysql/GetUserWithEamil"), zap.Error(err))
		return user, fmt.Errorf("in GetUserWithEmail can't get user")
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

func (r *Repository) GetUserByID(id int) (user *repo.User, err error) {
	user = &repo.User{}

	query := "SELECT * FROM user WHERE user_id=?"

	err = r.db.Get(user, query, id)

	if err != nil {
		return nil, err
	}

	user.Roles, err = r.GetRoles(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *Repository) UpdateUser(id int, updateU *repo.User) (err error) {
	var setter strings.Builder
	setter.WriteString("SET ")

	elem := reflect.ValueOf(updateU).Elem()
	for i := 0; i < elem.NumField(); i++ {
		dbColName := elem.Type().Field(i).Tag.Get("db")
		value, ok := elem.Field(i).Interface().(string)
		if ok && value != "" {
			setter.WriteString(dbColName + "=\"" + value + "\",")
		}
	}

	setCols := setter.String()[0 : setter.Len()-1]
	if setCols == "SET" {
		return nil
	}

	query := "UPDATE user " + setCols + "WHERE user_id=?"
	res, err := r.db.Exec(query, id)

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

func (r *Repository) GetUserById(id int) (*repo.User, error) {
	var user repo.User

	query := "SELECT * FROM user WHERE user_id=?"

	err := r.db.Get(&user, query, id)

	if err != nil {
		r.log.Debug("Error during select", zap.String("func", "mysql/GetUserById"), zap.Error(err))
		return nil, fmt.Errorf("in GetUserById can't get user")
	}

	return &user, err
}

func (r *Repository) UpdatePassword(id int, pwd []byte) error {
	query := "UPDATE user SET password=? WHERE user_id=?"

	res, err := r.db.Exec(query, pwd, id)

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
