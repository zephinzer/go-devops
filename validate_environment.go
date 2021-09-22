package devops

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EnvType string

const (
	TypeAny    EnvType = "any"
	TypeString EnvType = "string"
	TypeInt    EnvType = "int"
	TypeUint   EnvType = "uint"
	TypeFloat  EnvType = "float"
	TypeBool   EnvType = "bool"
)

type EnvironmentKey struct {
	Name string
	Type EnvType
}

type EnvironmentKeys []EnvironmentKey

type ValidateEnvironmentOpts struct {
	Keys EnvironmentKeys
}

func (o *ValidateEnvironmentOpts) SetDefaults() {
	for i := 0; i < len(o.Keys); i++ {
		if o.Keys[i].Type == "" {
			o.Keys[i].Type = TypeAny
		}
	}
}

func ValidateEnvironment(opts ValidateEnvironmentOpts) error {
	errors := []string{}

	opts.SetDefaults()
	expectedKeyMap := map[string]EnvType{}
	for _, expectedKey := range opts.Keys {
		expectedKeyMap[expectedKey.Name] = expectedKey.Type
	}

	keyMap := map[string]interface{}{}
	environment := os.Environ()
	for _, environmentEntry := range environment {
		keyValuePair := strings.SplitN(environmentEntry, "=", 2)
		if len(keyValuePair) == 2 {
			keyMap[keyValuePair[0]] = keyValuePair[1]
		} else if len(keyValuePair) == 1 {
			keyMap[keyValuePair[0]] = nil
		}
	}

	for expectedKey, expectedType := range expectedKeyMap {
		if value, ok := keyMap[expectedKey]; ok { // key exists at least
			switch expectedType {
			case TypeAny:
			case TypeString:
				fmt.Println("------------------")
				if value == nil {
					errors = append(errors, fmt.Sprintf("key[%s]:string was nil", expectedKey))
				} else if val := value.(string); val == "" {
					errors = append(errors, fmt.Sprintf("key[%s]:string was empty", expectedKey))
				}
				fmt.Println(value)
			case TypeInt:
				if _, err := strconv.Atoi(value.(string)); err != nil {
					errors = append(errors, fmt.Sprintf("key[%s]:integer was '%s'", expectedKey, value.(string)))
				}
			case TypeUint:
				if _, err := strconv.ParseUint(value.(string), 10, 0); err != nil {
					errors = append(errors, fmt.Sprintf("key[%s]:uint was '%s'", expectedKey, value.(string)))
				}
			case TypeFloat:
				if _, err := strconv.ParseFloat(value.(string), 0); err != nil {
					errors = append(errors, fmt.Sprintf("key[%s]:float was '%s'", expectedKey, value.(string)))
				}
			case TypeBool:
				if _, err := strconv.ParseBool(value.(string)); err != nil {
					errors = append(errors, fmt.Sprintf("key[%s]:bool was '%s'", expectedKey, value.(string)))
				}
			default:
				errors = append(errors, fmt.Sprintf("key[%s] has unknown type '%s'", expectedKey, expectedType))
			}
		} else {
			errors = append(errors, fmt.Sprintf("key[%s] does not exist", expectedKey))
			continue
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate environment: ['%s']", strings.Join(errors, "', '"))
	}

	return nil
}
