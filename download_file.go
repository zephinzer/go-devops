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

// DownloadFileOpts presents configuration for the
// DownloadFile method
type DownloadFileOpts struct {
	BasicAuth       *BasicAuth
	Client          *http.Client
	Headers         map[string][]string
	DestinationPath string
	Overwrite       bool
	URL             *url.URL
}

// SetDefaults sets defaults for this object instance
func (o *DownloadFileOpts) SetDefaults() {
	if o.Client == nil {
		o.Client = &http.Client{}
	}
}

// Validate verifies that this object instance is usable
// by the DownloadFile method
func (o DownloadFileOpts) Validate() error {
	errors := []string{}

	if o.DestinationPath == "" {
		errors = append(errors, "missing destination file path")
	}

	if o.URL == nil {
		errors = append(errors, "missing url")
	} else if o.URL.Host == "" {
		errors = append(errors, "missing host")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}

	return nil
}

// DownloadFile downloads a file as directed by the configuration set
// in the options object `opts`
func DownloadFile(opts DownloadFileOpts) (err error) {
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

	res, err := SendHTTPRequest(SendHTTPRequestOpts{
		BasicAuth: opts.BasicAuth,
		Headers:   opts.Headers,
		Method:    http.MethodGet,
		URL:       opts.URL,
	})
	if err != nil {
		return fmt.Errorf("failed to make request: %s", err)
	}
	defer res.Body.Close()

	/* #nosec - this is required to write the file */
	fileHandle, err := os.OpenFile(fileDestination, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file at '%s': %s", fileDestination, err)
	}
	defer func() {
		if e := fileHandle.Close(); e != nil {
			closeError := fmt.Errorf("failed to close file at '%s': %s", fileDestination, e)
			if err != nil {
				err = fmt.Errorf("%s (previous error: %s)", closeError.Error(), err.Error())
			} else {
				err = closeError
			}
		}
	}()
	_, err = io.Copy(fileHandle, res.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file at '%s': %s", fileDestination, err)
	}

	return nil
}
