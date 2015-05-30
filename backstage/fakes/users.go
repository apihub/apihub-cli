package fakes

import (
	"github.com/backstage/backstage-client/backstage"
)

type Users struct {
	storage map[string]backstage.User
}

func NewUsers() *Users {
	return &Users{
		storage: make(map[string]backstage.User),
	}
}

func (users *Users) Add(user backstage.User) {
	users.storage[user.Email] = user
}

func (users *Users) Get(email string) (backstage.User, bool) {
	user, ok := users.storage[email]
	return user, ok
}

func (users *Users) Delete(email string) {
	delete(users.storage, email)
}

func (users *Users) Reset() {
	users.storage = make(map[string]backstage.User)
}
