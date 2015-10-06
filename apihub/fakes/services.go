package fakes

import (
	"github.com/apihub/apihub-cli/apihub"
)

type Services struct {
	storage map[string]apihub.Service
}

func NewServices() *Services {
	return &Services{
		storage: make(map[string]apihub.Service),
	}
}

func (services *Services) Add(service apihub.Service) {
	services.storage[service.Subdomain] = service
}

func (services *Services) Get(subdomain string) (apihub.Service, bool) {
	service, ok := services.storage[subdomain]
	return service, ok
}

func (services *Services) List() []apihub.Service {
	var ts []apihub.Service
	for _, t := range services.storage {
		ts = append(ts, t)
	}
	return ts
}

func (services *Services) Delete(subdomain string) {
	delete(services.storage, subdomain)
}

func (services *Services) Reset() {
	services.storage = make(map[string]apihub.Service)
}
