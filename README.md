# go-devops

This repository exports a package `devops` that simplifies writing of Go applications as internal tooling "glue".

- [go-devops](#go-devops)
- [Why you might want to use this](#why-you-might-want-to-use-this)
- [Design principles](#design-principles)
- [Example usage](#example-usage)
  - [Commands](#commands)
    - [Listing a directory contents](#listing-a-directory-contents)
    - [Pulling `go` dependencies](#pulling-go-dependencies)
    - [Pulling `node` dependencies](#pulling-node-dependencies)
- [Security](#security)
  - [SSH Keys](#ssh-keys)
    - [Retrieving the SSH key fingerprint](#retrieving-the-ssh-key-fingerprint)
- [Changelog](#changelog)
- [License](#license)

# Why you might want to use this

1. You spend your time writing shell scripts and are sick of having untestable code (without significant effort)
2. You are in a DevOps team moving towards a product way of doing things and have picked up Go an want to rewrite your shell scripts using Go

# Design principles

1. All `New[.]*` functions will return an `interface` as far as possible, while this could hide data, it also prevents state related errors from modifications after initialisation. This is in turn supported by validation checks that run during the initialisation process
2. All `New[.]*` functions will perform a sanity check on provided options and return an `error` if checks are not successful. While this could be annoying, this encourages lazy-instantiation so that assigned properties do not become stale
3. Rather than just providing methods to run a function, which would easily solve problems addressed above, we require a constructor for most objects via a method named `New[.]*` to allow for passing the instance to another controller, which means with this separation you can also separate your data access/creation and controller code by passing an instance to a controller for processing

# Example usage

## Commands

### Listing a directory contents

> A working example is available at [`./cmd/ls`](./cmd/ls)

The following runs `ls -al`:

```go
import "gitlab.com/zephinzer/go-devops"

// ...

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
import "gitlab.com/zephinzer/go-devops"

// ...

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
import "gitlab.com/zephinzer/go-devops"

// ...

func main() {
  installNodeDeps, _ := devops.NewCommand(devops.NewCommandOpts{
    Command: "npm",
    Arguments: "install",
  })
  installNodeDeps.Run()
}
```

# Security

## SSH Keys

### Retrieving the SSH key fingerprint

```go
import "gitlab.com/zephinzer/go-devops"

// ...

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

# Changelog

| Version  | Changes                                              |
| -------- | ---------------------------------------------------- |
| `v0.0.3` | Added SSH fingerprinting API. Also started changelog |

# License

Code is licensed under the MIT license. [See full license here](./LICENSE).
