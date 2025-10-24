package remove

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCmdRemove(t *testing.T) {
	cmd := NewCmdRemove()

	if cmd == nil {
		t.Fatal("NewCmdRemove() returned nil")
	}

	if cmd.Use != "remove [flags] <organization>" {
		t.Errorf("Expected Use to be 'remove [flags] <organization>', got %s", cmd.Use)
	}

	// Update expected Short description to match actual implementation
	if cmd.Short != "Remove repo access for repository collaborators." {
		t.Errorf("Expected Short description, got %s", cmd.Short)
	}

	// Check that command requires at least 1 argument
	if cmd.Args == nil {
		t.Error("Expected Args to be set")
	}
}

func TestRemoveCommandFlags(t *testing.T) {
	cmd := NewCmdRemove()

	// Test that all expected flags exist
	expectedFlags := map[string]string{
		"token":     "t",
		"hostname":  "",
		"from-file": "f",
		"debug":     "d",
	}

	for flag, shorthand := range expectedFlags {
		f := cmd.Flag(flag)
		if f == nil {
			t.Errorf("Expected flag '%s' to exist", flag)
			continue
		}

		if shorthand != "" && f.Shorthand != shorthand {
			t.Errorf("Expected flag '%s' to have shorthand '%s', got '%s'", flag, shorthand, f.Shorthand)
		}
	}
}

func TestRemoveCommandRequiredFlags(t *testing.T) {
	cmd := NewCmdRemove()

	// Test that the from-file flag is required
	fromFileFlag := cmd.Flag("from-file")
	if fromFileFlag == nil {
		t.Fatal("Expected 'from-file' flag to exist")
	}

	// Check if the flag is marked as required
	if fromFileFlag.Annotations == nil {
		t.Error("Expected 'from-file' flag to have annotations")
	} else if _, ok := fromFileFlag.Annotations[cobra.BashCompOneRequiredFlag]; !ok {
		t.Error("Expected 'from-file' flag to be marked as required")
	}
}

func TestRemoveCommandDefaultValues(t *testing.T) {
	cmd := NewCmdRemove()

	// Check default hostname
	hostnameFlag := cmd.Flag("hostname")
	if hostnameFlag != nil && hostnameFlag.DefValue != "github.com" {
		t.Errorf("Expected default hostname to be 'github.com', got %s", hostnameFlag.DefValue)
	}

	// Check default debug value
	debugFlag := cmd.Flag("debug")
	if debugFlag != nil && debugFlag.DefValue != "false" {
		t.Errorf("Expected default debug to be 'false', got %s", debugFlag.DefValue)
	}
}

func TestRemoveCmdFlagsStruct(t *testing.T) {
	flags := cmdFlags{
		token:    "test-token",
		hostname: "github.enterprise.com",
		fileName: "test.csv",
		debug:    true,
	}

	if flags.token != "test-token" {
		t.Errorf("Expected token to be 'test-token', got %s", flags.token)
	}

	if flags.hostname != "github.enterprise.com" {
		t.Errorf("Expected hostname to be 'github.enterprise.com', got %s", flags.hostname)
	}

	if flags.fileName != "test.csv" {
		t.Errorf("Expected fileName to be 'test.csv', got %s", flags.fileName)
	}

	if !flags.debug {
		t.Error("Expected debug to be true")
	}
}

func TestRemoveCommandArgsValidation(t *testing.T) {
	cmd := NewCmdRemove()

	// Test with no args - should fail
	cmd.SetArgs([]string{})
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("Expected error with no arguments, got nil")
	}

	// Test with one arg - should pass
	err = cmd.Args(cmd, []string{"test-org"})
	if err != nil {
		t.Errorf("Expected no error with one argument, got %v", err)
	}

	// Test with multiple args - should pass (MinimumNArgs(1))
	err = cmd.Args(cmd, []string{"test-org", "extra-arg"})
	if err != nil {
		t.Errorf("Expected no error with multiple arguments, got %v", err)
	}
}
