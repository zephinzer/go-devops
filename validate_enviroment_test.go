package devops

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidateEnvironmentTests struct {
	suite.Suite
	StringKey string
	IntKey    string
	UintKey   string
	FloatKey  string
	BoolKey   string
}

func TestValidateEnvironment(t *testing.T) {
	testSuite := &ValidateEnvironmentTests{
		StringKey: "STRING",
		IntKey:    "INT",
		UintKey:   "UINT",
		FloatKey:  "FLOAT",
		BoolKey:   "BOOL",
	}
	suite.Run(t, testSuite)
}

func (s ValidateEnvironmentTests) AfterTest(suiteName, testName string) {
	os.Unsetenv(s.StringKey)
	os.Unsetenv(s.IntKey)
	os.Unsetenv(s.UintKey)
	os.Unsetenv(s.FloatKey)
	os.Unsetenv(s.BoolKey)
}

func (s ValidateEnvironmentTests) TestValidateEnvironment_basic() {
	os.Setenv("HEY", "hello world")
	err := ValidateEnvironment(ValidateEnvironmentOpts{
		Keys: EnvironmentKeys{{Name: "HEY"}},
	})
	s.Nil(err)
}

func (s ValidateEnvironmentTests) TestValidateEnvironment_typesAny() {
	os.Setenv(s.StringKey, "LALALA")
	os.Setenv(s.IntKey, "-25092019")
	os.Setenv(s.UintKey, "25092019")
	os.Setenv(s.FloatKey, "3.142")
	os.Setenv(s.BoolKey, "true")
	err := ValidateEnvironment(ValidateEnvironmentOpts{
		Keys: EnvironmentKeys{
			{Name: s.StringKey, Type: TypeAny},
			{Name: s.IntKey, Type: TypeAny},
			{Name: s.UintKey, Type: TypeAny},
			{Name: s.FloatKey, Type: TypeAny},
			{Name: s.BoolKey, Type: TypeAny},
		},
	})
	s.Nil(err)
}

func (s ValidateEnvironmentTests) TestValidateEnvironment_typesUnknown() {
	os.Setenv(s.StringKey, "HELLO WORLD")
	err := ValidateEnvironment(ValidateEnvironmentOpts{
		Keys: EnvironmentKeys{
			{Name: s.StringKey, Type: EnvType("custom")},
		},
	})
	s.NotNil(err)
	s.Contains(err.Error(), "key[STRING]")
}

func (s ValidateEnvironmentTests) TestValidateEnvironment_typesHappy() {
	os.Setenv(s.StringKey, "LALALA")
	os.Setenv(s.IntKey, "-25092019")
	os.Setenv(s.UintKey, "25092019")
	os.Setenv(s.FloatKey, "3.142")
	os.Setenv(s.BoolKey, "true")
	err := ValidateEnvironment(ValidateEnvironmentOpts{
		Keys: EnvironmentKeys{
			{Name: s.StringKey, Type: TypeString},
			{Name: s.IntKey, Type: TypeInt},
			{Name: s.UintKey, Type: TypeUint},
			{Name: s.FloatKey, Type: TypeFloat},
			{Name: s.BoolKey, Type: TypeBool},
		},
	})
	s.Nil(err)
}

func (s ValidateEnvironmentTests) TestValidateEnvironment_typesSad() {
	os.Setenv(s.IntKey, "not an int")
	os.Setenv(s.UintKey, "-123")
	os.Setenv(s.FloatKey, "not a float")
	os.Setenv(s.BoolKey, "maybe")
	err := ValidateEnvironment(ValidateEnvironmentOpts{
		Keys: EnvironmentKeys{
			{Name: s.StringKey, Type: TypeString},
			{Name: s.IntKey, Type: TypeInt},
			{Name: s.UintKey, Type: TypeUint},
			{Name: s.FloatKey, Type: TypeFloat},
			{Name: s.BoolKey, Type: TypeBool},
		},
	})
	s.NotNil(err)
	s.Contains(err.Error(), "key[STRING]")
	s.Contains(err.Error(), "key[INT]")
	s.Contains(err.Error(), "key[UINT]")
	s.Contains(err.Error(), "key[FLOAT]")
	s.Contains(err.Error(), "key[BOOL]")
}
