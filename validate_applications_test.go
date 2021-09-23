package devops

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidateApplicationsTest struct {
	suite.Suite
}

func TestValidateApplications(t *testing.T) {
	suite.Run(t, &ValidateApplicationsTest{})
}

func (s ValidateApplicationsTest) TestValidateApplications() {
	expectedSuccessful := []string{"printf", "echo", "./tests/paths/to/a/script.sh"}
	expectedFailure := []string{"doesnotexist", "alsodoesntexist", "./tests/paths/to/a/file"}

	err := ValidateApplications(ValidateApplicationsOpts{
		Paths: expectedSuccessful,
	})
	s.Nil(err)

	err = ValidateApplications(ValidateApplicationsOpts{
		Paths: expectedFailure,
	})
	s.NotNil(err)
	errs := err.(ValidateApplicationsErrors)

	for _, expectedFail := range expectedFailure {
		s.Contains(errs.Errors, expectedFail)
	}
}
