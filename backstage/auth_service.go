package backstage

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthService struct {
	client HTTPClient
}

func NewAuthService(client HTTPClient) *AuthService {
	return &AuthService{
		client: client,
	}
}

func (s AuthService) Login(email, password string) (TokenInfo, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "POST",
		Path:           "/auth/login",
		Body: User{
			Email:    email,
			Password: password,
		},
	})

	if err != nil {
		return TokenInfo{}, err
	}

	var token TokenInfo
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}

	if err = WriteToken(fmt.Sprintf("%s %s", token.Type, token.Token)); err != nil {
		return TokenInfo{}, err
	}

	return token, nil
}

func (s AuthService) Logout() error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusNoContent,
		Method:         "DELETE",
		Path:           "/auth/login",
	})

	filesystem().Remove(TokenFileName)
	return err
}

func (s AuthService) ChangePassword(email, password, newPassword, confirmationPassword string) error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusNoContent,
		Method:         "PUT",
		Path:           "/api/password",
		Body: User{
			Email:                email,
			Password:             password,
			NewPassword:          newPassword,
			ConfirmationPassword: confirmationPassword,
		},
	})

	return err
}
