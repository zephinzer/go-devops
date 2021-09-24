package devops

import (
	"errors"
	"fmt"
)

// Environment holds a structure for retrieving environment values
// loaded from a call to LoadEnvironment
type Environment interface {
	// Get is the shorthand convenience method for .GetString
	Get(name string) string

	// GetE is the shorthand convenience method for .GetStringE
	GetE(name string) (string, error)

	// GetBool retrieves the environment variable with key `name` and
	// returns it as a boolean value. Panics if the value is not
	// parseable as a boolean
	GetBool(name string) bool

	// GetBoolE retrieves the environment variable with key `name` and
	// returns it as a boolean value. Returns an error in the return
	// values if the value is not parseable as a boolean
	GetBoolE(name string) (bool, error)

	// GetFloat retrieves the environment variable with key `name` and
	// returns it as a float value. Panics if the value is not
	// parseable as a float value
	GetFloat(name string) float64

	// GetFloatE retrieves the environment variable with key `name` and
	// returns it as a float value. Returns an error in the return
	// values if the value is not parseable as a float value
	GetFloatE(name string) (float64, error)

	// GetInt retrieves the environment variable with key `name` and
	// returns it as an integer value. Panics if the value is not
	// parseable as an integer
	GetInt(name string) int

	// GetIntE retrieves the environment variable with key `name` and
	// returns it as an integer. Returns an error in the return
	// values if the value is not parseable as an integer
	GetIntE(name string) (int, error)

	// GetString retrieves the environment variable with key `name` and
	// returns it as a string value. Panics if the value is not
	// parseable as a string
	GetString(name string) string

	// GetStringE retrieves the environment variable with key `name` and
	// returns it as a string. Returns an error in the return
	// values if the value is not parseable as a string
	GetStringE(name string) (string, error)

	// GetUint retrieves the environment variable with key `name` and
	// returns it as an unsigned integer value. Panics if the value is
	// not parseable as an unsigned integer
	GetUint(name string) uint

	// GetUintE retrieves the environment variable with key `name` and
	// returns it as an unsigned integer. Returns an error in the return
	// values if the value is not parseable as an unsigned integer
	GetUintE(name string) (uint, error)
}

type environment map[string]interface{}

func (e environment) Get(name string) string {
	return e.GetString(name)
}

func (e environment) GetE(name string) (string, error) {
	return e.GetStringE(name)
}

func (e environment) GetBool(name string) bool {
	if val, ok := e[name].(bool); ok {
		return val
	}
	panic(fmt.Sprintf("'%s' is not a boolean", name))
}

func (e environment) GetBoolE(name string) (val bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	return e.GetBool(name), nil
}

func (e environment) GetFloat(name string) float64 {
	if val, ok := e[name].(float64); ok {
		return val
	}
	panic(fmt.Sprintf("'%s' is not a float", name))
}

func (e environment) GetFloatE(name string) (val float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	return e.GetFloat(name), nil
}

func (e environment) GetInt(name string) int {
	if val, ok := e[name].(int); ok {
		return val
	}
	panic(fmt.Sprintf("'%s' is not an integer", name))
}

func (e environment) GetIntE(name string) (val int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	return e.GetInt(name), nil
}

func (e environment) GetString(name string) string {
	if val, ok := e[name].(string); ok {
		return val
	}
	panic(fmt.Sprintf("'%s' is not a string", name))
}

func (e environment) GetStringE(name string) (val string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	return e.GetString(name), nil
}

func (e environment) GetUint(name string) uint {
	if val, ok := e[name].(uint); ok {
		return val
	}
	panic(fmt.Sprintf("'%s' is not an unsigned integer", name))
}

func (e environment) GetUintE(name string) (val uint, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	return e.GetUint(name), nil
}
