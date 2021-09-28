package devops

import (
	"reflect"
	"unsafe"

	"github.com/zephinzer/go-strcase"
)

func newConfiguration(from interface{}) configuration {
	configValue := reflect.ValueOf(from)
	configType := reflect.TypeOf(from)
	if configValue.Kind() == reflect.Ptr {
		configType = configType.Elem()
	}
	configFields := []configurationField{}
	if configType.Kind() == reflect.Struct {
		for i := 0; i < configType.NumField(); i++ {
			field := configType.Field(i)
			configField := configurationField{
				Name:  field.Name,
				Tag:   field.Tag,
				Type:  field.Type,
				Value: configValue.Elem().FieldByName(field.Name),
			}
			configFields = append(configFields, configField)
		}
	} else {
		configFields = append(configFields, configurationField{
			Name:  "_",
			Tag:   "",
			Type:  configType,
			Value: configValue,
		})
	}
	return configuration{
		Value:  configValue,
		Type:   configType,
		Fields: configFields,
	}
}

type configuration struct {
	Value  reflect.Value
	Type   reflect.Type
	Fields []configurationField // map[string]reflect.Type //
}

func (c configuration) IsPointer() bool {
	return c.Value.Kind() == reflect.Ptr
}

func (c configuration) IsStruct() bool {
	return c.Type.Kind() == reflect.Struct
}

type configurationField struct {
	Name  string
	Tag   reflect.StructTag
	Type  reflect.Type
	Value reflect.Value
}

func (c configurationField) getPointer() unsafe.Pointer {
	/* #nosec - this is required to assign values to an untyped struct */
	return unsafe.Pointer(c.Value.UnsafeAddr())
}

func (c configurationField) GetDefaultValue() *string {
	if v, ok := c.Tag.Lookup("default"); ok {
		return &v
	}
	return nil
}

func (c configurationField) GetEnvironmentKey() string {
	if v, ok := c.Tag.Lookup("env"); ok {
		return v
	}
	return strcase.ToUpperSnake(c.Name)
}

func (c configurationField) SetBool(value bool) {
	c.Value.SetBool(value)
}

func (c configurationField) SetBoolPointer(value bool) {
	reflect.NewAt(c.Value.Type(), c.getPointer()).Elem().Set(reflect.ValueOf(&value))
}

func (c configurationField) SetInt(value int) {
	c.Value.SetInt(int64(value))
}

func (c configurationField) SetIntPointer(value int) {
	reflect.NewAt(c.Value.Type(), c.getPointer()).Elem().Set(reflect.ValueOf(&value))
}

func (c configurationField) SetString(value string) {
	c.Value.SetString(value)
}

func (c configurationField) SetStringPointer(value string) {
	reflect.NewAt(c.Value.Type(), c.getPointer()).Elem().Set(reflect.ValueOf(&value))
}

func (c configurationField) SetStringSlice(value []string) {
	reflect.NewAt(c.Value.Type(), c.getPointer()).Elem().Set(reflect.ValueOf(value))
}

func (c configurationField) SetStringSlicePointer(value []string) {
	reflect.NewAt(c.Value.Type(), c.getPointer()).Elem().Set(reflect.ValueOf(&value))
}
