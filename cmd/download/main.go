package main

import (
	"devops"
	"fmt"
	"net/url"
	"time"
)

func main() {
	targetURL, err := url.Parse("https://google.com")
	if err != nil {
		panic(err)
	}
	now := time.Now().UTC().UnixMicro()
	if err := devops.DownloadFile(devops.DownloadFileOpts{
		DestinationPath: fmt.Sprintf("./tests/downloads/%v.txt", now),
		URL:             targetURL,
	}); err != nil {
		panic(err)
	}
}
