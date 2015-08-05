package apihub_test

import (
	"testing"

	"github.com/apihub/apihub-cli/maestro"
	"github.com/apihub/apihub-cli/maestro/fakes"
	. "gopkg.in/check.v1"
)

var apihubServer *fakes.ApiHubServer
var httpClient apihub.HTTPClient

var unsupportedPayload = func() {}

var authService *apihub.AuthService
var appService *apihub.AppService
var serviceService *apihub.ServiceService
var teamService *apihub.TeamService
var userService *apihub.UserService

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	apihubServer = fakes.NewApiHubServer()
}

func (s *S) SetUpTest(c *C) {
	apihubServer.Reset()

	httpClient = apihub.NewHTTPClient(apihubServer.URL())
	authService = apihub.NewAuthService(httpClient)
	appService = apihub.NewAppService(httpClient)
	serviceService = apihub.NewServiceService(httpClient)
	teamService = apihub.NewTeamService(httpClient)
	userService = apihub.NewUserService(httpClient)
}

func (s *S) TearDownSuite(c *C) {
	apihubServer.Stop()
}
