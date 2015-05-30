package backstage

import (
	"encoding/json"
	"net/http"
)

type UserService struct {
	client HttpClient
}

func NewUserService(client HttpClient) *UserService {
	return &UserService{
		client: client,
	}
}

func (s UserService) Create(name, username, email, password string) (User, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/api/users",
		Body: User{
			Name:     name,
			Username: username,
			Email:    email,
			Password: password,
		},
	})

	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	return user, nil
}

func (s UserService) Delete() error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           "/api/users",
	})

	DeleteToken()
	return err
}