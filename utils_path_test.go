package devops

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsPathTest struct {
	suite.Suite
}

func TestUtilsPath(t *testing.T) {
	suite.Run(t, &UtilsPathTest{})
}

func (s UtilsPathTest) TestNormalizeLocalPath_absoluteReference() {
	pathOfInterest := "/path/that/is/absolute"
	p, e := NormalizeLocalPath(pathOfInterest)
	s.Nil(e)
	s.Equal(pathOfInterest, p)
}

func (s UtilsPathTest) TestNormalizeLocalPath_homeReference() {
	pathOfInterest := "path/somewhere"
	p, e := NormalizeLocalPath(fmt.Sprintf("~/%s", pathOfInterest))
	s.Nil(e)
	uhd, e := os.UserHomeDir()
	s.Nil(e)
	s.Contains(p, uhd)
	s.Contains(p, pathOfInterest)
}

func (s UtilsPathTest) TestNormalizeLocalPath_relativeReference() {
	pathOfInterest := "path/somewhere"
	p, e := NormalizeLocalPath(fmt.Sprintf("%s", pathOfInterest))
	s.Nil(e)
	wd, e := os.Getwd()
	s.Nil(e)
	s.Contains(p, wd)
	s.Contains(p, pathOfInterest)
	fmt.Println(p)
}
