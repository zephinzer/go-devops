package devops

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// InputHooks is a convenience reference for a slice of InputHook instances
type InputHooks []InputHook

// InputHook defines a structure for responding to a byte sequence
// of `.On` using the byte sequence .Send. Receiver and signaller is
// left to the controller to implement
type InputHook struct {
	On   []byte
	Send []byte
}

// CommandFlagSet defines a set of boolean configuration flags for the
// Command class
type CommandFlagset struct {
	// HideStdout indicates whether STDOUT should be printed to the terminal
	HideStdout bool
	// HideStderr indicates whether STDERR should be printed to the terminal
	HideStderr bool
	// UseGlobalEnvironment indicates whether the child process should
	// inherit the parent's environment
	UseGlobalEnvironment bool
	// UseTTY enables use of STDIN
	UseTTY bool
}

// NewCommandOpts defines a set of options for use with the `NewCommand()`
// intiializer method
type NewCommandOpts struct {
	// Arguments lists the parameters to pass to the `.Command`
	Arguments []string

	// Command is an invocation that is in the $PATH OR a path to the
	// binary to execute (relative paths will be resolved to absolute)
	Command string

	// Environment is a mapping of key=value that will be injected into
	// the child process's environment. Use the `.Flag.UseGlobalEnvironment`
	// configuration flag to inject the parent environment into the child's
	Environment map[string]string

	// WorkingDir indicates the working directory of the child process.
	// If not an absolute path, this will be resolved to its absolute
	// one before the process begins
	WorkingDir string

	// StdoutHooks allows you to send a []byte data structure to STDIN
	// when receiving a predefined string from STDOUT. The `.Flag.UseTTY`
	// has to be enabled for this to work
	//
	// NOTE: If you have defined any `StdanyHooks`, those take execution
	// precedence
	StdoutHooks InputHooks

	// StderrHooks allows you to send a []byte data structure to STDIN
	// when receiving a predefined string from STDERR. The `.Flag.UseTTY`
	// has to be enabled for this to work
	//
	// NOTE: If you have defined any `StdanyHooks`, those take execution
	// precedence
	StderrHooks InputHooks

	// StdanyHooks allows you to send a []byte data structure to STDIN
	// when receiving a predefined string from both STDOUT and STDERR.
	// The `.Flag.UseTTY` has to be enabled for this to work
	//
	// NOTE: If you have defined any `StdoutHooks` or `StderrHooks` that
	// overlap with hooks defined in `StdanyHooks`, the hooks from
	// `StdanyHooks` will be executed first
	StdanyHooks InputHooks

	// Flag defines a boolean configuration flagset
	Flag CommandFlagset
}

// Validate returns an error if a combination of the provided options will
// cause problems during execution or just plain invalid
func (nco NewCommandOpts) Validate() error {
	errors := []string{}

	if nco.Command == "" {
		errors = append(errors, "missing .Command")
	}

	if nco.StdoutHooks != nil && len(nco.StdoutHooks) > 0 && !nco.Flag.UseTTY {
		errors = append(errors, ".Flag.UseTTY should be true if .StdoutHooks is defined")
	}

	if nco.StderrHooks != nil && len(nco.StderrHooks) > 0 && !nco.Flag.UseTTY {
		errors = append(errors, ".Flag.UseTTY should be true if .StderrHooks is defined")
	}

	if nco.StdanyHooks != nil && len(nco.StdanyHooks) > 0 && !nco.Flag.UseTTY {
		errors = append(errors, ".Flag.UseTTY should be true if .StdanyHooks is defined")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate NewCommandOpts: ['%s']", strings.Join(errors, "', '"))
	}
	return nil
}

func NewCommand(opts NewCommandOpts) (Command, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create Command: %s", err)
	}
	cmd := exec.Cmd{}

	invocation, err := exec.LookPath(opts.Command)
	if err != nil {
		return nil, fmt.Errorf("failed to find binary '%s' in $PATH: %s", opts.Command, err)
	}
	cmd.Path = invocation

	arguments := []string{opts.Command}
	arguments = append(arguments, opts.Arguments...)
	cmd.Args = arguments

	var workingDir string
	currentDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %s", err)
	}
	if opts.WorkingDir != "" {
		workingDir = opts.WorkingDir
		if !path.IsAbs(workingDir) {
			workingDir = path.Join(currentDirectory, workingDir)
		}
		fileInfo, err := os.Lstat(workingDir)
		if err != nil {
			return nil, fmt.Errorf("failed to get information about path '%s': %s", workingDir, err)
		}
		if !fileInfo.IsDir() {
			return nil, fmt.Errorf("failed to find a directory at path '%s'", workingDir)
		}
	} else {
		workingDir = currentDirectory
	}

	environment := []string{}
	for key, value := range opts.Environment {
		environment = append(environment, fmt.Sprintf("%s=%s", key, value))
	}
	if opts.Flag.UseGlobalEnvironment {
		globalEnvironment := os.Environ()
		for _, globalEnvironmentKeyValue := range globalEnvironment {
			environment = append(environment, globalEnvironmentKeyValue)
		}
	}
	cmd.Env = environment

	var stdoutReader io.Reader
	var stdoutWriter io.Writer
	var stdoutOutput bytes.Buffer
	stdoutReader, stdoutWriter = io.Pipe()
	stdoutWriter = io.MultiWriter(stdoutWriter, &stdoutOutput)
	if !opts.Flag.HideStdout {
		stdoutWriter = io.MultiWriter(os.Stdout, stdoutWriter)
	}
	cmd.Stdout = stdoutWriter

	var stderrReader io.Reader
	var stderrWriter io.Writer
	var stderrOutput bytes.Buffer
	stderrReader, stderrWriter = io.Pipe()
	stderrWriter = io.MultiWriter(stderrWriter, &stderrOutput)
	if !opts.Flag.HideStderr {
		stderrWriter = io.MultiWriter(os.Stderr, stderrWriter)
	}
	cmd.Stderr = stderrWriter

	stdoutHooks := opts.StdoutHooks
	stderrHooks := opts.StderrHooks
	stdanyHooks := opts.StdanyHooks

	var stdin io.WriteCloser = nil
	if opts.Flag.UseTTY {
		if len(stdoutHooks) == 0 && len(stderrHooks) == 0 && len(stdanyHooks) == 0 {
			cmd.Stdin = os.Stdin
		} else {
			stdin, err = cmd.StdinPipe()
			if err != nil {
				return nil, fmt.Errorf("failed to provision a tty: %s", err)
			}
		}
	}

	return &command{
		Cmd:          cmd,
		stdout:       stdoutReader,
		stdoutHooks:  stdoutHooks,
		stdoutOutput: &stdoutOutput,
		stderr:       stderrReader,
		stderrHooks:  stderrHooks,
		stderrOutput: &stderrOutput,
		stdanyHooks:  stdanyHooks,
		stdin:        stdin,
	}, nil
}

type Command interface {
	Run() error
	GetOutput() []byte
	GetError() []byte
}

type command struct {
	exec.Cmd
	stdout       io.Reader
	stdoutHooks  InputHooks
	stdoutOutput io.Reader
	stderr       io.Reader
	stderrHooks  InputHooks
	stderrOutput io.Reader
	stdanyHooks  InputHooks
	stdin        io.WriteCloser
}

func (c *command) GetError() []byte {
	stderr, _ := ioutil.ReadAll(c.stderrOutput)
	return stderr
}

func (c *command) GetOutput() []byte {
	stdout, _ := ioutil.ReadAll(c.stdoutOutput)
	return stdout
}

func (c *command) Run() error {
	if c.stdin != nil {
		defer c.stdin.Close()
	}
	go c.hook(c.stdout, c.stdoutHooks, c.stdin)
	go c.hook(c.stderr, c.stderrHooks, c.stdin)
	go func() {
		for {
			<-time.After(200 * time.Millisecond)
			if c.Cmd.ProcessState.ExitCode() > 0 {
				c.stdout.(*io.PipeReader).Close()
				c.stderr.(*io.PipeReader).Close()
				return
			}
		}
	}()
	return c.Cmd.Run()
}

func (c *command) hook(from io.Reader, onto InputHooks, writer io.WriteCloser) {
	incoming := make([]byte, 64)
	for {
		_, err := from.Read(incoming)
		if err == io.EOF {
			return
		}
		for _, inputHook := range c.stdanyHooks {
			if bytes.Contains(incoming, inputHook.On) {
				writer.Write(inputHook.Send)
			}
		}
		if onto == nil {
			continue
		}
		for _, inputHook := range onto {
			if bytes.Contains(incoming, inputHook.On) {
				writer.Write(inputHook.Send)
			}
		}
	}
}
