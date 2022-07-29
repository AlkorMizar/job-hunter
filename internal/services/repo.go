package services

type Repository interface {
	UserManagment
	Close() error
}
