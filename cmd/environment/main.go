package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/zephinzer/go-strcase"
)

func parseConfiguration(config interface{}) error {
	configValue := reflect.ValueOf(config)
	isPointer := configValue.Kind() == reflect.Ptr
	if !isPointer {
		return fmt.Errorf("failed to receive a valid pointer")
	}
	configType := reflect.TypeOf(config).Elem()
	isStruct := configType.Kind() == reflect.Struct
	if !isStruct {
		return fmt.Errorf("failed to receive a valid struct")
	}

	fieldsCount := configType.NumField()
	fieldTypeMap := map[string]reflect.Type{}
	for i := 0; i < fieldsCount; i++ {
		fieldTypeMap[configType.Field(i).Name] = configType.Field(i).Type
	}
	values := make([]string, fieldsCount)
	for fieldName, fieldType := range fieldTypeMap {
		field := configValue.Elem().FieldByName(fieldName)
		var defaultValue *string
		typ, found := configType.FieldByName(fieldName)
		if v, ok := typ.Tag.Lookup("default"); found && ok {
			defaultValue = &v
		}
		environmentName := strcase.ToUpperSnake(fieldName)
		environmentValue, isValueDefined := os.LookupEnv(environmentName)
		switch fieldType.String() {
		case "string":
			if defaultValue != nil {
				field.SetString(*defaultValue)
			}
			if len(environmentValue) > 0 {
				field.SetString(environmentValue)
			}
		case "*string":
			if defaultValue != nil {
				reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(defaultValue))
			}
			if isValueDefined && len(environmentValue) > 0 {
				reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(&environmentValue))
			}
		case "bool":
			var boolValue bool
			var err error
			if defaultValue != nil {
				boolValue, err = strconv.ParseBool(*defaultValue)
			}
			if len(environmentValue) > 0 {
				boolValue, err = strconv.ParseBool(environmentValue)
				if err != nil {
					boolValue = false
				}
			}
			field.SetBool(boolValue)
		case "*bool":
			var boolValue bool
			var err error
			if defaultValue != nil {
				boolValue, err = strconv.ParseBool(*defaultValue)
			}
			if isValueDefined && len(environmentValue) > 0 {
				boolValue, err = strconv.ParseBool(environmentValue)
				if err != nil {
					boolValue = false
				}
			}
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(&boolValue))
		case "int":
			var intValue int64
			var err error
			if defaultValue != nil {
				intValue, err = strconv.ParseInt(*defaultValue, 10, 0)
				if err != nil {
					intValue = 0
				}
			}
			if len(environmentValue) > 0 {
				intValue, err = strconv.ParseInt(environmentValue, 10, 0)
				if err != nil {
					intValue = 0
				}
			}
			field.SetInt(intValue)
		case "*int":
			var intValue int64
			var err error
			if defaultValue != nil {
				intValue, err = strconv.ParseInt(*defaultValue, 10, 0)
				if err != nil {
					intValue = 0
				}
			}
			if isValueDefined && len(environmentValue) > 0 {
				intValue, err = strconv.ParseInt(environmentValue, 10, 0)
				if err != nil {
					intValue = 0
				}
			}
			noramlisedIntValue := int(intValue)
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(&noramlisedIntValue))
		}
		values = append(values, environmentValue)
	}
	return nil
}

type configuration struct {
	RequiredString string  `default:"hello world"`
	RequiredInt    int     `default:"1"`
	RequiredBool   bool    `default:"true"`
	OptionalString *string `default:"hola mundo"`
	OptionalInt    *int    `default:"2"`
	OptionalBool   *bool   `default:"true"`
}

func main() {
	c := configuration{}
	if err := parseConfiguration(&c); err != nil {
		log.Fatal(err)
	}
	log.Printf("RequiredBool:   '%v'", c.RequiredBool)
	log.Printf("OptionalBool:   '%v' (ptr)", c.OptionalBool)
	if c.OptionalBool != nil {
		log.Printf("OptionalBool:   '%v' (value)", *c.OptionalBool)
	}
	log.Printf("RequiredInt:    '%v'", c.RequiredInt)
	log.Printf("OptionalInt:    '%v' (ptr)", c.OptionalInt)
	if c.OptionalInt != nil {
		log.Printf("OptionalInt:    '%v' (value)", *c.OptionalInt)
	}
	log.Printf("RequiredString: '%s'", c.RequiredString)
	log.Printf("OptionalString: '%s' (ptr)", c.OptionalString)
	if c.OptionalString != nil {
		log.Printf("OptionalString: '%s' (value)", *c.OptionalString)
	}
}
