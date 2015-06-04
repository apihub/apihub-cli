package backstage

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ClientService struct {
	client HTTPClient
}

func NewClientService(client HTTPClient) *ClientService {
	return &ClientService{
		client: client,
	}
}

func (s ClientService) Create(team, clientID, name, redirectURI, secret string) (Client, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/api/clients",
		Body: Client{
			Team:        team,
			ID:          clientID,
			Name:        name,
			RedirectURI: redirectURI,
			Secret:      secret,
		},
	})

	if err != nil {
		return Client{}, err
	}

	var client Client
	err = json.Unmarshal(body, &client)
	if err != nil {
		panic(err)
	}

	return client, nil
}

func (s ClientService) Update(team, clientID, name, redirectURI, secret string) (Client, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           fmt.Sprintf("/api/clients/%s", clientID),
		Body: Client{
			Team:        team,
			Name:        name,
			RedirectURI: redirectURI,
			Secret:      secret,
		},
	})

	if err != nil {
		return Client{}, err
	}

	var client Client
	err = json.Unmarshal(body, &client)
	if err != nil {
		panic(err)
	}

	return client, nil
}

func (s ClientService) Info(clientID string) (Client, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "GET",
		Path:           fmt.Sprintf("/api/clients/%s", clientID),
	})

	if err != nil {
		return Client{}, err
	}

	var client Client
	err = json.Unmarshal(body, &client)
	if err != nil {
		panic(err)
	}

	return client, nil
}

func (s ClientService) Delete(team, clientID string) error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           fmt.Sprintf("/api/clients/%s", clientID),
	})

	return err
}
