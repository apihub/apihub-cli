package fakes

import (
	"encoding/json"
	"net/http"

	"github.com/backstage/backstage-cli/backstage"
)

func (fake *BackstageServer) Login(w http.ResponseWriter, req *http.Request) {
	var user backstage.User
	var j []byte
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	_, found := fake.Users.Get(user.Email)
	if !found {
		errorResponse := backstage.ErrorResponse{
			Type:        "bad_request",
			Description: "Authentication failed.",
		}
		j, _ = json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	token := backstage.TokenInfo{
		Type:      "Token",
		Token:     "RpOMQwiTMtxH6abgwonjBrVhBlrE1jbOxsk86UD_trI=",
		Expires:   86400,
		CreatedAt: "2015-05-29T01:05:45Z",
	}

	fake.Tokens.Add(token.Token, user)
	j, _ = json.Marshal(token)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (fake *BackstageServer) Logout(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (fake *BackstageServer) ChangePassword(w http.ResponseWriter, req *http.Request) {
	var user backstage.User
	var j []byte
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	if user.NewPassword != user.ConfirmationPassword {
		errorResponse := backstage.ErrorResponse{
			Type:        "bad_request",
			Description: "Your new password and confirmation password do not match.",
		}
		j, _ = json.Marshal(errorResponse)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(j)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
