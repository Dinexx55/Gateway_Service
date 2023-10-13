package repository

import (
	"GatewayService/internal/service"
	"errors"
)

type MockUserRepository struct {
	users []service.User
}

func NewMockUserRepository() *MockUserRepository {
	repo := &MockUserRepository{
		users: []service.User{
			{Login: "user1", Password: "password1"},
			{Login: "user2", Password: "password2"},
			{Login: "user3", Password: "password3"},
		},
	}
	return repo
}

func (r *MockUserRepository) GetUserByLogin(login string) (*service.User, error) {
	for _, user := range r.users {
		if user.Login == login {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}
