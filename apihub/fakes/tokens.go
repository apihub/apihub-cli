package fakes

import (
	"github.com/apihub/apihub-cli/apihub"
)

type Tokens struct {
	storage map[string]apihub.User
}

func NewTokens() *Tokens {
	return &Tokens{
		storage: make(map[string]apihub.User),
	}
}

func (tokens *Tokens) Add(token string, user apihub.User) {
	tokens.storage[token] = user
}

func (tokens *Tokens) Reset() {
	tokens.storage = make(map[string]apihub.User)
}
