package backstage

type App struct {
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
	Name         string   `json:"name,omitempty"`
	RedirectURIs []string `json:"redirect_uris,omitempty"`
	Team         string   `json:"team,omitempty"`
}
