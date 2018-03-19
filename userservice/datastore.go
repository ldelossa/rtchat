package userservice

type DataStore interface {
	Init() error
	AddUser(u User) error
	UpdateUser(u User) error
	GetUserByID(UID string) (*User, error)
	GetUserByUserName(un string) (*User, error)
	DeleteUserByID(UID string) error
}
