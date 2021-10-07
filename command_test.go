package devops

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommandTests struct {
	suite.Suite
}

func TestCommand(t *testing.T) {
	suite.Run(t, &CommandTests{})
}

// Test_Basic checks for basic expectations of a command concept
// like invocation with arguments
func (s CommandTests) Test_Basic() {
	scriptPath := "tests/command/basic.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	command, err := NewCommand(NewCommandOpts{
		Command:   scriptPath,
		Arguments: []string{"arg1", "arg2", "arg3"},
	})
	s.Nil(err, "this command should be created successfully")
	err = command.Run()
	s.Nil(err, "the command should have ran successfully")

	stdout := string(command.GetStdout())
	s.NotEmpty(stdout, "the stdout output should not be empty")
	s.Contains(stdout, "$@: arg1 arg2 arg3")
	s.Contains(stdout, scriptPath)
	s.Contains(stdout, "$1: arg1")
	s.Contains(stdout, "$2: arg2")
	s.Contains(stdout, "$3: arg3")
	s.Contains(stdout, "this prints to stdout")
	stderr := string(command.GetStderr())
	s.NotEmpty(stderr, "the stderr output should not be empty")
	s.Contains(stderr, "this prints to stderr")
}

func (s CommandTests) Test_Environment() {
	scriptPath := "tests/command/env.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	envKey1 := "ENV_VALUE_1"
	envKey2 := "ENV_VALUE_2"
	envValue1 := "say"
	envValue2 := "hello world"
	command, err := NewCommand(NewCommandOpts{
		Command: scriptPath,
		Environment: map[string]string{
			envKey1: envValue1,
			envKey2: envValue2,
		},
	})
	s.Nil(err, "this command should be created successfully")
	err = command.Run()
	s.Nil(err, "the command should have ran successfully")
	output := string(command.GetStdout())
	s.NotEmpty(output, "the output should not be empty")
	s.Contains(output, fmt.Sprintf("$%s:%s", envKey1, envValue1))
	s.Contains(output, fmt.Sprintf("$%s:%s", envKey2, envValue2))
}

func (s CommandTests) Test_EnvironmentGlobal() {
	scriptPath := "tests/command/env.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	command, err := NewCommand(NewCommandOpts{
		Command: scriptPath,
		Flag: CommandFlagset{
			UseGlobalEnvironment: true,
		},
	})
	s.Nil(err, "this command should be created successfully")
	err = command.Run()
	s.Nil(err, "the command should have ran successfully")
	env := command.GetEnvironment()
	_, ok := env["USER"] // should be a universal environment key?
	s.True(ok)
}

func (s CommandTests) Test_Pwd() {
	scriptPath := "tests/command/pwd.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	targetDirectory := "./tests"
	wd, _ := os.Getwd()
	expectedResolvedDirectory := path.Join(wd, targetDirectory)
	command, err := NewCommand(NewCommandOpts{
		Command:    scriptPath,
		WorkingDir: targetDirectory,
	})
	s.Nil(err, "this command should be created successfully but failed with: %s", err)
	err = command.Run()
	s.Nil(err, "the command should have ran successfully but failed with: %s", err)
	output := string(command.GetStdout())
	s.NotEmpty(output, "the output should not be empty")
	s.Contains(output, fmt.Sprintf("$(pwd):%s", expectedResolvedDirectory))
}

func (s CommandTests) Test_StdanyHook() {
	scriptPath := "tests/command/stderr.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	expectedText := "for stdany"
	command, err := NewCommand(NewCommandOpts{
		Command: scriptPath,
		Flag: CommandFlagset{
			UseTTY: true,
		},
		StdanyHooks: InputHooks{
			InputHook{
				On:   []byte("prints to stderr"),
				Send: []byte(expectedText + "\n"),
			},
		},
	})
	s.Nil(err, "this command should be created successfully but failed with: %s", err)
	err = command.Run()
	s.Nil(err, "the command should have ran successfully but failed with: %s", err)

	output := string(command.GetStderr())
	s.NotEmpty(output, "the stderr output should not be empty")
	s.Contains(output, expectedText)
}

func (s CommandTests) Test_StderrHook() {
	scriptPath := "tests/command/stderr.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	expectedText := "for stderr"
	command, err := NewCommand(NewCommandOpts{
		Command: scriptPath,
		Flag: CommandFlagset{
			UseTTY: true,
		},
		StderrHooks: InputHooks{
			InputHook{
				On:   []byte("prints to stderr"),
				Send: []byte(expectedText + "\n"),
			},
		},
	})
	s.Nil(err, "this command should be created successfully but failed with: %s", err)
	err = command.Run()
	s.Nil(err, "the command should have ran successfully but failed with: %s", err)

	output := string(command.GetStderr())
	s.NotEmpty(output, "the stderr output should not be empty")
	s.Contains(output, expectedText)
}

func (s CommandTests) Test_StdoutHook() {
	scriptPath := "tests/command/stdout.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.True(scriptPathInfo.Mode().Perm()&0111 > 0, "script should be executable in the first place")

	expectedText := "for stdout"
	command, err := NewCommand(NewCommandOpts{
		Command: scriptPath,
		Flag: CommandFlagset{
			UseTTY: true,
		},
		StdoutHooks: InputHooks{
			InputHook{
				On:   []byte("prints to stdout"),
				Send: []byte(expectedText + "\n"),
			},
		},
	})
	s.Nil(err, "this command should be created successfully but failed with: %s", err)
	err = command.Run()
	s.Nil(err, "the command should have ran successfully but failed with: %s", err)

	stdout := string(command.GetStdout())
	s.NotEmpty(stdout, "the stdout output should not be empty")
	s.Contains(stdout, expectedText)
}

func (s CommandTests) Test_command_String() {
	scriptPath := "tests/command/stdout.sh"
	scriptPathInfo, err := os.Lstat(scriptPath)
	s.Nil(err, "file should exist and be accessible")
	s.False(scriptPathInfo.IsDir(), "path should not be to a directory")

	arguments := []string{"1", "2", "3"}
	command, err := NewCommand(NewCommandOpts{
		Command:   scriptPath,
		Arguments: arguments,
	})
	s.Nil(err)
	s.Contains(command.String(), scriptPath+" \""+strings.Join(arguments, "\" \"")+"\"")
}

func (s CommandTests) Test_NewCommand_pathNotFound() {
	command, err := NewCommand(NewCommandOpts{
		Command: "thisbinarydoesnotexist",
	})
	s.NotNil(err)
	s.Nil(command)
}

func (s CommandTests) Test_NewCommandOpts_Validate() {
	opts := NewCommandOpts{}
	err := opts.Validate()
	s.Contains(err.Error(), ".Command")
	opts.Command = "test"

	opts.StdanyHooks = InputHooks{InputHook{}}
	err = opts.Validate()
	s.Contains(err.Error(), ".Flag.UseTTY")
	opts.StdanyHooks = nil

	opts.StderrHooks = InputHooks{InputHook{}}
	err = opts.Validate()
	s.Contains(err.Error(), ".Flag.UseTTY")
	opts.StderrHooks = nil

	opts.StdoutHooks = InputHooks{InputHook{}}
	err = opts.Validate()
	s.Contains(err.Error(), ".Flag.UseTTY")
	opts.StdoutHooks = nil

}
