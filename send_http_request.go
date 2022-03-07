package devops

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SendHTTPRequestOpts presents options for the SendHTTPRequest
// method
type SendHTTPRequestOpts struct {
	// BasicAuth defines user credentials for use with the
	// request similar to curl's --basic flag. If left nil,
	// basic auth will not be used
	BasicAuth *BasicAuth

	// Body defines the body data to be sent with the request.
	// If left nil, a request without a body will be sent
	Body []byte

	// Client defines the HTTP client to use. If left nil,
	// defaults to http.DefaultClient
	Client *http.Client

	// Headers defines the headers to be sent along with this
	// request. If left nil, no headers will be sent
	Headers map[string][]string

	// Method defines the HTTP method to make the request with
	Method string

	// URL defines the endpoint to call
	URL *url.URL
}

// SetDefaults sets defaults for the options object instance
func (o *SendHTTPRequestOpts) SetDefaults() {
	if o.Client == nil {
		o.Client = http.DefaultClient
	}

	if o.Method == "" {
		o.Method = http.MethodGet
	}
}

// Validate validates the options to check if this object
// instance is usable by SendHTTPRequest
func (o SendHTTPRequestOpts) Validate() error {
	errors := []string{}

	if o.Client == nil {
		errors = append(errors, "missing client")
	}

	if o.Method == "" {
		errors = append(errors, "missing method")
	}

	if o.URL == nil {
		errors = append(errors, "missing url")
	} else if o.URL.Host == "" {
		errors = append(errors, "missing host in url")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}
	return nil
}

// SendHTTPRequest performs a HTTP request as configured by the provided
// options object instance `opts`.
func SendHTTPRequest(opts SendHTTPRequestOpts) (*http.Response, error) {
	opts.SetDefaults()
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("failed to send http request: %s", err)
	}
	if opts.BasicAuth != nil {
		opts.URL.User = url.UserPassword(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	var body io.Reader
	var req *http.Request
	if opts.Body != nil {
		body = bytes.NewReader(opts.Body)
	}
	req, err := http.NewRequest(opts.Method, opts.URL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request object: %s", err)
	}
	if opts.Headers != nil {
		req.Header = opts.Headers
	}
	res, err := opts.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to start download from '%s': %s", opts.URL.String(), err)
	}
	return res, nil
}
