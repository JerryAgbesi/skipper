package connect

import (
	"testing"
)

func TestParseTarget_TrailingColon(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"trailing colon plain host", "user@host:"},
		{"trailing colon IPv6", "user@[::1]:"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.input)
			if err == nil {
				t.Errorf("expected error for input %q, got host %+v", tc.input, got)
			}
		})
	}
}

func TestParseTarget_Valid(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		wantUser     string
		wantHostname string
		wantPort     int
	}{
		{"host only", "alice@example.com", "alice", "example.com", 0},
		{"host with port", "alice@example.com:2222", "alice", "example.com", 2222},
		{"IPv6 bracketed no port", "alice@[::1]", "alice", "::1", 0},
		{"IPv6 bracketed with port", "alice@[::1]:22", "alice", "::1", 22},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.input)
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if got.User != tc.wantUser {
				t.Errorf("user: got %q, want %q", got.User, tc.wantUser)
			}
			if got.Hostname != tc.wantHostname {
				t.Errorf("hostname: got %q, want %q", got.Hostname, tc.wantHostname)
			}
			if got.Port != tc.wantPort {
				t.Errorf("port: got %d, want %d", got.Port, tc.wantPort)
			}
		})
	}
}

func TestParseTarget_Invalid(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"no at-sign", "hostonly"},
		{"empty user", "@host"},
		{"empty host", "user@"},
		{"double at-sign", "user@host@extra"},
		{"port out of range", "user@host:99999"},
		{"port zero", "user@host:0"},
		{"non-numeric port", "user@host:abc"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.input)
			if err == nil {
				t.Errorf("expected error for input %q, got host %+v", tc.input, got)
			}
		})
	}
}