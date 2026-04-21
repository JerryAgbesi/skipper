package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jerryagbesi/skipper/internal/sshconfig"
)

func TestAddCommandWritesHost(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")

	out := new(bytes.Buffer)
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{"add", "devone", "user@10.0.0.8:9000", "-c", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected add command to succeed, got %v", err)
	}

	hosts, err := sshconfig.ParseHosts(path)
	if err != nil {
		t.Fatalf("expected config to parse, got %v", err)
	}

	if len(hosts) != 1 || hosts[0].Alias != "devone" {
		t.Fatalf("expected single devone host, got %+v", hosts)
	}
}

func TestAddCommandRejectsWrongArgCount(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")

	out := new(bytes.Buffer)
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{"add", "devone", "-c", path})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing target argument")
	}

	if !strings.Contains(err.Error(), "accepts 2 arg") {
		t.Fatalf("expected arg-count error, got %v", err)
	}
}
