package devops

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidateConnectionTest struct {
	suite.Suite
}

func TestValidateConnection(t *testing.T) {
	suite.Run(t, &ValidateApplicationsTest{})
}

func (s ValidateApplicationsTest) TestValidateConnection() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer server.CloseClientConnections()
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	s.Nil(err)
	host := serverURL.Hostname()
	port, err := strconv.Atoi(serverURL.Port())
	s.Nil(err)
	validateConnectionOpts := ValidateConnectionOpts{
		Hostname: host,
		Port:     uint16(port),
	}
	ok, err := ValidateConnection(validateConnectionOpts)
	s.True(ok)
	s.Nil(err)
}
