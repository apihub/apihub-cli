package apihub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TeamService struct {
	client HTTPClient
}

func NewTeamService(client HTTPClient) *TeamService {
	return &TeamService{
		client: client,
	}
}

func (s TeamService) Create(name, alias string) (Team, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/api/teams",
		Body: Team{
			Name:  name,
			Alias: alias,
		},
	})

	if err != nil {
		return Team{}, err
	}

	var team Team
	err = json.Unmarshal(body, &team)
	if err != nil {
		panic(err)
	}

	return team, nil
}

func (s TeamService) Update(name, alias string) (Team, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           fmt.Sprintf("/api/teams/%s", alias),
		Body: Team{
			Name:  name,
			Alias: alias,
		},
	})

	if err != nil {
		return Team{}, err
	}

	var team Team
	err = json.Unmarshal(body, &team)
	if err != nil {
		panic(err)
	}

	return team, nil
}

func (s TeamService) Info(alias string) (Team, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "GET",
		Path:           fmt.Sprintf("/api/teams/%s", alias),
	})

	if err != nil {
		return Team{}, err
	}

	var team Team
	err = json.Unmarshal(body, &team)
	if err != nil {
		panic(err)
	}

	return team, nil
}

func (s TeamService) List() ([]Team, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "GET",
		Path:           "/api/teams",
	})

	if err != nil {
		return []Team{}, err
	}

	collection := struct {
		Items []Team `json:"items"`
		Count int    `json:"item_count"`
	}{}

	err = json.Unmarshal(body, &collection)
	if err != nil {
		panic(err)
	}

	return collection.Items, nil
}

func (s TeamService) Delete(alias string) error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           fmt.Sprintf("/api/teams/%s", alias),
	})

	return err
}

func (s TeamService) AddUser(alias, email string) (bool, error) {
	users := struct {
		Users []string `json:"users"`
	}{Users: []string{email}}

	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "POST",
		Path:           fmt.Sprintf("/api/teams/%s/users", alias),
		Body:           users,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s TeamService) RemoveUser(alias, email string) (bool, error) {
	users := struct {
		Users []string `json:"users"`
	}{Users: []string{email}}

	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           fmt.Sprintf("/api/teams/%s/users", alias),
		Body:           users,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
