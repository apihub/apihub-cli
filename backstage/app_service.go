package backstage

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppService struct {
	client HTTPClient
}

func NewAppService(client HTTPClient) *AppService {
	return &AppService{
		client: client,
	}
}

func (s AppService) Create(team, clientID, name, redirectURI, secret string) (App, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/api/apps",
		Body: App{
			Team:         team,
			ClientID:     clientID,
			Name:         name,
			RedirectURI:  redirectURI,
			ClientSecret: secret,
		},
	})

	if err != nil {
		return App{}, err
	}

	var app App
	err = json.Unmarshal(body, &app)
	if err != nil {
		panic(err)
	}

	return app, nil
}

func (s AppService) Update(team, clientID, name, redirectURI, secret string) (App, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           fmt.Sprintf("/api/apps/%s", clientID),
		Body: App{
			Team:         team,
			Name:         name,
			RedirectURI:  redirectURI,
			ClientSecret: secret,
		},
	})

	if err != nil {
		return App{}, err
	}

	var app App
	err = json.Unmarshal(body, &app)
	if err != nil {
		panic(err)
	}

	return app, nil
}

func (s AppService) Info(clientID string) (App, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "GET",
		Path:           fmt.Sprintf("/api/apps/%s", clientID),
	})

	if err != nil {
		return App{}, err
	}

	var app App
	err = json.Unmarshal(body, &app)
	if err != nil {
		panic(err)
	}

	return app, nil
}

func (s AppService) Delete(team, clientID string) error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           fmt.Sprintf("/api/apps/%s", clientID),
	})

	return err
}
