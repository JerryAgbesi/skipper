package sshconfig

import (
	"testing"
)

func TestStandardConfigParse(t *testing.T) {
	hosts, err := ParseHosts("../../testdata/config1")
	if err != nil {
		t.Fatal(err)
	}

	if len(hosts) != 6 {
		t.Fatalf("expected 6 hosts, got %d", len(hosts))
	}
}

func TestEmptyConfig(t *testing.T) {
	hosts, err := ParseHosts("../../testdata/config")
	if err != nil {
		t.Fatal(err)
	}

	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(hosts))
	}
}

func TestFileNotFound(t *testing.T) {
	hosts, err := ParseHosts("404file")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(hosts))
	}
}
