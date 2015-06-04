package backstage

type Client struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Team        string `json:"team,omitempty"`
}
