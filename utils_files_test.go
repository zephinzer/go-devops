package devops

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsFilesTest struct {
	suite.Suite
}

func TestUtilsFiles(t *testing.T) {
	suite.Run(t, &UtilsFilesTest{})
}

func (s UtilsFilesTest) Test_recursivelyGetExtensionsCount() {
	extensions, err := recursivelyGetExtensionsCount("./tests/projecttypes")
	s.Nil(err)
	s.Equal(1, extensions[".c"])
	s.Equal(1, extensions[".java"])
	s.Equal(1, extensions[".js"])
	s.Equal(1, extensions[".py"])
	s.Equal(1, extensions[".rb"])
	s.Equal(1, extensions[".rs"])
	s.Equal(1, extensions[".ts"])
}
