package main

import "gitlab.com/zephinzer/go-devops"

func main() {
	ls, _ := devops.NewCommand(devops.NewCommandOpts{
		Command:   "ls",
		Arguments: []string{"-a", "-l"},
	})
	ls.Run()
}
