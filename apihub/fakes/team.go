package fakes

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/apihub/apihub-cli/apihub"
)

func (fake *ApiHubServer) CreateTeam(w http.ResponseWriter, req *http.Request) {
	var team apihub.Team
	err := json.NewDecoder(req.Body).Decode(&team)
	if err != nil {
		panic(err)
	}

	team.Owner = "alice@example.org"
	team.Users = []string{team.Owner}

	if team.Name == "" {
		errorResponse := apihub.ErrorResponse{
			Type:        "bad_request",
			Description: "Name cannot be empty.",
		}
		fake.Error(w, http.StatusBadRequest, errorResponse)
		return
	}

	fake.Teams.Add(team)
	response, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (fake *ApiHubServer) UpdateTeam(w http.ResponseWriter, req *http.Request) {
	teamAlias := strings.TrimPrefix(req.URL.Path, "/api/teams/")

	teamFound, ok := fake.Teams.Get(teamAlias)
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}

	var team apihub.Team
	err := json.NewDecoder(req.Body).Decode(&team)
	if err != nil {
		panic(err)
	}
	team.Alias = teamFound.Alias

	fake.Teams.Add(team)
	response, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *ApiHubServer) TeamInfo(w http.ResponseWriter, req *http.Request) {
	teamAlias := strings.TrimPrefix(req.URL.Path, "/api/teams/")

	team, ok := fake.Teams.Get(teamAlias)
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}

	response, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *ApiHubServer) GetTeams(w http.ResponseWriter, req *http.Request) {
	teams := fake.Teams.List()

	collection := struct {
		Items []apihub.Team `json:"items"`
		Count int           `json:"item_count"`
	}{}
	collection.Items = teams
	collection.Count = len(teams)

	response, err := json.Marshal(collection)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *ApiHubServer) DeleteTeam(w http.ResponseWriter, req *http.Request) {
	teamAlias := strings.TrimPrefix(req.URL.Path, "/api/teams/")

	team, ok := fake.Teams.Get(teamAlias)
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}

	response, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	fake.Teams.Delete(teamAlias)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *ApiHubServer) AddUsersToTeam(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/api/teams/(.*)/users$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	team, ok := fake.Teams.Get(matches[1])
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}

	collection := struct {
		Users []string `json:"users"`
	}{}
	err := json.NewDecoder(req.Body).Decode(&collection)
	if err != nil {
		panic(err)
	}

	for _, u := range collection.Users {
		team.Users = append(team.Users, u)
	}

	fake.Teams.Add(team)
	response, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *ApiHubServer) RemoveUserFromTeam(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/api/teams/(.*)/users$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	team, ok := fake.Teams.Get(matches[1])
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}

	collection := struct {
		Users []string `json:"users"`
	}{}
	err := json.NewDecoder(req.Body).Decode(&collection)
	if err != nil {
		panic(err)
	}

	for _, u := range collection.Users {
		if i, ok := team.ContainsUserByEmail(u); ok {
			hi := len(team.Users) - 1
			if hi > i {
				team.Users[i] = team.Users[hi]
			}
			team.Users = team.Users[:hi]
		}
	}

	fake.Teams.Add(team)
	response, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
