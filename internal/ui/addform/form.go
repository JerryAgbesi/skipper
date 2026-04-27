// Package addform renders the interactive prompt used by `skipper add`
// when invoked without positional arguments.
package addform

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/huh"
)

type Input struct {
	Alias    string
	User     string
	Hostname string
	Port     string
}

type Result struct {
	Alias     string
	User      string
	Hostname  string
	Port      int
	Cancelled bool
}

// Run displays the interactive add-host form and returns the collected values.
// If the user aborts (Ctrl+C / Esc) Result.Cancelled is true and err is nil.
func Run(in Input) (Result, error) {
	var (
		alias    = in.Alias
		user     = in.User
		hostname = in.Hostname
		port     = in.Port
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Alias").
				Description("Optional. Defaults to host name (or host-port).").
				Placeholder("devone").
				Value(&alias).
				Validate(validateAlias),
			huh.NewInput().
				Title("User").
				Placeholder("root").
				Value(&user).
				Validate(validateUser),
			huh.NewInput().
				Title("HostName").
				Placeholder("10.0.0.8").
				Value(&hostname).
				Validate(validateHostname),
			huh.NewInput().
				Title("Port").
				Placeholder("22 (leave blank to omit)").
				Value(&port).
				Validate(validatePort),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return Result{Cancelled: true}, nil
		}
		return Result{}, err
	}

	parsedPort, err := parsePort(port)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Alias:    strings.TrimSpace(alias),
		User:     strings.TrimSpace(user),
		Hostname: strings.TrimSpace(hostname),
		Port:     parsedPort,
	}, nil
}

func validateAlias(value string) error    { return validateOptionalField("alias", value) }
func validateUser(value string) error     { return validateRequiredField("user", value) }
func validateHostname(value string) error { return validateRequiredField("host name", value) }

func validatePort(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	port, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("port must be a number")
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func parsePort(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func validateRequiredField(name, value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fmt.Errorf("%s is required", name)
	}
	if strings.ContainsFunc(trimmed, unicode.IsSpace) {
		return fmt.Errorf("%s cannot contain whitespace", name)
	}
	return nil
}

func validateOptionalField(name, value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	if strings.ContainsFunc(trimmed, unicode.IsSpace) {
		return fmt.Errorf("%s cannot contain whitespace", name)
	}
	return nil
}
