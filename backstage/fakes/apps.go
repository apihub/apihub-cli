package fakes

import (
	"github.com/backstage/backstage-cli/backstage"
)

type Apps struct {
	storage map[string]backstage.App
}

func NewApps() *Apps {
	return &Apps{
		storage: make(map[string]backstage.App),
	}
}

func (apps *Apps) Add(app backstage.App) {
	apps.storage[app.ClientID] = app
}

func (apps *Apps) Get(id string) (backstage.App, bool) {
	app, ok := apps.storage[id]
	return app, ok
}

func (apps *Apps) Delete(id string) {
	delete(apps.storage, id)
}

func (apps *Apps) Reset() {
	apps.storage = make(map[string]backstage.App)
}
