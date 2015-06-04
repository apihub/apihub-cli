package fakes

import (
	"encoding/json"
	"net/http"
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
	clientID := strings.TrimPrefix(req.URL.Path, "/api/clients/")

	clientFound, ok := fake.Clients.Get(clientID)
	if !ok {
		fake.notFound(w, "Client not found.")
		return
	}

	var client backstage.Client
	err := json.NewDecoder(req.Body).Decode(&client)
	if err != nil {
		panic(err)
	}

	client.ID = clientFound.ID
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
	id := strings.TrimPrefix(req.URL.Path, "/api/clients/")

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
	clientID := strings.TrimPrefix(req.URL.Path, "/api/clients/")

	fake.Clients.Delete(clientID)

	w.WriteHeader(http.StatusOK)
}
