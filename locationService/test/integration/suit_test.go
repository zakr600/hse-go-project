package integration

import (
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/config"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/repository"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(ServerSuit))
}

func (s *ServerSuit) TestServer() {

	c := http.Client{}
	r, err := c.Get(fmt.Sprintf("http://localhost:%d/drivers?lat=0.1&lng=0.2&radius=0.5", s.serverCfg.Port))
	assert.Equal(s.T(), err, nil)
	assert.Equal(s.T(), r.StatusCode, http.StatusNotFound)

	r, err = c.Post(fmt.Sprintf("http://localhost:%d/drivers/driver1/location", s.serverCfg.Port),
		"application/json",
		strings.NewReader("{\"lat\": 0.1,\n\"lng\": 0.2}"))
	assert.Equal(s.T(), err, nil)
	assert.Equal(s.T(), r.StatusCode, http.StatusOK)

	r, err = c.Get(fmt.Sprintf("http://localhost:%d/drivers?lat=0.1&lng=0.2&radius=0.5", s.serverCfg.Port))
	assert.Equal(s.T(), err, nil)
	assert.Equal(s.T(), r.StatusCode, http.StatusOK)

	r, err = c.Get(fmt.Sprintf("http://localhost:%d/drivers?lat=10&lng=20&radius=0.5", s.serverCfg.Port))
	assert.Equal(s.T(), err, nil)
	assert.Equal(s.T(), r.StatusCode, http.StatusNotFound)
}

type ServerSuit struct {
	suite.Suite
	serverCfg config.ServerConfig
	server    *internal.App
}

func (s *ServerSuit) SetupSuite() {
	s.serverCfg = config.ServerConfig{Port: 8080, ApiVersion: "1.0", Debug: true}
	serverApp := internal.NewApplication(s.serverCfg, service.CreateMainService(repository.CreateMapRepository()))
	go func() {
		serverApp.Run()
	}()
	s.server = serverApp

	s.T().Log("Suite setup is done")
}

func (s *ServerSuit) TearDownSuite() {
	s.T().Log("Suite stop is done")
}

func (s *ServerSuit) BeforeTest(suiteName, testName string) {}

func (s *ServerSuit) AfterTest(suiteName, testName string) {}
