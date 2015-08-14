package interfaces

import (
	"errors"

	"github.com/gh-service/domain"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

var counter = 0
var users = []domain.User{}

func (repo UserRepo) Store(user domain.User) error {
	user.ID = int64(counter)
	counter++
	users = append(users, user)
	return nil
}

func (repo UserRepo) RetrieveByID(id int64) (*domain.User, error) {

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("User not found")
}

func (repo UserRepo) RetrieveByUserName(username string) (*domain.User, error) {

	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("User not found")
}

func (repo UserRepo) Update(incUser domain.User) error {

	for index, user := range users {
		for user.ID == incUser.ID {
			users[index] = incUser
			return nil
		}
	}
	return errors.New("User not found")
}
