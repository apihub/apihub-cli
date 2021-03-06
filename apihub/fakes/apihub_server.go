package fakes

import (
	"net/http"
	"net/http/httptest"

	"github.com/apihub/apihub-cli/apihub"
	"github.com/gorilla/mux"
)

type ApiHubServer struct {
	server *httptest.Server

	Apps     *Apps
	Services *Services
	Teams    *Teams
	Tokens   *Tokens
	Users    *Users
}

func NewApiHubServer() *ApiHubServer {
	fake := &ApiHubServer{
		Apps:     NewApps(),
		Services: NewServices(),
		Teams:    NewTeams(),
		Tokens:   NewTokens(),
		Users:    NewUsers(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/auth/login", fake.Login).Methods("POST")
	router.HandleFunc("/auth/password", fake.ChangePassword).Methods("PUT")
	router.HandleFunc("/auth/login", fake.Logout).Methods("DELETE")
	router.HandleFunc("/auth/signup", fake.CreateUser).Methods("POST")
	router.HandleFunc("/auth/signup", fake.DeleteUser).Methods("DELETE")

	router.HandleFunc("/api/teams", fake.CreateTeam).Methods("POST")
	router.HandleFunc("/api/teams", fake.GetTeams).Methods("GET")
	router.HandleFunc("/api/teams/{alias}", fake.UpdateTeam).Methods("PUT")
	router.HandleFunc("/api/teams/{alias}", fake.TeamInfo).Methods("GET")
	router.HandleFunc("/api/teams/{alias}", fake.DeleteTeam).Methods("DELETE")
	router.HandleFunc("/api/teams/{alias}/users", fake.AddUsersToTeam).Methods("PUT")
	router.HandleFunc("/api/teams/{alias}/users", fake.RemoveUserFromTeam).Methods("DELETE")

	router.HandleFunc("/api/apps", fake.CreateApp).Methods("POST")
	router.HandleFunc("/api/apps/{id}", fake.UpdateApp).Methods("PUT")
	router.HandleFunc("/api/apps/{id}", fake.DeleteApp).Methods("DELETE")
	router.HandleFunc("/api/apps/{id}", fake.AppInfo).Methods("GET")

	router.HandleFunc("/api/services", fake.CreateService).Methods("POST")
	router.HandleFunc("/api/services/{subdomain}", fake.UpdateService).Methods("PUT")
	router.HandleFunc("/api/services/{subdomain}", fake.DeleteService).Methods("DELETE")
	fake.server = httptest.NewServer(router)
	return fake
}

func (fake *ApiHubServer) Stop() {
	fake.server.Close()
}

func (fake *ApiHubServer) URL() string {
	return fake.server.URL
}

func (fake *ApiHubServer) Reset() {
	fake.Apps.Reset()
	fake.Services.Reset()
	fake.Teams.Reset()
	fake.Tokens.Reset()
	fake.Users.Reset()
}

func (fake *ApiHubServer) notFound(w http.ResponseWriter, message string) {
	errorResponse := apihub.ErrorResponse{
		Type:        "not_found",
		Description: message,
	}
	fake.Error(w, http.StatusBadRequest, errorResponse)
}
