package fakes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/backstage/backstage-cli/backstage"
)

func (fake *BackstageServer) CreateService(w http.ResponseWriter, req *http.Request) {
	var service backstage.Service
	err := json.NewDecoder(req.Body).Decode(&service)
	if err != nil {

		panic(err)
	}

	_, ok := fake.Teams.Get(service.Team)
	if !ok {
		fake.notFound(w, "Team not found.")
		return
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
	subdomain := strings.TrimPrefix(req.URL.Path, "/api/services/")

	_, ok := fake.Services.Get(subdomain)
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
	subdomain := strings.TrimPrefix(req.URL.Path, "/api/services/")

	service, ok := fake.Services.Get(subdomain)
	if !ok {
		fake.notFound(w, "Service not found.")
		return
	}

	response, err := json.Marshal(service)
	if err != nil {
		panic(err)
	}

	fake.Services.Delete(service.Subdomain)

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
