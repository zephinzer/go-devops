package devops

import (
	"fmt"
	"os/exec"
	"strings"
)

type ValidateApplicationsErrors struct {
	Errors []string
}

func (e *ValidateApplicationsErrors) Push(err string) {
	e.Errors = append(e.Errors, err)
}

func (e ValidateApplicationsErrors) Len() int {
	return len(e.Errors)
}

func (e ValidateApplicationsErrors) Error() string {
	errors := []string{}
	if e.Len() > 0 {
		for _, err := range e.Errors {
			errors = append(errors, fmt.Sprintf("%s was not found", err))
		}
		return fmt.Sprintf("failed to validate following applications: ['%s']", strings.Join(errors, "', '"))
	}
	return ""
}

type ValidateApplicationsOpts struct {
	Paths []string
}

func ValidateApplications(opts ValidateApplicationsOpts) error {
	errors := ValidateApplicationsErrors{}

	for _, application := range opts.Paths {
		_, err := exec.LookPath(application)
		if err != nil {
			errors.Push(application)
		}
	}

	if errors.Len() > 0 {
		return errors
	}

	return nil
}
