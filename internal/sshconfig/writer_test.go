package sshconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddHostRejectsUnsafeWhitespaceInWrittenFields(t *testing.T) {
	tests := []struct {
		name        string
		host        Host
		wantMessage string
	}{
		{
			name: "alias with space",
			host: Host{
				Alias:    "jump box",
				Hostname: "example.com",
				User:     "alice",
			},
			wantMessage: "alias cannot contain whitespace",
		},
		{
			name: "hostname with newline",
			host: Host{
				Alias:    "jump-box",
				Hostname: "example.com\nProxyCommand yes",
				User:     "alice",
			},
			wantMessage: "host name cannot contain whitespace",
		},
		{
			name: "user with tab",
			host: Host{
				Alias:    "jump-box",
				Hostname: "example.com",
				User:     "alice\tadmin",
			},
			wantMessage: "user cannot contain whitespace",
		},
		{
			name: "identity file with space",
			host: Host{
				Alias:        "jump-box",
				Hostname:     "example.com",
				User:         "alice",
				IdentityFile: "/Users/test/.ssh/my key",
			},
			wantMessage: "identity file cannot contain whitespace",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "config")

			_, err := AddHost(configPath, test.host)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), test.wantMessage) {
				t.Fatalf("expected error containing %q, got %v", test.wantMessage, err)
			}

			content, readErr := os.ReadFile(configPath)
			if !os.IsNotExist(readErr) {
				t.Fatalf("expected config file to not be created, readErr=%v content=%q", readErr, string(content))
			}
		})
	}
}

func TestAddHostRejectsUnsafeWhitespaceInResolvedAlias(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config")

	_, err := AddHost(configPath, Host{
		Hostname: "example host",
		User:     "alice",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "alias cannot contain whitespace") {
		t.Fatalf("expected resolved alias validation error, got %v", err)
	}
}

func TestAddHostWritesSafeIdentityFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config")

	_, err := AddHost(configPath, Host{
		Alias:        "jump-box",
		Hostname:     "example.com",
		User:         "alice",
		IdentityFile: "/Users/test/.ssh/id_ed25519",
	})
	if err != nil {
		t.Fatalf("expected host to be added, got %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("expected config to be readable, got %v", err)
	}

	if !strings.Contains(string(content), "  IdentityFile /Users/test/.ssh/id_ed25519\n") {
		t.Fatalf("expected identity file to be written, got content:\n%s", string(content))
	}
}
