package backstage

type Client struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	RedirectUri string `json:"redirect_uri,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Team        string `json:"team,omitempty"`
}
