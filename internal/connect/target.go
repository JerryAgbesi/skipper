package connect

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/jerryagbesi/skipper/internal/sshconfig"
)

func ParseTarget(target string) (*sshconfig.Host, error) {
	user, rawHost, ok := strings.Cut(strings.TrimSpace(target), "@")
	if !ok || user == "" || rawHost == "" || strings.Contains(rawHost, "@") {
		return nil, fmt.Errorf("target must be in the format user@host[:port]")
	}

	hostname, port, err := parseHostPort(rawHost)
	if err != nil {
		return nil, err
	}

	return &sshconfig.Host{
		User:     user,
		Hostname: hostname,
		Port:     port,
	}, nil
}

func parseHostPort(rawHost string) (string, int, error) {
	rawHost = strings.TrimSpace(rawHost)
	if rawHost == "" {
		return "", 0, fmt.Errorf("target must include a host")
	}

	if strings.HasPrefix(rawHost, "[") && strings.HasSuffix(rawHost, "]") {
		return strings.Trim(rawHost, "[]"), 0, nil
	}

	hostname, portText, err := net.SplitHostPort(rawHost)
	if err == nil {
		port, err := strconv.Atoi(portText)
		if err != nil || port <= 0 || port > 65535 {
			return "", 0, fmt.Errorf("port must be between 1 and 65535")
		}

		return strings.Trim(hostname, "[]"), port, nil
	}

	var addrErr *net.AddrError
	if errors.As(err, &addrErr) && strings.Contains(addrErr.Err, "missing port in address") {
		if strings.HasSuffix(rawHost, ":") {
			return "", 0, fmt.Errorf("target must be in the format user@host[:port]")
		}
		return strings.Trim(rawHost, "[]"), 0, nil
	}

	if strings.Contains(rawHost, ":") {
		return "", 0, fmt.Errorf("target must be in the format user@host[:port]")
	}

	return rawHost, 0, nil
}
