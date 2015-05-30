package backstage_test

import (
	"testing"

	"github.com/backstage/backstage-client/backstage"
	"github.com/backstage/backstage-client/backstage/fakes"
	. "gopkg.in/check.v1"
)

var backstageServer *fakes.BackstageServer
var httpClient backstage.HttpClient

var unsupportedPayload = func() {}

var authService *backstage.AuthService
var clientService *backstage.ClientService
var serviceService *backstage.ServiceService
var teamService *backstage.TeamService
var userService *backstage.UserService

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	backstageServer = fakes.NewBackstageServer()
}

func (s *S) SetUpTest(c *C) {
	backstageServer.Reset()

	httpClient = backstage.NewHttpClient(backstageServer.URL())
	authService = backstage.NewAuthService(httpClient)
	clientService = backstage.NewClientService(httpClient)
	serviceService = backstage.NewServiceService(httpClient)
	teamService = backstage.NewTeamService(httpClient)
	userService = backstage.NewUserService(httpClient)
}

func (s *S) TearDownSuite(c *C) {
	backstageServer.Stop()
}
