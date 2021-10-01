package devops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DownloadFileTests struct {
	suite.Suite
	Data     []byte
	Headers  map[string][]string
	Password string
	Username string
}

func TestDownloadFile(t *testing.T) {
	suite.Run(t, &DownloadFileTests{
		Data: []byte(fmt.Sprintf("%v", time.Now().UTC().UnixMicro())),
		Headers: map[string][]string{
			"a": {"b", "c"},
			"b": {"c", "d", "e"},
		},
		Password: "password",
		Username: "username",
	})
}

func (s DownloadFileTests) TestDownloadFile() {
	authKey := "auth"
	dataKey := "data"
	headersKey := "headers"
	usernameKey := "username"
	passwordKey := "password"

	testFilePath := "./tests/downloads/TestDownloadFile"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{}
		if username, password, ok := r.BasicAuth(); ok {
			response[authKey] = map[string]string{
				usernameKey: username,
				passwordKey: password,
			}
		}
		response[headersKey] = r.Header
		response[dataKey] = string(s.Data)
		responseBytes, err := json.MarshalIndent(response, "", "  ")
		s.Nil(err)
		w.Write(responseBytes)
	}))
	defer server.CloseClientConnections()
	defer server.Close()
	client := server.Client()
	serverURL, err := url.Parse(server.URL)
	s.Nil(err)

	options := DownloadFileOpts{
		BasicAuth: &BasicAuth{
			Username: s.Username,
			Password: s.Password,
		},
		Client:          client,
		DestinationPath: testFilePath,
		Headers:         s.Headers,
		Overwrite:       true,
		URL:             serverURL,
	}
	err = DownloadFile(options)
	s.Nil(err)
	fileContent, err := ioutil.ReadFile(testFilePath)
	s.Nil(err)
	var testResponse map[string]interface{}
	err = json.Unmarshal(fileContent, &testResponse)
	s.Nil(err)

	auth := testResponse[authKey]
	basicAuth := auth.(map[string]interface{})
	s.Equal(s.Username, basicAuth[usernameKey])
	s.Equal(s.Password, basicAuth[passwordKey])

	data := testResponse[dataKey].(string)
	s.Equal(string(s.Data), data)

	headers := testResponse[headersKey].(map[string]interface{})
	s.EqualValues(fmt.Sprintf("%s", s.Headers["a"]), fmt.Sprintf("%s", headers["A"]))
	s.EqualValues(fmt.Sprintf("%s", s.Headers["b"]), fmt.Sprintf("%s", headers["B"]))
}

func (s DownloadFileTests) TestDownloadFile_badUrlError() {
	targetURL, err := url.Parse("http://definitely.no.where.valid")
	s.Nil(err)
	options := DownloadFileOpts{
		DestinationPath: "./tests/downloads/TestDownloadFile_badUrlError",
		URL:             targetURL,
	}
	err = DownloadFile(options)
	s.NotNil(err)
	s.Contains(err.Error(), "failed to start download")
}

func (s DownloadFileTests) TestDownloadFile_directoryError() {
	targetURL, err := url.Parse("https://google.com")
	s.Nil(err)
	options := DownloadFileOpts{
		DestinationPath: "./tests/paths/to/a/directory",
		URL:             targetURL,
	}
	err = DownloadFile(options)
	s.NotNil(err)
	s.Contains(err.Error(), "it's a directory")
}

func (s DownloadFileTests) TestDownloadFile_overwriteError() {
	targetURL, err := url.Parse("https://google.com")
	s.Nil(err)
	options := DownloadFileOpts{
		DestinationPath: "./tests/paths/to/a/file",
		URL:             targetURL,
	}
	err = DownloadFile(options)
	s.NotNil(err)
	s.Contains(err.Error(), "refusing to overwrite")
}

func (s DownloadFileTests) TestDownloadFile_invalidOptions() {
	options := DownloadFileOpts{}
	err := DownloadFile(options)
	s.Contains(err.Error(), "missing destination file path")
	s.Contains(err.Error(), "missing url")

	options.URL = &url.URL{}
	err = DownloadFile(options)
	s.Contains(err.Error(), "missing host")
}
