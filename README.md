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
  - [Input validation](#input-validation)
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

## Input validation

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

## Security

### Retrieving the SSH key fingerprint

```go
func main() {
  keyPath := "./tests/sshkeys/id_rsa_1024.pub"
  fingerprint, err := GetSshKeyFingerprint(GetSshKeyFingerprintOpts{
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

| Version  | Changes                                               |
| -------- | ----------------------------------------------------- |
| `v0.0.6` | Added `.ValidateEnvironment`                          |
| `v0.0.5` | Added `.Confirm`                                      |
| `v0.0.4` | Added inline code comments for documentation          |
| `v0.0.3` | Added `.GetSshKeyFingerprint`. Also started changelog |

# License

Code is licensed under the MIT license. [See full license here](./LICENSE).
