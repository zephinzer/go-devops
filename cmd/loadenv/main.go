package main

import (
	"fmt"

	"gitlab.com/zephinzer/go-devops"
)

func main() {
	env, err := devops.LoadEnvironment(devops.LoadEnvironmentOpts{
		{
			Key:     "SOME_BOOLEAN_KEY",
			Type:    devops.TypeBool,
			Default: true,
		},
		{
			Key:     "SOME_FLOAT_KEY",
			Type:    devops.TypeFloat,
			Default: 3.142,
		},
		{
			Key:     "SOME_INTEGER_KEY",
			Type:    devops.TypeInt,
			Default: -123456,
		},
		{
			Key:     "SOME_DEFAULT_KEY",
			Default: "hola mundo",
		},
		{
			Key:     "SOME_STRING_KEY",
			Type:    devops.TypeString,
			Default: "hello world",
		},
		{
			Key:     "SOME_UNSIGNED_INTEGER_KEY",
			Type:    devops.TypeUint,
			Default: 123456,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("SOME_BOOLEAN_KEY: %v\n", env.GetBool("SOME_BOOLEAN_KEY"))
	fmt.Printf("SOME_FLOAT_KEY: %v\n", env.GetFloat("SOME_FLOAT_KEY"))
	fmt.Printf("SOME_DEFAULT_KEY: %v\n", env.Get("SOME_DEFAULT_KEY"))
	fmt.Printf("SOME_INTEGER_KEY: %v\n", env.GetInt("SOME_INTEGER_KEY"))
	fmt.Printf("SOME_STRING_KEY: %v\n", env.GetString("SOME_STRING_KEY"))
	fmt.Printf("SOME_UNSIGNED_INTEGER_KEY: %v\n", env.GetUint("SOME_UNSIGNED_INTEGER_KEY"))
}
