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

type LoadConfigurationError struct {
	Code    int
	Message string
}

func (e LoadConfigurationError) Error() string {
	return fmt.Sprintf("LoadConfiguration/err[%v]: %s", e.Code, e.Message)
}

func LoadConfiguration(config interface{}) error {
	c := newConfiguration(config)
	if !c.IsPointer() {
		return LoadConfigurationError{ErrorLoadConfigurationPrereqs, "failed to receive a valid pointer"}
	}
	if !c.IsStruct() {
		return LoadConfigurationError{ErrorLoadConfigurationPrereqs, "failed to receive a valid struct"}
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
				return LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (string)", field.Name, environmentKey),
				}
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
				return LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (string)", field.Name, environmentKey),
				}
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
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", *defaultValue, field.Name),
					}
				}
			} else if !isEnvironmentDefined {
				return LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (bool)", field.Name, environmentKey),
				}
			}
			if isEnvironmentDefined {
				if boolValue, err = strconv.ParseBool(environmentValue); err != nil {
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", environmentValue, field.Name),
					}
				}
			}
			field.SetBool(boolValue)
		case "*bool":
			var boolValue bool
			var err error
			if defaultValue != nil {
				if boolValue, err = strconv.ParseBool(*defaultValue); err != nil {
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", *defaultValue, field.Name),
					}
				}
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				if boolValue, err = strconv.ParseBool(environmentValue); err != nil {
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as a boolean for loading '%s'", environmentValue, field.Name),
					}
				}
			}
			field.SetBoolPointer(boolValue)
		case "int":
			var intValue int64
			var err error
			if defaultValue != nil {
				if intValue, err = strconv.ParseInt(*defaultValue, 10, 0); err != nil {
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", *defaultValue, field.Name),
					}
				}
			} else if !isEnvironmentDefined {
				return LoadConfigurationError{
					ErrorLoadConfigurationNotFound,
					fmt.Sprintf("failed to load '%s' via \"${%s}\" (int)", field.Name, environmentKey),
				}
			}
			if isEnvironmentDefined {
				if intValue, err = strconv.ParseInt(environmentValue, 10, 0); err != nil {
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", environmentValue, field.Name),
					}
				}
			}
			field.SetInt(int(intValue))
		case "*int":
			var intValue int64
			var err error
			if defaultValue != nil {
				if intValue, err = strconv.ParseInt(*defaultValue, 10, 0); err != nil {
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", *defaultValue, field.Name),
					}
				}
			} else if !isEnvironmentDefined {
				break
			}
			if isEnvironmentDefined {
				if intValue, err = strconv.ParseInt(environmentValue, 10, 0); err != nil {
					return LoadConfigurationError{
						ErrorLoadConfigurationInvalidValue,
						fmt.Sprintf("failed to parse '%s' as an int for loading '%s'", environmentValue, field.Name),
					}
				}
			}
			field.SetIntPointer(int(intValue))
		default:
			return LoadConfigurationError{
				ErrorLoadConfigurationInvalidType,
				fmt.Sprintf("failed to load '%s' (via \"${%s}\") of type '%s'", field.Name, environmentKey, field.Type.String()),
			}
		}
	}
	return nil
}
