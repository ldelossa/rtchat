package test

import (
	us "github.com/ldelossa/rtchat/userservice"
)

type GoodDataStore struct{}

func (g *GoodDataStore) Init() error {
	panic("not implemented")
}

func (g *GoodDataStore) AddUser(u us.User) error {
	panic("not implemented")
}

func (g *GoodDataStore) UpdateUser(u us.User) error {
	panic("not implemented")
}

func (g *GoodDataStore) GetUserByID(UID string) (*us.User, error) {
	panic("not implemented")
}

func (g *GoodDataStore) DeleteUserByID(UID string) error {
	panic("not implemented")
}
