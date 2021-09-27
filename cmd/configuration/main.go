package main

import (
	"log"
	"os"
	"strings"

	"gitlab.com/zephinzer/go-devops"
)

type configuration struct {
	RequiredStringSlice []string  `default:"a,b,c" delimiter:","`
	RequiredString      string    `default:"hello world"`
	RequiredInt         int       `default:"1"`
	RequiredBool        bool      `default:"true"`
	OptionalString      *string   `default:"hola mundo"`
	OptionalStringSlice *[]string `default:"d,e,f" delimiter:","`
	OptionalInt         *int      `default:"2"`
	OptionalBool        *bool     `default:"true"`
}

func main() {
	c := configuration{}
	if err := devops.LoadConfiguration(&c); err != nil {
		log.Println(err)
		os.Exit(err.(devops.LoadConfigurationError).Code)
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
	log.Printf("RequiredStringSlice: ['%s'] (len: %v)", strings.Join(c.RequiredStringSlice, "', '"), len(c.RequiredStringSlice))
	log.Printf("OptionalStringSlice: '%s' (ptr)", c.OptionalStringSlice)
	if c.OptionalStringSlice != nil {
		log.Printf("OptionalStringSlice: ['%s'] (value) (len: %v)", strings.Join(*c.OptionalStringSlice, "', '"), len(*c.OptionalStringSlice))
	}
}
