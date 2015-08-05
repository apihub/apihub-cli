package fakes

import (
	"github.com/apihub/apihub-cli/maestro"
)

type Teams struct {
	storage map[string]apihub.Team
}

func NewTeams() *Teams {
	return &Teams{
		storage: make(map[string]apihub.Team),
	}
}

func (teams *Teams) Add(team apihub.Team) {
	teams.storage[team.Alias] = team
}

func (teams *Teams) Get(alias string) (apihub.Team, bool) {
	team, ok := teams.storage[alias]
	return team, ok
}

func (teams *Teams) List() []apihub.Team {
	var ts []apihub.Team
	for _, t := range teams.storage {
		ts = append(ts, t)
	}
	return ts
}

func (teams *Teams) Delete(alias string) {
	delete(teams.storage, alias)
}

func (teams *Teams) Reset() {
	teams.storage = make(map[string]apihub.Team)
}
