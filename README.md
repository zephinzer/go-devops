# go-devops

This repository exports a package `devops` that simplifies writing of Go applications as internal tooling "glue".

- [go-devops](#go-devops)
- [Why you might want to use this](#why-you-might-want-to-use-this)
- [Example usage](#example-usage)
  - [Listing a directory contents](#listing-a-directory-contents)
  - [Pulling `go` dependencies](#pulling-go-dependencies)
  - [Pulling `node` dependencies](#pulling-node-dependencies)

# Why you might want to use this

1. You spend your time writing shell scripts and are sick of having untested code
2. You are in a DevOps team moving towards a product way of doing things

# Example usage

## Listing a directory contents

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

## Pulling `go` dependencies

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

## Pulling `node` dependencies

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
