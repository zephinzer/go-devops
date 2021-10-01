# go-devops

This repository exports a package `devops` that simplifies writing of Go applications as internal tooling "glue".

- [go-devops](#go-devops)
- [Why you might want to use this](#why-you-might-want-to-use-this)
- [Design principles](#design-principles)
- [Usage and Examples](#usage-and-examples)
  - [Commands](#commands)
    - [Running a command](#running-a-command)
  - [Input data](#input-data)
    - [Download files](#download-files)
    - [Get data from a HTTP endpoint](#get-data-from-a-http-endpoint)
    - [Load configuration](#load-configuration)
      - [Notes on loading configuration](#notes-on-loading-configuration)
  - [Input validation](#input-validation)
    - [Validating applications](#validating-applications)
    - [Validating connections](#validating-connections)
    - [Validating the environment](#validating-the-environment)
  - [Security](#security)
    - [Generating an SSH keypair](#generating-an-ssh-keypair)
    - [Retrieving the SSH key fingerprint](#retrieving-the-ssh-key-fingerprint)
  - [User interactions](#user-interactions)
    - [Confirmation dialog](#confirmation-dialog)
- [Changelog](#changelog)
- [License](#license)

# Why you might want to use this

1. You spend your time writing shell scripts and are sick of having untestable code (without significant effort)
2. You are in a DevOps team moving towards a product way of doing things and have picked up Go an want to rewrite your shell scripts using Go

# Design principles

1. All `New[.]*` functions will return an `interface` as far as possible, while this could hide data, it also prevents state related errors from modifications after initialisation. This is in turn supported by validation checks that run during the initialisation process
2. All `New[.]*` functions will perform a sanity check on provided options and return an `error` if checks are not successful. While this could be annoying, this encourages lazy-instantiation so that assigned properties do not become stale
3. Rather than just providing methods to run a function, which would easily solve problems addressed above, we require a constructor for most objects via a method named `New[.]*` to allow for passing the instance to another controller, which means with this separation you can also separate your data access/creation and controller code by passing an instance to a controller for processing

# Usage and Examples

All examples assume the importing of this package using:

```go
// ...
import "gitlab.com/zephinzer/go-devops"
// ...
```

## Commands

### Running a command

> A working example is available at [`./cmd/command`](./cmd/command)

The following runs `ls -al`:

```go
func main() {
  ls, _ := devops.NewCommand(devops.NewCommandOpts{
    Command: "ls",
    Arguments: []string{"-a", "-l"},
  })
  ls.Run()
}
```

The following runs `go mod vendor` and pulls in dependencies for a Go project:

```go
func main() {
  installGoDeps, _ := devops.NewCommand(devops.NewCommandOpts{
    Command: "go",
    Arguments: []string{"mod", "vendor"},
  })
  installGoDeps.Run()
}
```

The following runs `npm install` and pulls in dependencies for a Node project:

```go
func main() {
  installNodeDeps, _ := devops.NewCommand(devops.NewCommandOpts{
    Command: "npm",
    Arguments: "install",
  })
  installNodeDeps.Run()
}
```

## Input data

### Download files

> A working example is available at [`./cmd/download`](./cmd/download)

The `.DownloadFile` method downloads the source code from Google into a specified `DestinationPath`:

```go
func main() {
	targetURL, err := url.Parse("https://google.com")
  if err != nil {
    panic(err)
  }
	if err = devops.DownloadFile(DownloadFileOpts{
		DestinationPath: "./google.com.src.txt",
		URL:             targetURL,
	}); err != nil {
    panic(err)
  }
}
```

### Get data from a HTTP endpoint

> A working example is available at [`./cmd/curl`](./cmd/curl)

The `.SendHTTPRequest` method can be used in place of `cURL` to make a HTTP request:

```go
func main() {
	targetURL, err := url.Parse("https://httpbin.org/uuid")
	if err != nil {
		panic(err)
	}
	response, err := devops.SendHTTPRequest(devops.SendHTTPRequestOpts{
		URL: targetURL,
	})
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("uuid: %s", string(responseBody))
}
```

`.SendHTTPRequest` supports all common `curl` flags via the `SendHTTPRequestOpts` object.

### Load configuration

> A working example is available at [`./cmd/configuration`](./cmd/configuration)

The `.LoadConfiguration` method allows you to load from environment variables using your own `struct` definition:

```go
type configuration struct {
  // CustomEnvString will be read using os.Getenv("USE_THIS_INSTEAD")
  CustomEnvString     string `env:"USE_THIS_INSTEAD"`
  // RequiredStringSlice will be read using os.Getenv("REQUIRED_STRING_SLICE")
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
    // use it like an error
    log.Println(err)

    // consolidated errors
    errs := err.(devops.LoadConfigurationErrors)
    log.Printf("error code   : %v", errs.GetCode())
    log.Printf("error message: %s", errs.GetMessage())

    // individual errors
    log.Println("errors follow")
    for _, errInstance := range err.(devops.LoadConfigurationErrors) {
      log.Printf("code   : %v", errInstance.Code)
      log.Printf("message: %s", errInstance.Message)
    }

		os.Exit(errs.GetCode())
	}
}
```

#### Notes on loading configuration

1. Property names are automagically converted to `UPPER_SNAKE_CASE` and these are used to load values from the environment using `os.Getenv`
2. To define a custom environment key for the property, use the `env:"READ_FROM_THIS_INSTEAD"` struct tag
3. To define a default value for the property, use the `default:"default value"` struct tag
4. To indiciate a configuration property is **REQUIRED**, specify the type as a `value` type. If the environment does not contain the environment key, an error is returned
5. To indiciate a configuration property is **OPTIONAL**, specify the type as a `*pointer` type. If the environment does not contain the environment key, the value is set to `nil`
6. When defining a slice of strings, use the `delimiter:","` struct tag to define the character sequence used to indicate boundaries between sequential strings
7. The returned `error` can be type-asserted into a `LoadConfigurationErrors` structure which provides both a `GetCode()` and a `GetMessage()` method you can use for assessing errors, you could `range` through it to get individual errors or just call `.Error()` to get a collated error message

## Input validation

### Validating applications

The `.ValidateApplications` function can be used to validate that paths provided are executable or in the system's `$PATH` variable.

A full example follows:

```go
func main() {
	err := devops.ValidateApplications(ValidateApplicationsOpts{
		Paths: []string{"thisappdoesnotexist"},
	})
  if err != nil {
    if _, ok := err.(devops.ValidateApplicationsErrors); ok {
      panic(fmt.Sprintf("failed to find applications: ['%s']", strings.Join(err.Errors, "', '")))
    }
  }
}
```

### Validating connections

The `.ValidateConnection` function can be used to validate that a provided hostname and port is reachable and listening for requests.

A full example follows:

```go
func main() {
  err := devops.ValidateConnection(ValidateConnectionOpts{
    Hostname: "google.com",
    Port: 80,
  })
  if err != nil {
    panic(err)
  }
}
```

### Validating the environment

The `.ValidateEnvironment` can be used to validate that certain keys of interest are defined in the enviornment and returns an error if it doesn't.

A full example follows:

```go
func main() {
	err := ValidateEnvironment(ValidateEnvironmentOpts{
		Keys: EnvironmentKeys{
			{Name: "STRING", Type: TypeString},
			{Name: "INT", Type: TypeInt},
			{Name: "UINT", Type: TypeUint},
			{Name: "FLOAT", Type: TypeFloat},
			{Name: "BOOL", Type: TypeBool},
			{Name: "ANY", Type: TypeAny},
		},
	})
  if err != nil {
    panic(err)
  }
}
```

If the `Type` property is not set, it defaults to `TypeAny`

For custom parsing of error, you can do a type assertion on the `error` interface to `ValidateEnvironmentErrors` and retrieve the error keys/types using the `.Errors` property:

```go
func main() {
  err := devops.ValidateEnvironment(ValidateEnvironmentOpts{
    Keys: EnvironmentKeys{
			{Name: "STRING", Type: TypeString},
			{Name: "INT", Type: TypeInt},
			{Name: "UINT", Type: TypeUint},
			{Name: "FLOAT", Type: TypeFloat},
			{Name: "BOOL", Type: TypeBool},
			{Name: "ANY", Type: TypeAny},
    },
  })
  if err != nil {
    errs, _ := err.(devops.ValidateEnvironmentErrors)
    for _, errInstance := range errs.Errors {
      fmt.Printf(
        "key[%s] errored (expected type: %s, observed value: %s)",
        errInstance.Key,
        errInstance.ExpectedType,
        errInstance.Value,
      )
    }
  }
}
```

## Security

### Generating an SSH keypair

To generate an SSH keypair, you can use the `.NewSSHKeypair` function.

```go
func main() {
  keypair, err := NewSSHKeypair(NewSSHKeypairOpts{
    Bytes: 4096,
  })
  if err != nil {
    panic(err)
  }
  // this prints the keys, you can write it to a file instead
  fmt.Printf("private key: %s\n", string(keypair.Private))
  fmt.Printf("public key : %s\n", string(keypair.Public))
}
```


### Retrieving the SSH key fingerprint

```go
func main() {
  keyPath := "./tests/sshkeys/id_rsa_1024.pub"
  fingerprint, err := devops.GetSshKeyFingerprint(devops.GetSshKeyFingerprintOpts{
    IsPublicKey: true,
    Path:        keyPath,
  })

  fmt.Printf("md5 hash   : %s\n", fingerprint.GetMD5())
  // above outputs 'aa:bb:cc:dd ...'

  fmt.Printf("sha256 hash: %s\n", fingerprint.GetSHA256())
  // above outputs 'sha256 hash: SHA256:AbCdEf ...'
}
```

To run this on a private key, set the `IsPublicKey` to `false` (or leave it unset) and set `IsPrivateKey` property to true.

To specify a password, set the `Passphrase` property of the `GetSshKeyFingerprintOpts` instance.

## User interactions

### Confirmation dialog

To trigger a confirmation dialog in the terminal with the user, use the `.Confirm` method.

> A working example is available at [`./cmd/confirm`](./cmd/confirm)

```go
func main() {
  yes, err := devops.Confirm(devops.ConfirmOpts{
    Question:   "exact match",
    MatchExact: "yes",
  })
  if err != nil {
    log.Fatalf("failed to get user input: %s", err)
  }
  log.Printf("user confirmed: %v\n", yes)
}
```

# Changelog

| Version   | Changes                                                                                                                                 |
| --------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| `v0.2.3`  | Added `.SendHTTPRequest`, improved inline documentation                                                                                                 |
| `v0.2.2`  | Added `.NewSSHKeypair`                                                                                                 |
| `v0.2.1`  | Fixed issues coming from `gosec`                                                                                                        |
| `v0.2.0`  | Updated `error` return of `.LoadConfiguration` to return `LoadConfigurationErrors` instead so that all errors can be made known at once |
| `v0.1.0`  | **Removed `.LoadEnvironment`** and added `.LoadConfiguration` which is a better and cleaner way of doing things                         |
| `v0.0.13` | Formatting fixes                                                                                                                        |
| `v0.0.12` | Added `.LoadEnvironment`                                                                                                                |
| `v0.0.11` | Renamed module for being able to import it via its Gitlab URL                                                                           |
| `v0.0.10` | Added `.ValidateConnection`                                                                                                             |
| `v0.0.9`  | Added `.ValidateApplications`                                                                                                           |
| `v0.0.8`  | Added `.DownloadFile`                                                                                                                   |
| `v0.0.7`  | Added custom error parsing for `.ValidateEnvironment`                                                                                   |
| `v0.0.6`  | Added `.ValidateEnvironment`                                                                                                            |
| `v0.0.5`  | Added `.Confirm`                                                                                                                        |
| `v0.0.4`  | Added inline code comments for documentation                                                                                            |
| `v0.0.3`  | Added `.GetSshKeyFingerprint`. Also started changelog                                                                                   |

# License

Code is licensed under the MIT license. [See full license here](./LICENSE).
