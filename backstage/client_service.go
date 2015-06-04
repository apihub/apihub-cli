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

func (s ClientService) Create(team, clientId, name, redirectUri, secret string) (Client, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/api/clients",
		Body: Client{
			Team:        team,
			ID:          clientId,
			Name:        name,
			RedirectURI: redirectUri,
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

func (s ClientService) Update(team, clientId, name, redirectUri, secret string) (Client, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           fmt.Sprintf("/api/clients/%s", clientId),
		Body: Client{
			Team:        team,
			Name:        name,
			RedirectURI: redirectUri,
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

func (s ClientService) Info(clientId string) (Client, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "GET",
		Path:           fmt.Sprintf("/api/clients/%s", clientId),
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

func (s ClientService) Delete(team, clientId string) error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           fmt.Sprintf("/api/clients/%s", clientId),
	})

	return err
}
