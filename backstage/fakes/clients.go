package fakes

import (
	"github.com/backstage/backstage-client/backstage"
)

type Clients struct {
	storage map[string]backstage.Client
}

func NewClients() *Clients {
	return &Clients{
		storage: make(map[string]backstage.Client),
	}
}

func (clients *Clients) Add(client backstage.Client) {
	clients.storage[client.Id] = client
}

func (clients *Clients) Get(id string) (backstage.Client, bool) {
	client, ok := clients.storage[id]
	return client, ok
}

func (clients *Clients) Delete(id string) {
	delete(clients.storage, id)
}

func (clients *Clients) Reset() {
	clients.storage = make(map[string]backstage.Client)
}
