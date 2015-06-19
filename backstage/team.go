package backstage

type Team struct {
	Name     string     `json:"name"`
	Alias    string     `json:"alias"`
	Users    []string   `json:"users"`
	Owner    string     `json:"owner"`
	Services []*Service `json:"services,omitempty"`
	Apps     []*App     `json:"apps,omitempty"`
}

func (t *Team) ContainsUserByEmail(email string) (int, bool) {
	for i, u := range t.Users {
		if u == email {
			return i, true
		}
	}

	return -1, false
}
