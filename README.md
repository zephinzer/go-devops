# go-devops

This repository exports a package `devops` that simplifies writing of Go applications as internal tooling "glue".

- [go-devops](#go-devops)
- [Why you might want to use this](#why-you-might-want-to-use-this)
- [Design principles](#design-principles)
- [Usage and Examples](#usage-and-examples)
  - [Commands](#commands)
    - [Listing a directory contents](#listing-a-directory-contents)
    - [Pulling `go` dependencies](#pulling-go-dependencies)
    - [Pulling `node` dependencies](#pulling-node-dependencies)
  - [Input data](#input-data)
    - [Download files](#download-files)
    - [Load environment](#load-environment)
  - [Input validation](#input-validation)
    - [Validating applications](#validating-applications)
    - [Validating connections](#validating-connections)
    - [Validating the environment](#validating-the-environment)
  - [Security](#security)
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

### Listing a directory contents

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

### Pulling `go` dependencies

The following runs `go mod vendor`:

```go
func main() {
  installGoDeps, _ := devops.NewCommand(devops.NewCommandOpts{
    Command: "go",
    Arguments: []string{"mod", "vendor"},
  })
  installGoDeps.Run()
}
```

### Pulling `node` dependencies

The following runs `npm install`:

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

The following downloads the source code from Google:

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

### Load environment

The `.LoadEnvironment` method allows you to load typed environment variables:

```go
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
```

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

| Version   | Changes                                               |
| --------- | ----------------------------------------------------- |
| `v0.0.13` | Formatting fixes  |
| `v0.0.12` | Added `.LoadEnvironment`  |
| `v0.0.11` | Renamed module for being able to import it via its Gitlab URL  |
| `v0.0.10` | Added `.ValidateConnection`                           |
| `v0.0.9`  | Added `.ValidateApplications`                         |
| `v0.0.8`  | Added `.DownloadFile`                                 |
| `v0.0.7`  | Added custom error parsing for `.ValidateEnvironment` |
| `v0.0.6`  | Added `.ValidateEnvironment`                          |
| `v0.0.5`  | Added `.Confirm`                                      |
| `v0.0.4`  | Added inline code comments for documentation          |
| `v0.0.3`  | Added `.GetSshKeyFingerprint`. Also started changelog |

# License

Code is licensed under the MIT license. [See full license here](./LICENSE).
