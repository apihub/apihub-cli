package fakes

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/backstage-cli/backstage"
)

func (fake *BackstageServer) CreateUser(w http.ResponseWriter, req *http.Request) {
	var user backstage.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	if user.Name == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		errorResponse := backstage.ErrorResponse{
			Type:        "bad_request",
			Description: "Name/Email/Username/Password cannot be empty.",
		}
		j, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	fake.Users.Add(user)
	response, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (fake *BackstageServer) DeleteUser(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
