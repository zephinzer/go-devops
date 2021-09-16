package devops

import (
	"bytes"
	"errors"
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

	// fmt.Println(nco.StdoutHooks != nil)
	// fmt.Println(len(nco.StdoutHooks))
	// fmt.Println(!nco.Flag.UseTTY)
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
	currentDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %s", err)
	}

	cmd := exec.Cmd{}

	invocation, err := exec.LookPath(opts.Command)
	if err != nil {
		if !path.IsAbs(invocation) {
			invocation = path.Join(currentDirectory, invocation)
			invocationInfo, err := os.Lstat(invocation)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return nil, fmt.Errorf("failed to find a file system listing at %s", invocation)
				}
				return nil, fmt.Errorf("failed to access file at %s", invocation)
			}
			if invocationInfo.IsDir() {
				return nil, fmt.Errorf("failed to find a file at %s", invocation)
			}
		}
		return nil, fmt.Errorf("failed to find binary '%s' in $PATH: %s", opts.Command, err)
	}
	if strings.Contains(invocation, "/") {
		if !path.IsAbs(invocation) {
			invocation = path.Join(currentDirectory, invocation)
		}
	}
	cmd.Path = invocation

	arguments := []string{opts.Command}
	arguments = append(arguments, opts.Arguments...)
	cmd.Args = arguments

	var workingDir string
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
	cmd.Dir = workingDir

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
	if opts.Flag.HideStdout {
		// this is basically untestsable without extensive instrumentation
		// so it's not tested, but the behaviour should be obvious
		stdoutReader, stdoutWriter = &bytes.Buffer{}, &bytes.Buffer{}
	} else {
		stdoutReader, stdoutWriter = io.Pipe()
		stdoutWriter = io.MultiWriter(os.Stdout, stdoutWriter)
	}
	stdoutWriter = io.MultiWriter(stdoutWriter, &stdoutOutput)
	cmd.Stdout = stdoutWriter

	var stderrReader io.Reader
	var stderrWriter io.Writer
	var stderrOutput bytes.Buffer
	if opts.Flag.HideStderr {
		// this is basically untestsable without extensive instrumentation
		// so it's not tested, but the behaviour should be obvious
		stderrReader, stderrWriter = &bytes.Buffer{}, &bytes.Buffer{}
	} else {
		stderrReader, stderrWriter = io.Pipe()
		stderrWriter = io.MultiWriter(os.Stderr, stderrWriter)
	}
	stderrWriter = io.MultiWriter(stderrWriter, &stderrOutput)
	cmd.Stderr = stderrWriter

	stdoutHooks := opts.StdoutHooks
	stderrHooks := opts.StderrHooks
	stdanyHooks := opts.StdanyHooks

	var stdin io.WriteCloser = nil
	if opts.Flag.UseTTY {
		if len(stdoutHooks) == 0 && len(stderrHooks) == 0 && len(stdanyHooks) == 0 {
			// another non-testable one, any ideas how to verify stdin on
			// an actual terminal?
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
		stdanyHooks:  stdanyHooks,
		stderr:       stderrReader,
		stderrHooks:  stderrHooks,
		stderrOutput: &stderrOutput,
		stdin:        stdin,
		stdout:       stdoutReader,
		stdoutHooks:  stdoutHooks,
		stdoutOutput: &stdoutOutput,
	}, nil
}

type Command interface {
	Bytes() []byte
	GetEnvironment() map[string]string
	GetStderr() []byte
	GetStdout() []byte
	Run() error
	String() string
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

func (c *command) Bytes() []byte {
	var output bytes.Buffer
	output.WriteString(c.Cmd.Path + " ")
	for _, argument := range c.Cmd.Args[1:] {
		output.WriteString("\"" + argument + "\" ")
	}
	return output.Bytes()
}

func (c *command) GetEnvironment() map[string]string {
	envKeyValueMap := map[string]string{}
	for _, envKeyValuePair := range c.Cmd.Env {
		pair := strings.SplitN(envKeyValuePair, "=", 2)
		if len(pair) == 2 {
			envKeyValueMap[pair[0]] = pair[1]
		}
	}
	return envKeyValueMap
}

func (c *command) GetStderr() []byte {
	stderr, _ := ioutil.ReadAll(c.stderrOutput)
	return stderr
}

func (c *command) GetStdout() []byte {
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

func (c *command) String() string {
	return string(c.Bytes())
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
		if onto == nil || len(onto) == 0 {
			continue
		}
		for _, inputHook := range onto {
			if bytes.Contains(incoming, inputHook.On) {
				writer.Write(inputHook.Send)
			}
		}
	}
}
