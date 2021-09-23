package devops

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type ConnectionProtocol string

const (
	ConnectionTCP             ConnectionProtocol = "tcp"
	ConnectionUDP             ConnectionProtocol = "udp"
	DefaultConnectionProtocol                    = ConnectionTCP
)

const (
	DefaultTimeout = 3 * time.Second
)

type ValidateConnectionOpts struct {
	Hostname      string
	IsIpV6        bool
	Port          uint16
	Protocol      ConnectionProtocol
	RetryInterval time.Duration
	RetryLimit    uint
	Timeout       time.Duration
}

func (o *ValidateConnectionOpts) SetDefaults() {
	if o.Protocol == "" {
		o.Protocol = DefaultConnectionProtocol
	}

	if o.Timeout == 0 {
		o.Timeout = DefaultTimeout
	}
}

func (o ValidateConnectionOpts) Validate() error {
	errors := []string{}

	if o.Hostname == "" {
		errors = append(errors, "missing hostname")
	}

	if o.Port == 0 {
		errors = append(errors, "missing port")
	}

	if o.Protocol == "" {
		errors = append(errors, "missing protocol")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to validate options: ['%s']", strings.Join(errors, "', '"))
	}

	return nil
}

func ValidateConnection(opts ValidateConnectionOpts) (bool, error) {
	opts.SetDefaults()
	if err := opts.Validate(); err != nil {
		return false, fmt.Errorf("failed to validate connection: %s", err)
	}

	address := net.JoinHostPort(opts.Hostname, strconv.Itoa(int(opts.Port)))
	connection, err := net.DialTimeout(string(opts.Protocol), address, opts.Timeout)
	if err != nil {
		if te, ok := err.(interface{ Temporary() bool }); ok {
			if te.Temporary() {
				if opts.RetryLimit > 0 {
					<-time.After(opts.RetryInterval)
					newOpts := opts
					newOpts.RetryLimit -= 1
					return ValidateConnection(newOpts)
				}
			}
		}
		return false, fmt.Errorf("failed to connect to '%s': %s", address, err)
	}
	if connection != nil {
		defer connection.Close()
		return true, nil
	}

	return false, nil
}
