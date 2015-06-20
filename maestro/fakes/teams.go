package fakes

import (
	"github.com/backstage/backstage-cli/maestro"
)

type Teams struct {
	storage map[string]backstage.Team
}

func NewTeams() *Teams {
	return &Teams{
		storage: make(map[string]backstage.Team),
	}
}

func (teams *Teams) Add(team backstage.Team) {
	teams.storage[team.Alias] = team
}

func (teams *Teams) Get(alias string) (backstage.Team, bool) {
	team, ok := teams.storage[alias]
	return team, ok
}

func (teams *Teams) List() []backstage.Team {
	var ts []backstage.Team
	for _, t := range teams.storage {
		ts = append(ts, t)
	}
	return ts
}

func (teams *Teams) Delete(alias string) {
	delete(teams.storage, alias)
}

func (teams *Teams) Reset() {
	teams.storage = make(map[string]backstage.Team)
}
