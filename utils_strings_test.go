package devops

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsStringsTest struct {
	suite.Suite
}

func TestUtilsStrings(t *testing.T) {
	suite.Run(t, &UtilsStringsTest{})
}

func (s UtilsStringsTest) Test_containsAllStrings() {
	s.False(containsAllStrings(nil, []string{"hola"}))
	s.False(containsAllStrings([]string{"hola", "mundo"}, nil))
	s.False(containsAllStrings([]string{"hola", "mundo"}, []string{"hola", "world"}))
	s.True(containsAllStrings([]string{"hola", "mundo"}, []string{"hola", "mundo"}))
}

func (s UtilsStringsTest) Test_containsAnyString() {
	s.False(containsAnyString(nil, []string{"hola"}))
	s.False(containsAnyString([]string{"hola", "mundo"}, nil))
	s.True(containsAnyString([]string{"hola", "mundo"}, []string{"hola", "world"}))
	s.True(containsAnyString([]string{"hola", "mundo"}, []string{"hello", "mundo"}))
	s.True(containsAnyString([]string{"hola", "mundo"}, []string{"hola", "mundo"}))
}
