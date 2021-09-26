package devops

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigurationTest struct {
	suite.Suite
}

func TestConfiguration(t *testing.T) {
	suite.Run(t, &ConfigurationTest{})
}

func (s ConfigurationTest) Test_newConfiguration() {
	type testStruct struct {
		Bool   bool
		Float  float64
		Int    int
		String string
		Uint   uint
	}
	var config configuration
	s.NotPanics(func() {
		config = newConfiguration(&testStruct{})
	})
	s.Len(config.Fields, 5)
}

func (s ConfigurationTest) Test_configuration_IsPointer() {
	type testStruct struct{}
	config := newConfiguration(testStruct{})
	s.False(config.IsPointer())
	config = newConfiguration(&testStruct{})
	s.True(config.IsPointer())
}

func (s ConfigurationTest) Test_configuration_IsStruct() {
	type testStruct struct{}
	config := newConfiguration(testStruct{})
	s.True(config.IsStruct())
	config = newConfiguration(&testStruct{})
	s.True(config.IsStruct())
}

func (s ConfigurationTest) Test_configurationField_GetDefaultValue() {
	type testStruct struct {
		DefaultField     string `default:"default"`
		DefaultlessField string
	}
	config := newConfiguration(&testStruct{})
	defaultValue := config.Fields[0].GetDefaultValue()
	s.NotNil(defaultValue)
	s.Equal("default", *defaultValue)
	defaultValue = config.Fields[1].GetDefaultValue()
	s.Nil(defaultValue)
}

func (s ConfigurationTest) Test_configurationField_GetEnvironmentKey() {
	type testStruct struct {
		WithEnv    string `env:"env"`
		WithoutEnv string
	}
	config := newConfiguration(&testStruct{})
	environmentKey := config.Fields[0].GetEnvironmentKey()
	s.Equal("env", environmentKey)
	environmentKey = config.Fields[1].GetEnvironmentKey()
	s.Equal("WITHOUT_ENV", environmentKey)
}

func (s ConfigurationTest) Test_configurationField_Setters() {
	type testStruct struct {
		OptionalBool        *bool
		RequiredBool        bool
		OptionalInt         *int
		RequiredInt         int
		OptionalString      *string
		RequiredString      string
		OptionalStringSlice *[]string
		RequiredStringSlice []string
	}
	testStructInstance := testStruct{}
	config := newConfiguration(&testStructInstance)
	config.Fields[0].SetBoolPointer(true)
	s.Equal(true, *testStructInstance.OptionalBool)
	config.Fields[1].SetBool(true)
	s.Equal(true, testStructInstance.RequiredBool)
	config.Fields[2].SetIntPointer(-1)
	s.Equal(-1, *testStructInstance.OptionalInt)
	config.Fields[3].SetInt(-2)
	s.Equal(-2, testStructInstance.RequiredInt)
	config.Fields[4].SetStringPointer("hello")
	s.Equal("hello", *testStructInstance.OptionalString)
	config.Fields[5].SetString("world")
	s.Equal("world", testStructInstance.RequiredString)
	config.Fields[6].SetStringSlicePointer([]string{"hola", "mundo"})
	s.EqualValues([]string{"hola", "mundo"}, *testStructInstance.OptionalStringSlice)
	config.Fields[7].SetStringSlice([]string{"hola", "para", "ti"})
	s.EqualValues([]string{"hola", "para", "ti"}, testStructInstance.RequiredStringSlice)
}
