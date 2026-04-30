package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jerryagbesi/skipper/internal/sshconfig"
	"github.com/jerryagbesi/skipper/internal/ui/addform"
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

func TestAddCommandRejectsSingleArg(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")

	out := new(bytes.Buffer)
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{"add", "devone", "-c", path})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for partial argument set")
	}

	if !strings.Contains(err.Error(), "no arguments or both") {
		t.Fatalf("expected partial-args error, got %v", err)
	}
}

func TestAddCommandInteractiveWritesHost(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")

	originalRunner := addFormRunner
	t.Cleanup(func() { addFormRunner = originalRunner })

	addFormRunner = func(_ addform.Input) (addform.Result, error) {
		return addform.Result{
			Alias:    "interactive",
			User:     "alice",
			Hostname: "10.0.0.42",
			Port:     2222,
		}, nil
	}

	out := new(bytes.Buffer)
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{"add", "-c", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected interactive add to succeed, got %v", err)
	}

	hosts, err := sshconfig.ParseHosts(path)
	if err != nil {
		t.Fatalf("expected config to parse, got %v", err)
	}

	if len(hosts) != 1 {
		t.Fatalf("expected 1 host, got %d (%+v)", len(hosts), hosts)
	}
	got := hosts[0]
	if got.Alias != "interactive" || got.User != "alice" || got.Hostname != "10.0.0.42" || got.Port != 2222 {
		t.Fatalf("unexpected host written: %+v", got)
	}
}

func TestAddCommandInteractiveDerivesAliasWhenBlank(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")

	originalRunner := addFormRunner
	t.Cleanup(func() { addFormRunner = originalRunner })

	addFormRunner = func(_ addform.Input) (addform.Result, error) {
		return addform.Result{
			User:     "alice",
			Hostname: "10.0.0.42",
			Port:     2222,
		}, nil
	}

	out := new(bytes.Buffer)
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{"add", "-c", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected interactive add to succeed, got %v", err)
	}

	hosts, err := sshconfig.ParseHosts(path)
	if err != nil {
		t.Fatalf("expected config to parse, got %v", err)
	}

	if len(hosts) != 1 {
		t.Fatalf("expected 1 host, got %d (%+v)", len(hosts), hosts)
	}
	if hosts[0].Alias != "10.0.0.42-2222" {
		t.Fatalf("expected alias derived from host:port, got %q", hosts[0].Alias)
	}
}

func TestAddCommandInteractiveCancelledWritesNothing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")

	originalRunner := addFormRunner
	t.Cleanup(func() { addFormRunner = originalRunner })

	addFormRunner = func(_ addform.Input) (addform.Result, error) {
		return addform.Result{Cancelled: true}, nil
	}

	out := new(bytes.Buffer)
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs([]string{"add", "-c", path})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected cancelled add to succeed, got %v", err)
	}

	if _, err := sshconfig.ParseHosts(path); err == nil {
		t.Fatal("expected no config file to be created when form is cancelled")
	}
}
