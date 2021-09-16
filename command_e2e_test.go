package devops

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommandEndToEndTestsSuite struct {
	suite.Suite
}

func TestCommandEndToEnd(t *testing.T) {
	suite.Run(t, &CommandEndToEndTestsSuite{})
}

// Test_Basic checks for basic expectations of a command concept
// like invocation with arguments
func (s CommandEndToEndTestsSuite) Test_Basic() {
	scriptPath := "tests/command/basic.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	command, err := NewCommand(NewCommandOpts{
		Command:   scriptPath,
		Arguments: []string{"arg1", "arg2", "arg3"},
	})
	s.Nil(err, "this command should be created successfully")
	err = command.Run()
	s.Nil(err, "the command should have ran successfully")
	output := string(command.GetOutput())
	s.NotEmpty(output, "the output should not be empty")
	s.Contains(output, "$@: arg1 arg2 arg3")
	s.Contains(output, fmt.Sprintf("$0: %s", scriptPath))
	s.Contains(output, "$1: arg1")
	s.Contains(output, "$2: arg2")
	s.Contains(output, "$3: arg3")
}
