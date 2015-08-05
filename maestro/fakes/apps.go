package fakes

import (
	"github.com/apihub/apihub-cli/maestro"
)

type Apps struct {
	storage map[string]apihub.App
}

func NewApps() *Apps {
	return &Apps{
		storage: make(map[string]apihub.App),
	}
}

func (apps *Apps) Add(app apihub.App) {
	apps.storage[app.ClientID] = app
}

func (apps *Apps) Get(id string) (apihub.App, bool) {
	app, ok := apps.storage[id]
	return app, ok
}

func (apps *Apps) Delete(id string) {
	delete(apps.storage, id)
}

func (apps *Apps) Reset() {
	apps.storage = make(map[string]apihub.App)
}
