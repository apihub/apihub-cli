package fakes

import (
	"net/http"
	"net/http/httptest"

	"github.com/backstage/backstage-client/backstage"
	"github.com/gorilla/mux"
)

type BackstageServer struct {
	server *httptest.Server

	Clients  *Clients
	Services *Services
	Teams    *Teams
	Tokens   *Tokens
	Users    *Users
}

func NewBackstageServer() *BackstageServer {
	fake := &BackstageServer{
		Clients:  NewClients(),
		Services: NewServices(),
		Teams:    NewTeams(),
		Tokens:   NewTokens(),
		Users:    NewUsers(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/login", fake.Login).Methods("POST")
	router.HandleFunc("/api/password", fake.ChangePassword).Methods("PUT")
	router.HandleFunc("/api/logout", fake.Logout).Methods("DELETE")
	router.HandleFunc("/api/users", fake.CreateUser).Methods("POST")
	router.HandleFunc("/api/users", fake.DeleteUser).Methods("DELETE")

	router.HandleFunc("/api/teams", fake.CreateTeam).Methods("POST")
	router.HandleFunc("/api/teams", fake.GetTeams).Methods("GET")
	router.HandleFunc("/api/teams/{alias}", fake.UpdateTeam).Methods("PUT")
	router.HandleFunc("/api/teams/{alias}", fake.TeamInfo).Methods("GET")
	router.HandleFunc("/api/teams/{alias}", fake.DeleteTeam).Methods("DELETE")
	router.HandleFunc("/api/teams/{alias}/users", fake.AddUsersToTeam).Methods("POST")
	router.HandleFunc("/api/teams/{alias}/users", fake.RemoveUserFromTeam).Methods("DELETE")

	router.HandleFunc("/api/clients", fake.CreateClient).Methods("POST")
	router.HandleFunc("/api/clients/{id}", fake.UpdateClient).Methods("PUT")
	router.HandleFunc("/api/clients/{id}", fake.DeleteClient).Methods("DELETE")
	router.HandleFunc("/api/clients/{id}", fake.ClientInfo).Methods("GET")

	router.HandleFunc("/api/services", fake.CreateService).Methods("POST")
	router.HandleFunc("/api/services/{subdomain}", fake.UpdateService).Methods("PUT")
	router.HandleFunc("/api/services/{subdomain}", fake.DeleteService).Methods("DELETE")
	fake.server = httptest.NewServer(router)
	return fake
}

func (fake *BackstageServer) Stop() {
	fake.server.Close()
}

func (fake *BackstageServer) URL() string {
	return fake.server.URL
}

func (fake *BackstageServer) Reset() {
	fake.Clients.Reset()
	fake.Services.Reset()
	fake.Teams.Reset()
	fake.Tokens.Reset()
	fake.Users.Reset()
}

func (fake *BackstageServer) notFound(w http.ResponseWriter, message string) {
	errorResponse := backstage.ErrorResponse{
		Type:        "not_found",
		Description: message,
	}
	fake.Error(w, http.StatusBadRequest, errorResponse)
}
