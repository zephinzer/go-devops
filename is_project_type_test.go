package devops

import (
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IsProjectTypeTest struct {
	suite.Suite
	TestDirectory string
}

func TestIsProjectType(t *testing.T) {
	suite.Run(t, &IsProjectTypeTest{
		TestDirectory: "./tests/projecttypes",
	})
}

func (s IsProjectTypeTest) Test_invalid() {
	testPath := path.Join(s.TestDirectory, "invalidtype")
	yes, err := IsProjectType(testPath, TypeC)
	s.NotNil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_C() {
	testPath := path.Join(s.TestDirectory, "c")
	yes, err := IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeGo)
	s.Nil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_Go() {
	testPath := path.Join(s.TestDirectory, "go")
	yes, err := IsProjectType(testPath, TypeGo)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_Haskell() {
	testPath := path.Join(s.TestDirectory, "hs")
	yes, err := IsProjectType(testPath, TypeHaskell)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_Java() {
	testPath := path.Join(s.TestDirectory, "java")
	yes, err := IsProjectType(testPath, TypeJava)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_Javascript() {
	testPath := path.Join(s.TestDirectory, "js")
	yes, err := IsProjectType(testPath, TypeJavascript)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
	testPath = path.Join(s.TestDirectory, "js")
	yes, err = IsProjectType(testPath, TypeJavascript)
	s.Nil(err)
	s.True(yes)
}

func (s IsProjectTypeTest) Test_Python() {
	testPath := path.Join(s.TestDirectory, "py")
	yes, err := IsProjectType(testPath, TypePython)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_Ruby() {
	testPath := path.Join(s.TestDirectory, "rb")
	yes, err := IsProjectType(testPath, TypeRuby)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_Rust() {
	testPath := path.Join(s.TestDirectory, "rs")
	yes, err := IsProjectType(testPath, TypeRust)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
}

func (s IsProjectTypeTest) Test_Typescript() {
	testPath := path.Join(s.TestDirectory, "ts")
	yes, err := IsProjectType(testPath, TypeTypescript)
	s.Nil(err)
	s.True(yes)
	yes, err = IsProjectType(testPath, TypeC)
	s.Nil(err)
	s.False(yes)
}
