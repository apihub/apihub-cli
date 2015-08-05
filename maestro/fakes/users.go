package fakes

import (
	"github.com/apihub/apihub-cli/maestro"
)

type Users struct {
	storage map[string]apihub.User
}

func NewUsers() *Users {
	return &Users{
		storage: make(map[string]apihub.User),
	}
}

func (users *Users) Add(user apihub.User) {
	users.storage[user.Email] = user
}

func (users *Users) Get(email string) (apihub.User, bool) {
	user, ok := users.storage[email]
	return user, ok
}

func (users *Users) Delete(email string) {
	delete(users.storage, email)
}

func (users *Users) Reset() {
	users.storage = make(map[string]apihub.User)
}
