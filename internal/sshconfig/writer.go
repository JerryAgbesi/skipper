package sshconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func AddHost(path string, host Host) (*Host, error) {
	if strings.TrimSpace(host.Hostname) == "" {
		return nil, fmt.Errorf("host name is required")
	}

	if strings.TrimSpace(host.User) == "" {
		return nil, fmt.Errorf("user is required")
	}

	host.Alias = resolveAlias(host)
	if err := validateHostFields(host); err != nil {
		return nil, err
	}

	existingHosts, err := readExistingHosts(path)
	if err != nil {
		return nil, err
	}

	for _, existingHost := range existingHosts {
		if existingHost.Alias != host.Alias {
			continue
		}

		if existingHost.Hostname == host.Hostname && existingHost.User == host.User && existingHost.Port == host.Port {
			return &existingHost, nil
		}

		return nil, fmt.Errorf("host %q already exists in %s with different settings", host.Alias, path)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	currentContent, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read config file %q: %w", path, err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %q: %w", path, err)
	}
	defer file.Close()

	if len(currentContent) > 0 && !strings.HasSuffix(string(currentContent), "\n") {
		if _, err := file.WriteString("\n"); err != nil {
			return nil, fmt.Errorf("failed to prepare config file %q: %w", path, err)
		}
	}

	if len(strings.TrimSpace(string(currentContent))) > 0 {
		if _, err := file.WriteString("\n"); err != nil {
			return nil, fmt.Errorf("failed to separate config entries in %q: %w", path, err)
		}
	}

	if _, err := file.WriteString(formatHostEntry(host)); err != nil {
		return nil, fmt.Errorf("failed to write host %q to %s: %w", host.Alias, path, err)
	}

	return &host, nil
}

func readExistingHosts(path string) ([]Host, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to access config file %q: %w", path, err)
	}

	return ParseHosts(path)
}

func formatHostEntry(host Host) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Host %s\n", host.Alias))
	builder.WriteString(fmt.Sprintf("  HostName %s\n", host.Hostname))
	builder.WriteString(fmt.Sprintf("  User %s\n", host.User))
	if host.Port > 0 {
		builder.WriteString(fmt.Sprintf("  Port %d\n", host.Port))
	}
	if host.IdentityFile != "" {
		builder.WriteString(fmt.Sprintf("  IdentityFile %s\n", host.IdentityFile))
	}

	return builder.String()
}

func resolveAlias(host Host) string {
	if alias := strings.TrimSpace(host.Alias); alias != "" {
		return alias
	}

	if host.Port > 0 {
		return fmt.Sprintf("%s-%d", host.Hostname, host.Port)
	}

	return host.Hostname
}

func validateHostFields(host Host) error {
	fields := []struct {
		name     string
		value    string
		required bool
	}{
		{name: "alias", value: host.Alias, required: true},
		{name: "host name", value: host.Hostname, required: true},
		{name: "user", value: host.User, required: true},
		{name: "identity file", value: host.IdentityFile},
	}

	for _, field := range fields {
		if field.required && strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("%s is required", field.name)
		}

		if containsUnsafeWhitespace(field.value) {
			return fmt.Errorf("%s cannot contain whitespace", field.name)
		}
	}

	return nil
}

func containsUnsafeWhitespace(value string) bool {
	return strings.ContainsFunc(value, unicode.IsSpace)
}
