package devops

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfirmTests struct {
	suite.Suite
}

func TestConfirm(t *testing.T) {
	suite.Run(t, &ConfirmTests{})
}

func (s ConfirmTests) TestConfirm_MatchExact() {
	var input bytes.Buffer
	input.Write([]byte("yes"))
	options := ConfirmOpts{Input: &input, MatchExact: "yes"}
	result, err := Confirm(options)
	s.True(result)
	s.Nil(err)
}

func (s ConfirmTests) TestConfirm_MatchRegexp() {
	var input bytes.Buffer
	input.Write([]byte("yes"))
	options := ConfirmOpts{Input: &input, MatchRegexp: regexp.MustCompile("^yes$")}
	result, err := Confirm(options)
	s.True(result)
	s.Nil(err)
}

func (s ConfirmTests) TestConfirm_Output() {
	var input bytes.Buffer
	var output bytes.Buffer
	input.Write([]byte("yes"))
	options := ConfirmOpts{Input: &input, Output: &output, MatchExact: "yes"}
	result, err := Confirm(options)
	s.True(result)
	s.Nil(err)
	s.Contains(output.String(), fmt.Sprintf(DefaultConfirmInputHint, "yes"))
}

func (s ConfirmTests) TestConfirm_Validation() {
	options := ConfirmOpts{
		MatchExact:  "yes",
		MatchRegexp: regexp.MustCompile("^yes$"),
	}
	result, err := Confirm(options)
	s.False(result)
	s.NotNil(err, "should fail if both matchers are defined")

	options = ConfirmOpts{}
	result, err = Confirm(options)
	s.False(result)
	s.NotNil(err, "should fail if both matchers are not defined")
}
