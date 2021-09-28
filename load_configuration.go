package devops

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultStringSliceDelimiter        = ","
	ErrorLoadConfigurationPrereqs      = 1 << iota
	ErrorLoadConfigurationNotFound     = 1 << iota
	ErrorLoadConfigurationInvalidType  = 1 << iota
	ErrorLoadConfigurationInvalidValue = 1 << iota
)

type LoadConfigurationErrors []LoadConfigurationError

func (e LoadConfigurationErrors) GetCode() int {
	code := 0
	for _, err := range e {
		code |= err.Code
	}
	return code
}

func (e LoadConfigurationErrors) GetMessage() string {
	messages := []string{}
	for _, err := range e {
		messages = append(messages, err.Message)
	}
	return fmt.Sprintf("['%s']", strings.Join(messages, "', '"))
}

func (e LoadConfigurationErrors) Error() string {
	codes := 0
	messages := []string{}
	for _, err := range e {
		codes |= err.Code
		messages = append(messages, err.Message)
	}
	return fmt.Sprintf("LoadConfiguration/err[%v]: ['%s']", codes, strings.Join(messages, "', '"))
}

type LoadConfigurationError struct {
	Code    int
	Message string
}

func (e LoadConfigurationError) Error() string {
	return fmt.Sprintf("LoadConfiguration/err[%v]: %s", e.Code, e.Message)
}

func LoadConfiguration(config interface{}) error {
	errors := LoadConfigurationErrors{}

	c := newConfiguration(config)
	if !c.IsPointer() {
		errors = append(errors, LoadConfigurationError{ErrorLoadConfigurationPrereqs, "failed to receive a valid pointer"})
	}
	if !c.IsStruct() {
		errors = append(errors, LoadConfigurationError{ErrorLoadConfigurationPrereqs, "failed to receive a valid struct"})
	}

	if len(errors) > 0 {
		return errors
	}

	for _, field := range c.Fields {
		environmentKey := field.GetEnvironmentKey()
		environmentValue, isEnvironmentDefined := os.LookupEnv(environmentKey)
		defaultValue := field.GetDefaultValue()
		fieldType := field.Type.String()
		switch fieldType {
		case "[]string":
			var stringSliceValue []string
			delimiter, found := field.Tag.Lookup("delimiter")
			if !found {
				delimiter = DefaultStringSliceDelimiter
			}
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (string)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			stringValue = strings.Trim(stringValue, delimiter)
			stringSliceValue = strings.Split(stringValue, delimiter)
			field.SetStringSlice(stringSliceValue)
		case "*[]string":
			var stringSliceValue []string
			delimiter, found := field.Tag.Lookup("delimiter")
			if !found {
				delimiter = DefaultStringSliceDelimiter
			}
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			stringValue = strings.Trim(stringValue, delimiter)
			stringSliceValue = strings.Split(stringValue, delimiter)
			field.SetStringSlicePointer(stringSliceValue)
		case "string":
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (string)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			field.SetString(stringValue)
		case "*string":
			var stringValue string
			if defaultValue != nil {
				stringValue = *defaultValue
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				stringValue = environmentValue
			}
			field.SetStringPointer(stringValue)
		case "bool":
			var boolValue bool
			var err error
			if defaultValue != nil {
				if boolValue, err = strconv.ParseBool(*defaultValue); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (bool)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				if boolValue, err = strconv.ParseBool(environmentValue); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetBool(boolValue)
		case "*bool":
			var boolValue bool
			var err error
			if defaultValue != nil {
				if boolValue, err = strconv.ParseBool(*defaultValue); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				if boolValue, err = strconv.ParseBool(environmentValue); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetBoolPointer(boolValue)
		case "int":
			var intValue int64
			var err error
			if defaultValue != nil {
				if intValue, err = strconv.ParseInt(*defaultValue, 10, 0); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				errors = append(errors, LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (int)", field.Name, environmentKey),
				})
			}
			if isEnvironmentDefined {
				if intValue, err = strconv.ParseInt(environmentValue, 10, 0); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetInt(int(intValue))
		case "*int":
			var intValue int64
			var err error
			if defaultValue != nil {
				if intValue, err = strconv.ParseInt(*defaultValue, 10, 0); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", *defaultValue, field.Name),
					})
				}
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				if intValue, err = strconv.ParseInt(environmentValue, 10, 0); err != nil {
					errors = append(errors, LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", environmentValue, field.Name),
					})
				}
			}
			field.SetIntPointer(int(intValue))
		default:
			errors = append(errors, LoadConfigurationError{
				ErrorLoadConfigurationInvalidType,
				fmt.Sprintf("failed to load '%s' (via \"${%s}\") of type '%s'", field.Name, environmentKey, field.Type.String()),
			})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
