package main

import (
	"log"
	"regexp"

	"gitlab.com/zephinzer/go-devops"
)

func main() {
	yes, err := devops.Confirm(devops.ConfirmOpts{
		Question:   "exact match",
		MatchExact: "yes",
	})
	if err != nil {
		log.Fatalf("failed to get user input: %s", err)
	}
	log.Printf("user confirmed: %v\n", yes)

	yes, err = devops.Confirm(devops.ConfirmOpts{
		Question:    "regexp match",
		MatchRegexp: regexp.MustCompile("^.+$"),
	})
	if err != nil {
		log.Fatalf("failed to get user input: %s", err)
	}
	log.Printf("user confirmed: %v\n", yes)
}
