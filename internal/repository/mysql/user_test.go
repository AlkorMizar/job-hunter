package mysql_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/AlkorMizar/job-hunter/internal/model/repo"
	"github.com/AlkorMizar/job-hunter/internal/repository/mysql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var pwd, _ = bcrypt.GenerateFromPassword([]byte("password"), 7)

func TestCreateUserCorrect(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:       1,
		Login:    "login",
		Email:    "eamil@d.com",
		Password: pwd,
		FullName: "full name",
	}

	mock.ExpectExec("INSERT INTO user").
		WithArgs(user.Login, user.Email, user.Password, user.FullName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	if err := rep.CreateUser(&user); err != nil {
		t.Errorf("error was not expected while inserting stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestCreateUserDuplicated(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:       1,
		Login:    "login",
		Email:    "eamil@d.com",
		Password: pwd,
		FullName: "full name",
	}

	mock.ExpectExec("INSERT INTO user").
		WithArgs(user.Login, user.Email, user.Password, user.FullName).
		WillReturnError(sql.ErrNoRows)

	// now we execute our method
	if err := rep.CreateUser(&user); err == nil {
		t.Errorf("error was expected erro while inserting stats")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestGetUserByEmailCorrect(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	email := "eamil@d.com"

	user := repo.User{
		ID:          1,
		Login:       "login",
		Email:       email,
		Password:    pwd,
		FullName:    "full name",
		DateCreated: time.Now().Round(time.Second).UTC(),
		LastCheck:   time.Now().Round(time.Second).UTC(),
	}

	mock.ExpectQuery("SELECT (.+) FROM user WHERE").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "login", "email", "full_name", "password", "date_created", "last_check"}).
			AddRow(user.ID, user.Login, user.Email, user.FullName, user.Password, user.DateCreated, user.LastCheck))

	userRes, err := rep.GetUserWithEamil(email)
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	if !reflect.DeepEqual(userRes, user) {
		t.Errorf("expected %v got %v", user, userRes)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByEmailDontExist(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	email := "eamil@d.com"

	user := repo.User{}

	mock.ExpectQuery("SELECT (.+) FROM user WHERE").
		WithArgs(email).
		WillReturnError(fmt.Errorf("user doesn't exist"))

	userRes, err := rep.GetUserWithEamil(email)
	if err == nil {
		t.Errorf("error was expected erro while inserting stats")
	}

	if !reflect.DeepEqual(userRes, user) {
		t.Errorf("expected %v got %v", user, userRes)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRolesCorrect(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:    1,
		Roles: map[string]struct{}{"role1": {}, "role2": {}, "role3": {}},
	}

	mock.ExpectQuery("SELECT role.name from role JOIN user_has_role ON").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow("role1").AddRow("role2").AddRow("role3"))

	roles, err := rep.GetRoles(&user)
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	if !reflect.DeepEqual(roles, user.Roles) {
		t.Errorf("expected %v got %v", user.Roles, roles)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRolesWithError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:    1,
		Roles: nil,
	}

	mock.ExpectQuery("SELECT role.name from role JOIN user_has_role ON").
		WithArgs(user.ID).
		WillReturnError(fmt.Errorf("internal error"))

	roles, err := rep.GetRoles(&user)
	if err == nil {
		t.Errorf("error was expected erro while inserting stats")
	}

	if !reflect.DeepEqual(roles, user.Roles) {
		t.Errorf("expected %v got %v", user.Roles, roles)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRolesNoRoles(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:    1,
		Roles: map[string]struct{}{},
	}

	mock.ExpectQuery("SELECT role.name from role JOIN user_has_role ON").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"name"}))

	roles, err := rep.GetRoles(&user)
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	if !reflect.DeepEqual(roles, user.Roles) {
		t.Errorf("expected %v got %v", user.Roles, roles)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

const insert = "insert into user_has_role (.+) select (.+) from role where (.+)"

func TestSetRolesCorrect(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:    1,
		Roles: map[string]struct{}{"role1": {}, "role2": {}},
	}

	mock.ExpectBegin()

	mock.ExpectExec("DELETE FROM user_has_role WHERE (.+)").WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(insert).WithArgs(user.ID, "role1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(insert).WithArgs(user.ID, "role2").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := rep.SetRoles(&user)
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSetRolesInsertErr(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:    1,
		Roles: map[string]struct{}{"role1": {}, "role2": {}},
	}

	mock.ExpectBegin()

	mock.ExpectExec("DELETE FROM user_has_role WHERE (.+)").WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(insert).WithArgs(user.ID, "role1").WillReturnError(fmt.Errorf("insert error"))

	mock.ExpectRollback()

	err := rep.SetRoles(&user)
	if err == nil {
		t.Errorf("error was expected erro while inserting stats")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSetRolesDeletetErr(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()

	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rep := mysql.NewMysqlRepository(sqlxDB, logging.NewMockLogger())

	user := repo.User{
		ID:    1,
		Roles: map[string]struct{}{"role1": {}, "role2": {}},
	}

	mock.ExpectBegin()

	mock.ExpectExec("DELETE FROM user_has_role WHERE (.+)").WithArgs(user.ID).WillReturnError(fmt.Errorf("delete error"))

	mock.ExpectRollback()

	err := rep.SetRoles(&user)
	if err == nil {
		t.Errorf("error was expected erro while inserting stats")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
