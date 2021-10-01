package main

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"gitlab.com/zephinzer/go-devops"
)

func main() {
	targetURL, err := url.Parse("https://httpbin.org/uuid")
	if err != nil {
		panic(err)
	}
	response, err := devops.SendHTTPRequest(devops.SendHTTPRequestOpts{
		URL: targetURL,
	})
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("uuid: %s", string(responseBody))
}
