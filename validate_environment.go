package devops

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EnvType string

const (
	TypeAny          EnvType = "any"
	TypeNil          EnvType = "nil"
	TypeString       EnvType = "string"
	TypeInt          EnvType = "int"
	TypeUint         EnvType = "uint"
	TypeFloat        EnvType = "float"
	TypeBool         EnvType = "bool"
	TypeErrorMissing EnvType = "__missing"
	TypeErrorUnknown EnvType = "__unknown"
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

type ValidateEnvironmentError struct {
	Key          string
	ExpectedType EnvType
	Value        string
}

type ValidateEnvironmentErrors struct {
	Errors []ValidateEnvironmentError
}

func (e ValidateEnvironmentErrors) Len() int {
	return len(e.Errors)
}

func (e *ValidateEnvironmentErrors) Push(err ValidateEnvironmentError) {
	e.Errors = append(e.Errors, err)
}

func (e ValidateEnvironmentErrors) Error() string {
	errors := []string{}
	for _, err := range e.Errors {
		if err.ExpectedType == TypeErrorUnknown {
			errors = append(errors, fmt.Sprintf("key[%s] has unknown type '%s'", err.Key, err.Value))
		} else if err.ExpectedType == TypeErrorMissing {
			errors = append(errors, fmt.Sprintf("key[%s] does not exist", err.Key))
		} else {
			errors = append(errors, fmt.Sprintf("key[%s]:%s was '%s'", err.Key, err.ExpectedType, err.Value))
		}
	}
	return fmt.Sprintf("failed to validate environment: ['%s']", strings.Join(errors, "', '"))
}

func ValidateEnvironment(opts ValidateEnvironmentOpts) error {
	errors := ValidateEnvironmentErrors{}

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
					errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: expectedType, Value: "nil"})
				} else if val := value.(string); val == "" {
					errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: expectedType, Value: "empty"})
				}
				fmt.Println(value)
			case TypeInt:
				if _, err := strconv.Atoi(value.(string)); err != nil {
					errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: expectedType, Value: value.(string)})
				}
			case TypeUint:
				if _, err := strconv.ParseUint(value.(string), 10, 0); err != nil {
					errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: expectedType, Value: value.(string)})
				}
			case TypeFloat:
				if _, err := strconv.ParseFloat(value.(string), 0); err != nil {
					errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: expectedType, Value: value.(string)})
				}
			case TypeBool:
				if _, err := strconv.ParseBool(value.(string)); err != nil {
					errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: expectedType, Value: value.(string)})
				}
			default:
				errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: TypeErrorUnknown, Value: string(expectedType)})
			}
		} else {
			errors.Push(ValidateEnvironmentError{Key: expectedKey, ExpectedType: TypeErrorMissing})
			continue
		}
	}

	if errors.Len() > 0 {
		return errors
	}

	return nil
}
