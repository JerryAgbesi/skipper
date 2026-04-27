package addform

import (
	"strings"
	"testing"
)

func TestValidateRequiredField(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError string
	}{
		{name: "non-empty no whitespace", value: "devone"},
		{name: "empty", value: "", wantError: "is required"},
		{name: "only whitespace", value: "   ", wantError: "is required"},
		{name: "internal whitespace", value: "dev one", wantError: "cannot contain whitespace"},
		{name: "tab character", value: "dev\tone", wantError: "cannot contain whitespace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequiredField("alias", tt.value)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateOptionalField(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError string
	}{
		{name: "empty allowed", value: ""},
		{name: "whitespace only allowed", value: "   "},
		{name: "valid value", value: "devone"},
		{name: "internal whitespace", value: "dev one", wantError: "cannot contain whitespace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOptionalField("alias", tt.value)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError string
	}{
		{name: "empty allowed", value: ""},
		{name: "whitespace only allowed", value: "   "},
		{name: "valid mid-range", value: "22"},
		{name: "valid edge low", value: "1"},
		{name: "valid edge high", value: "65535"},
		{name: "zero out of range", value: "0", wantError: "between 1 and 65535"},
		{name: "too high", value: "70000", wantError: "between 1 and 65535"},
		{name: "negative", value: "-1", wantError: "between 1 and 65535"},
		{name: "non-numeric", value: "abc", wantError: "must be a number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePort(tt.value)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
			}
		})
	}
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{name: "empty", value: "", want: 0},
		{name: "whitespace", value: "  ", want: 0},
		{name: "number", value: "2222", want: 2222},
		{name: "non-numeric", value: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePort(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
