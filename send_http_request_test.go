package devops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SendHTTPRequestTest struct {
	suite.Suite
	getHandler func(SendHTTPRequestTest) http.Handler
	observed   map[string]string
}

func TestSendHTTPRequest(t *testing.T) {
	suite.Run(t, &SendHTTPRequestTest{
		getHandler: func(s SendHTTPRequestTest) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				s.observed["method"] = r.Method
				s.observed["url"] = r.URL.String()

				if user, password, ok := r.BasicAuth(); ok {
					basicauthJSON, err := json.Marshal(map[string]string{
						"username": user,
						"password": password,
					})
					s.Nil(err, "basicauth should be writable")
					s.observed["basicauth"] = string(basicauthJSON)
				} else {
					s.observed["basicauth"] = "empty"
				}

				if r.Header != nil {
					headersJSON, err := json.Marshal(r.Header)
					s.Nil(err, "headers should be readable")
					s.observed["headers"] = string(headersJSON)
				} else {
					s.observed["headers"] = "empty"
				}

				if r.Body != nil {
					body, err := ioutil.ReadAll(r.Body)
					s.Nil(err, "body should be readable")
					s.observed["body"] = string(body)
				} else {
					s.observed["body"] = "empty"
				}
			})
		},
	})
}

func (s *SendHTTPRequestTest) BeforeTest(string, string) {
	s.observed = map[string]string{}
}

func (s SendHTTPRequestTest) TestSendHTTPRequestOpts_SetDefaults() {
	opts := SendHTTPRequestOpts{}
	opts.SetDefaults()

	s.NotNil(opts.Client)
	s.Equal(http.MethodGet, opts.Method)
}

func (s SendHTTPRequestTest) TestSendHTTPRequestOpts_Validate() {
	opts := SendHTTPRequestOpts{}
	err := opts.Validate()
	message := err.Error()
	s.Contains(message, "client")
	s.Contains(message, "method")
	s.Contains(message, "url")
	s.Contains(message, "failed to validate")

	opts.URL = &url.URL{}
	err = opts.Validate()
	message = err.Error()
	s.Contains(message, "host")
}

func (s SendHTTPRequestTest) TestSendHTTPRequest() {
	server := httptest.NewServer(s.getHandler(s))
	serverURL, err := url.Parse(server.URL)
	s.Nil(err)
	_, err = SendHTTPRequest(SendHTTPRequestOpts{
		BasicAuth: &BasicAuth{
			Username: "user",
			Password: "password",
		},
		Body: []byte("hello world"),
		Headers: map[string][]string{
			"one":   {"1", "uno", "une"},
			"two":   {"2", "dos", "deux"},
			"three": {"3", "tres", "trois"},
		},
		Method: http.MethodPost,
		URL:    serverURL,
	})
	s.Nil(err)

	var headers map[string][]string
	err = json.Unmarshal([]byte(s.observed["headers"]), &headers)
	s.Nil(err)
	s.EqualValues([]string{"1", "uno", "une"}, headers["One"])
	s.EqualValues([]string{"2", "dos", "deux"}, headers["Two"])
	s.EqualValues([]string{"3", "tres", "trois"}, headers["Three"])

	var basicauth map[string]string
	err = json.Unmarshal([]byte(s.observed["basicauth"]), &basicauth)
	s.Nil(err)
	s.Equal("user", basicauth["username"])
	s.Equal("password", basicauth["password"])
	s.Equal(http.MethodPost, s.observed["method"])
	s.Equal("/", s.observed["url"])
	s.Equal("hello world", s.observed["body"])
}

func (s SendHTTPRequestTest) TestSendHTTPRequest_noURL() {
	_, err := SendHTTPRequest(SendHTTPRequestOpts{})
	s.NotNil(err)
	s.Contains(err.Error(), "failed to send http request")
}
