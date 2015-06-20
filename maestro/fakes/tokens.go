package fakes

import (
	"github.com/backstage/backstage-cli/maestro"
)

type Tokens struct {
	storage map[string]backstage.User
}

func NewTokens() *Tokens {
	return &Tokens{
		storage: make(map[string]backstage.User),
	}
}

func (tokens *Tokens) Add(token string, user backstage.User) {
	tokens.storage[token] = user
}

func (tokens *Tokens) Reset() {
	tokens.storage = make(map[string]backstage.User)
}
