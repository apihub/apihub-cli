package fakes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/apihub/apihub-cli/maestro"
)

func (fake *ApiHubServer) CreateApp(w http.ResponseWriter, req *http.Request) {
	var app apihub.App
	err := json.NewDecoder(req.Body).Decode(&app)
	if err != nil {
		panic(err)
	}

	if app.Name == "" {
		errorResponse := apihub.ErrorResponse{
			Type:        "bad_request",
			Description: "Name cannot be empty.",
		}
		j, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	fake.Apps.Add(app)
	response, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (fake *ApiHubServer) UpdateApp(w http.ResponseWriter, req *http.Request) {
	appID := strings.TrimPrefix(req.URL.Path, "/api/apps/")

	appFound, ok := fake.Apps.Get(appID)
	if !ok {
		fake.notFound(w, "App not found.")
		return
	}

	var app apihub.App
	err := json.NewDecoder(req.Body).Decode(&app)
	if err != nil {
		panic(err)
	}

	app.ClientID = appFound.ClientID
	app.Team = appFound.Team
	fake.Apps.Add(app)

	response, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *ApiHubServer) AppInfo(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/apps/")

	app, ok := fake.Apps.Get(id)
	if !ok {
		fake.notFound(w, "App not found.")
		return
	}

	response, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (fake *ApiHubServer) DeleteApp(w http.ResponseWriter, req *http.Request) {
	appID := strings.TrimPrefix(req.URL.Path, "/api/apps/")

	fake.Apps.Delete(appID)

	w.WriteHeader(http.StatusOK)
}
