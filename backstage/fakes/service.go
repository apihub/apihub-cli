package fakes

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/backstage/backstage-client/backstage"
)

func (fake *BackstageServer) CreateService(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/api/teams/(.*)/services$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	_, ok := fake.Teams.Get(matches[1])
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}

	var service backstage.Service
	err := json.NewDecoder(req.Body).Decode(&service)
	if err != nil {

		panic(err)
	}

	if service.Subdomain == "" || service.Endpoint == "" {
		errorResponse := backstage.ErrorResponse{
			Type:        "bad_request",
			Description: "Subdomain/Endpoint cannot be empty.",
		}
		j, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	fake.Services.Add(service)
	response, err := json.Marshal(service)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (fake *BackstageServer) UpdateService(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/api/teams/(.*)/services/(.*)$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	_, ok := fake.Teams.Get(matches[1])
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}
	_, ok = fake.Services.Get(matches[2])
	if !ok {
		fake.notFound(w, "Service not found.")
		return
	}

	var service backstage.Service
	err := json.NewDecoder(req.Body).Decode(&service)
	if err != nil {
		panic(err)
	}

	if service.Subdomain == "" || service.Endpoint == "" {
		errorResponse := backstage.ErrorResponse{
			Type:        "bad_request",
			Description: "Subdomain/Endpoint cannot be empty.",
		}
		j, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	fake.Services.Add(service)

	response, err := json.Marshal(service)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *BackstageServer) DeleteService(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/api/teams/(.*)/services/(.*)$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	team, ok := fake.Teams.Get(matches[1])
	if !ok {
		fake.notFound(w, "Team not found.")
		return
	}
	service, ok := fake.Services.Get(matches[2])
	if !ok {
		fake.notFound(w, "Service not found.")
		return
	}

	response, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	fake.Services.Delete(service.Subdomain)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
