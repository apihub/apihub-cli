package backstage

type Service struct {
	Description   string   `json:"description,omitempty"`
	Disabled      bool     `json:"disabled,omitempty"`
	Documentation string   `json:"documentation,omitempty"`
	Endpoint      string   `json:"endpoint,omitempty"`
	Owner         string   `json:"owner,omitempty"`
	Subdomain     string   `json:"subdomain,omitempty"`
	Team          string   `json:"team,omitempty"`
	Timeout       int64    `json:"timeout,omitempty"`
	Transformers  []string `json:"transformers,omitempty"`
}
