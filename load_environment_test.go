package devops

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoadEnvironmentTest struct {
	suite.Suite
	opts LoadEnvironmentOpts
}

func TestLoadEnvironment(t *testing.T) {
	suite.Run(t, &LoadEnvironmentTest{
		opts: LoadEnvironmentOpts{
			{Description: "bool", Key: "bool", Type: TypeBool},
			{Description: "float", Key: "float", Type: TypeFloat},
			{Description: "int", Key: "int", Type: TypeInt},
			{Description: "uint", Key: "uint", Type: TypeUint},
			{Description: "string", Key: "string", Type: TypeString},
			{Description: "nil", Key: "nil", Type: TypeNil},
			{Description: "any", Key: "any", Type: TypeAny},
			{Description: "others", Key: "others", Type: EnvType("others")},
		},
	})
}

func (s LoadEnvironmentTest) TestLoadEnvironmentOpts() {
	expectedType := map[string]EnvType{
		"bool":   TypeBool,
		"float":  TypeFloat,
		"int":    TypeInt,
		"uint":   TypeUint,
		"string": TypeString,
		"nil":    TypeString,
		"any":    TypeString,
		"others": TypeString,
	}
	loadEnvironmentOpts := s.opts
	loadEnvironmentOpts.SetDefaults()
	for _, loadedEnv := range loadEnvironmentOpts {
		s.Equal(expectedType[loadedEnv.Key], loadedEnv.Type)
	}
}

func (s LoadEnvironmentTest) TestLoadEnvironment_bool() {
	os.Setenv("string", "false")
	os.Setenv("bool", "t")
	defer os.Unsetenv("bool")
	loadEnvironmentOpts := s.opts
	env, err := LoadEnvironment(loadEnvironmentOpts)
	s.Nil(err)
	s.NotPanics(func() { s.Equal(true, env.GetBool("bool")) })
	_, err = env.GetBoolE("bool")
	s.Nil(err)
	s.Panics(func() { env.GetBool("string") })
	_, err = env.GetBoolE("string")
	s.NotNil(err)
}

func (s LoadEnvironmentTest) TestLoadEnvironment_float() {
	os.Setenv("int", "3.142")
	os.Setenv("float", "3.142")
	defer os.Unsetenv("float")
	loadEnvironmentOpts := s.opts
	env, err := LoadEnvironment(loadEnvironmentOpts)
	s.Nil(err)
	s.NotPanics(func() { s.Equal(3.142, math.Ceil(env.GetFloat("float")*1000)/float64(1000)) })
	_, err = env.GetFloatE("float")
	s.Nil(err)
	s.Panics(func() { env.GetFloat("int") })
	_, err = env.GetFloatE("int")
	s.NotNil(err)
}

func (s LoadEnvironmentTest) TestLoadEnvironment_int() {
	os.Setenv("uint", "-123456")
	os.Setenv("int", "-123456")
	defer os.Unsetenv("uint")
	defer os.Unsetenv("int")
	loadEnvironmentOpts := s.opts
	env, err := LoadEnvironment(loadEnvironmentOpts)
	s.Nil(err)
	s.NotPanics(func() { s.Equal(-123456, env.GetInt("int")) })
	_, err = env.GetIntE("int")
	s.Nil(err)
	s.Panics(func() { env.GetInt("uint") })
	_, err = env.GetIntE("uint")
	s.NotNil(err)
}

func (s LoadEnvironmentTest) TestLoadEnvironment_uint() {
	os.Setenv("bool", "123456")
	os.Setenv("uint", "123456")
	defer os.Unsetenv("bool")
	defer os.Unsetenv("uint")
	loadEnvironmentOpts := s.opts
	env, err := LoadEnvironment(loadEnvironmentOpts)
	s.Nil(err)
	s.NotPanics(func() { s.EqualValues(123456, env.GetUint("uint")) })
	_, err = env.GetUintE("uint")
	s.Nil(err)
	s.Panics(func() { env.GetUint("bool") })
	_, err = env.GetUintE("bool")
	s.NotNil(err)
}

func (s LoadEnvironmentTest) TestLoadEnvironment_string() {
	os.Setenv("int", "hello world")
	os.Setenv("string", "hello world")
	defer os.Unsetenv("string")
	defer os.Unsetenv("int")
	loadEnvironmentOpts := s.opts
	env, err := LoadEnvironment(loadEnvironmentOpts)
	s.Nil(err)
	s.NotPanics(func() { s.Equal("hello world", env.GetString("string")) })
	_, err = env.GetStringE("string")
	s.Nil(err)
	s.Panics(func() { env.GetString("int") })
	_, err = env.GetStringE("int")
	s.NotNil(err)
}

func (s LoadEnvironmentTest) TestLoadEnvironment_others() {
	os.Setenv("int", "hello world")
	os.Setenv("others", "hello world")
	defer os.Unsetenv("string")
	defer os.Unsetenv("int")
	loadEnvironmentOpts := s.opts
	env, err := LoadEnvironment(loadEnvironmentOpts)
	s.Nil(err)
	s.NotPanics(func() { s.Equal("hello world", env.Get("others")) })
	_, err = env.GetE("others")
	s.Nil(err)
	s.Panics(func() { env.Get("int") })
	_, err = env.GetE("int")
	s.NotNil(err)
}

func (s LoadEnvironmentTest) TestLoadEnvironment_validate() {
	opts := LoadEnvironmentOpts{
		{Description: "no key"},
	}
	_, err := LoadEnvironment(opts)
	s.NotNil(err)
	s.Contains(err.Error(), "no key")

	opts = append(opts, EnvironmentOpts{Key: "dupe"})
	opts = append(opts, EnvironmentOpts{Key: "dupe"})
	_, err = LoadEnvironment(opts)
	s.NotNil(err)
	s.Contains(err.Error(), "duplicated")
}
