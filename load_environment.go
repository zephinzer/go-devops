package devops

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EnvironmentOpts struct {
	Description string
	Key         string
	Type        EnvType
	Default     interface{}
}

type LoadEnvironmentOpts []EnvironmentOpts

func (o LoadEnvironmentOpts) SetDefaults() {
	for i := 0; i < len(o); i++ {
		switch o[i].Type {
		case TypeBool:
			fallthrough
		case TypeFloat:
			fallthrough
		case TypeInt:
			fallthrough
		case TypeUint:
			break
		default:
			o[i].Type = TypeString
		}
	}
}

func (o LoadEnvironmentOpts) Validate() error {
	errors := []string{}

	keys := map[string]bool{}

	for _, env := range o {
		if env.Key == "" {
			errors = append(errors, fmt.Sprintf("missing key for '%s'", env.Description))
		} else if seen, ok := keys[env.Key]; ok && seen {
			errors = append(errors, fmt.Sprintf("'%s' is duplicated", env.Key))
		}
		keys[env.Key] = true
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}

	return nil
}

func LoadEnvironment(opts LoadEnvironmentOpts) (Environment, error) {
	opts.SetDefaults()
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("failed to load environment: %s", err)
	}
	envInstance := environment{}
	for _, envRequest := range opts {
		rawValue := os.Getenv(envRequest.Key)
		switch envRequest.Type {
		case TypeBool:
			if val, err := strconv.ParseBool(rawValue); err != nil {
				envInstance[envRequest.Key] = nil
				if val, ok := envRequest.Default.(bool); ok {
					envInstance[envRequest.Key] = val
				}
			} else {
				envInstance[envRequest.Key] = val
			}
		case TypeFloat:
			if val, err := strconv.ParseFloat(rawValue, 32); err != nil {
				envInstance[envRequest.Key] = nil
				if val, ok := envRequest.Default.(float64); ok {
					envInstance[envRequest.Key] = val
				}
			} else {
				envInstance[envRequest.Key] = val
			}
		case TypeInt:
			if val, err := strconv.ParseInt(rawValue, 10, 0); err != nil {
				envInstance[envRequest.Key] = nil
				if val, ok := envRequest.Default.(int); ok {
					envInstance[envRequest.Key] = val
				}
			} else {
				envInstance[envRequest.Key] = int(val)
			}
		case TypeUint:
			if val, err := strconv.ParseUint(rawValue, 10, 0); err != nil {
				if val, ok := envRequest.Default.(int); ok {
					envInstance[envRequest.Key] = uint(val)
				} else if val, ok := envRequest.Default.(uint); ok {
					envInstance[envRequest.Key] = val
				}
			} else {
				envInstance[envRequest.Key] = uint(val)
			}
		default:
			envInstance[envRequest.Key] = rawValue
			if len(rawValue) == 0 {
				if val, ok := envRequest.Default.(string); ok {
					envInstance[envRequest.Key] = val
				}
			}
		}
	}
	return envInstance, nil
}
