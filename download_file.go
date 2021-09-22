package devops

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DownloadStatusUpdate struct {
	Completed bool
	Progress  float32
}

type BasicAuthOpts struct {
	Username string
	Password string
}

type DownloadFileOpts struct {
	BasicAuth       *BasicAuthOpts
	Client          *http.Client
	Headers         map[string][]string
	DestinationPath string
	Overwrite       bool
	URL             *url.URL
}

func (o *DownloadFileOpts) SetDefaults() {
	if o.Client == nil {
		o.Client = &http.Client{}
	}
}

func (o DownloadFileOpts) Validate() error {
	errors := []string{}

	if o.DestinationPath == "" {
		errors = append(errors, "missing destination file path")
	}

	if o.URL.Host == "" {
		errors = append(errors, "missing host")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}

	return nil
}

func DownloadFile(opts DownloadFileOpts) error {
	opts.SetDefaults()
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("failed to download file: %s", err)
	}

	fileDestination, err := NormalizeLocalPath(opts.DestinationPath)
	if err != nil {
		return fmt.Errorf("failed to normalize path '%s': %s", opts.DestinationPath, err)
	}

	fileInfo, err := os.Lstat(fileDestination)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to access path '%s': %s", fileDestination, err)
		}
	}
	if err == nil {
		if fileInfo.IsDir() {
			return fmt.Errorf("failed to get a file at '%s': it's a directory", fileDestination)
		}
		if !opts.Overwrite {
			return fmt.Errorf("refusing to overwrite file at '%s' (set .Overwrite to true)", fileDestination)
		}
	}

	fileHandle, err := os.OpenFile(fileDestination, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file at '%s': %s", fileDestination, err)
	}
	defer fileHandle.Close()

	if opts.BasicAuth != nil {
		opts.URL.User = url.UserPassword(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    opts.URL,
	}
	if opts.Headers != nil {
		req.Header = opts.Headers
	}
	res, err := opts.Client.Do(&req)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to start download from '%s': %s", opts.URL.String(), err)
	}
	_, err = io.Copy(fileHandle, res.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file at '%s': %s", fileDestination, err)
	}

	return nil
}
