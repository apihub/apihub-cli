package backstage

import (
	"gopkg.in/mgo.v2/bson"
)

type Team struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string        `json:"name"`
	Alias    string        `json:"alias"`
	Users    []string      `json:"users"`
	Owner    string        `json:"owner"`
	Services []*Service    `json:"services,omitempty"`
	Clients  []*Client     `json:"clients,omitempty"`
}

func (t *Team) ContainsUserByEmail(email string) (int, bool) {
	for i, u := range t.Users {
		if u == email {
			return i, true
		}
	}

	return -1, false
}
