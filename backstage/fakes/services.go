package fakes

import (
	"github.com/backstage/backstage-cli/backstage"
)

type Services struct {
	storage map[string]backstage.Service
}

func NewServices() *Services {
	return &Services{
		storage: make(map[string]backstage.Service),
	}
}

func (services *Services) Add(service backstage.Service) {
	services.storage[service.Subdomain] = service
}

func (services *Services) Get(subdomain string) (backstage.Service, bool) {
	service, ok := services.storage[subdomain]
	return service, ok
}

func (services *Services) List() []backstage.Service {
	var ts []backstage.Service
	for _, t := range services.storage {
		ts = append(ts, t)
	}
	return ts
}

func (services *Services) Delete(subdomain string) {
	delete(services.storage, subdomain)
}

func (services *Services) Reset() {
	services.storage = make(map[string]backstage.Service)
}
