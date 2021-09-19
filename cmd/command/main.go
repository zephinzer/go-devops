package main

import "devops"

func main() {
	ls, _ := devops.NewCommand(devops.NewCommandOpts{
		Command:   "ls",
		Arguments: []string{"-a", "-l"},
	})
	ls.Run()
}
