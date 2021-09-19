package devops

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	DefaultConfirmInputHint = " (only '%s' will be accepted) "
)

type ConfirmOpts struct {
	// Question can optionally be specified for the .Confirm method to
	// print a string before requesting for confirmation. A space will
	// be added at the end of the provided .Question before the
	// string defined in .InputHint is added
	Question string

	// Input defines the input stream to read the input from
	//
	// Defaults to os.Stdin if not specified
	Input io.Reader

	// InputHint is a format string containing a single %s denoting
	// a string that needs to be matched for the confirmation to succeed.
	// Example. " (enter '%s' to continue)", the .Confirm method will
	// populate the %s with the correct matcher value based on
	// .MatchExact or .MatchRegexp
	//
	// Defaults to DefaultConfirmInputHint if not specified
	InputHint string

	// MatchExact defines an exact string match for the confirmation
	// to succeed.
	//
	// When this is defined, MatchRegexp CANNOT be defined
	MatchExact string

	// MatchRegexp defines a regular expression match for the
	// confirmation to succeed.
	//
	// When this is defined, MatchExact CANNOT be defined
	MatchRegexp *regexp.Regexp

	// Output defines the output stream to write output to
	//
	// Defaults to os.Stdin if not specified
	Output io.Writer
}

// SetDefaults checks for unspecified properties which have defaults
// and adds them
func (o *ConfirmOpts) SetDefaults() {
	if o.Input == nil {
		o.Input = os.Stdin
	}
	if o.InputHint == "" {
		o.InputHint = DefaultConfirmInputHint
	}
	if o.Output == nil {
		o.Output = os.Stdout
	}
}

// Validate runs validation checks against the provided options
func (o ConfirmOpts) Validate() error {
	errors := []string{}

	if o.MatchExact == "" && o.MatchRegexp == nil {
		errors = append(errors, "missing Match*")
	}
	if o.MatchExact != "" && o.MatchRegexp != nil {
		errors = append(errors, "only one Match* should be defined")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}
	return nil
}

// Confirm performs a user-terminal-input based confirmation. This
// can be used in situations where it could be useful for a user to
// manually verify a string such as a command to be run
func Confirm(opts ConfirmOpts) (bool, error) {
	opts.SetDefaults()
	if err := opts.Validate(); err != nil {
		return false, fmt.Errorf("failed to trigger confirmation: %s", err)
	}
	if opts.Question != "" {
		opts.Output.Write([]byte(opts.Question))
	}
	isUsingRegexp := opts.MatchRegexp != nil
	acceptedText := opts.MatchExact
	if isUsingRegexp {
		acceptedText = opts.MatchRegexp.String()
	}
	opts.Output.Write([]byte(fmt.Sprintf(opts.InputHint, acceptedText)))
	scanner := bufio.NewScanner(opts.Input)
	if scanner.Scan() {
		input := strings.Trim(scanner.Text(), " \n\t\r")
		if isUsingRegexp {
			return opts.MatchRegexp.Match([]byte(input)), nil
		}
		return opts.MatchExact == input, nil
	} else if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("failed to get user input: %s", err)
	}
	return false, nil
}
