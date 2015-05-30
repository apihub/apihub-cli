package fakes

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/backstage/backstage-client/backstage"
)

func (fake *BackstageServer) CreateClient(w http.ResponseWriter, req *http.Request) {
	var client backstage.Client
	err := json.NewDecoder(req.Body).Decode(&client)
	if err != nil {
		panic(err)
	}

	if client.Name == "" {
		errorResponse := backstage.ErrorResponse{
			Type:        "bad_request",
			Description: "Name cannot be empty.",
		}
		j, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	fake.Clients.Add(client)
	response, err := json.Marshal(client)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (fake *BackstageServer) UpdateClient(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/api/teams/(.*)/clients/(.*)$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	_, ok := fake.Teams.Get(matches[1])
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}

	clientFound, ok := fake.Clients.Get(matches[2])
	if !ok {
		fake.notFound(w, "Client not found.")
		return
	}

	var client backstage.Client
	err := json.NewDecoder(req.Body).Decode(&client)
	if err != nil {
		panic(err)
	}

	client.Id = clientFound.Id
	client.Team = clientFound.Team
	fake.Clients.Add(client)

	response, err := json.Marshal(client)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *BackstageServer) ClientInfo(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/teams/clients/")

	client, ok := fake.Clients.Get(id)
	if !ok {
		fake.notFound(w, "Client not found.")
		return
	}

	response, err := json.Marshal(client)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *BackstageServer) DeleteClient(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/api/teams/.*/clients/(.*)$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	fake.Clients.Delete(matches[1])

	w.WriteHeader(http.StatusOK)
}
