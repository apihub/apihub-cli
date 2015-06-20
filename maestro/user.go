package backstage

type User struct {
	Name                 string `json:"name,omitempty"`
	Email                string `json:"email,omitempty"`
	Username             string `json:"username,omitempty"`
	Password             string `json:"password,omitempty"`
	NewPassword          string `json:"new_password,omitempty"`
	ConfirmationPassword string `json:"confirmation_password,omitempty"`
}
