package list

import (
	"strings"
	"testing"
)

func TestNewCmdList(t *testing.T) {
	cmd := NewCmdList()

	if cmd == nil {
		t.Fatal("NewCmdList() returned nil")
	}

	if cmd.Use != "list [flags] <organization>" {
		t.Errorf("Expected Use to be 'list [flags] <organization>', got %s", cmd.Use)
	}

	// Fix: Update to match actual short description
	if cmd.Short != "Generate a report of repos that repository collaborators have access to." {
		t.Errorf("Expected Short description 'Generate a report of repos that repository collaborators have access to.', got %s", cmd.Short)
	}

	// Check that command requires at least 1 argument
	if cmd.Args == nil {
		t.Error("Expected Args to be set")
	}
}

func TestListCommandFlags(t *testing.T) {
	cmd := NewCmdList()

	// Test that all expected flags exist - using actual flag names from the implementation
	expectedFlags := map[string]string{
		"token":       "t",
		"hostname":    "",
		"username":    "u",
		"output-file": "o",
		"debug":       "d",
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

func TestListCommandRequiredFlags(t *testing.T) {
	cmd := NewCmdList()

	// The output-file flag has a default value so it's not marked as required
	// Let's test that it exists and has a default value
	outputFlag := cmd.Flag("output-file")
	if outputFlag == nil {
		t.Fatal("Expected 'output-file' flag to exist")
	}

	// Check that it has a default value (should contain "RepoCollaboratorsReport")
	if !strings.Contains(outputFlag.DefValue, "RepoCollaboratorsReport") {
		t.Errorf("Expected 'output-file' flag to have default value containing 'RepoCollaboratorsReport', got %s", outputFlag.DefValue)
	}
}

func TestListCommandDefaultValues(t *testing.T) {
	cmd := NewCmdList()

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

	// Check that output-file has a default value
	outputFlag := cmd.Flag("output-file")
	if outputFlag != nil && !strings.Contains(outputFlag.DefValue, "RepoCollaboratorsReport") {
		t.Errorf("Expected default output-file to contain 'RepoCollaboratorsReport', got %s", outputFlag.DefValue)
	}

	// Check that username flag has empty default
	usernameFlag := cmd.Flag("username")
	if usernameFlag != nil && usernameFlag.DefValue != "" {
		t.Errorf("Expected default username to be empty, got %s", usernameFlag.DefValue)
	}

	// Check that token flag has empty default
	tokenFlag := cmd.Flag("token")
	if tokenFlag != nil && tokenFlag.DefValue != "" {
		t.Errorf("Expected default token to be empty, got %s", tokenFlag.DefValue)
	}
}

func TestCmdFlagsStruct(t *testing.T) {
	flags := cmdFlags{
		token:    "test-token",
		hostname: "github.enterprise.com",
		listFile: "test-output.csv",
		username: "testuser",
		debug:    true,
	}

	if flags.token != "test-token" {
		t.Errorf("Expected token to be 'test-token', got %s", flags.token)
	}

	if flags.hostname != "github.enterprise.com" {
		t.Errorf("Expected hostname to be 'github.enterprise.com', got %s", flags.hostname)
	}

	if flags.listFile != "test-output.csv" {
		t.Errorf("Expected listFile to be 'test-output.csv', got %s", flags.listFile)
	}

	if flags.username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", flags.username)
	}

	if !flags.debug {
		t.Error("Expected debug to be true")
	}
}

func TestListCommandArgsValidation(t *testing.T) {
	cmd := NewCmdList()

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

func TestListCommandLongDescription(t *testing.T) {
	cmd := NewCmdList()

	// Check the long description
	if cmd.Long != "Generate a report of repos that repository collaborators have access to." {
		t.Errorf("Expected Long description 'Generate a report of repos that repository collaborators have access to.', got %s", cmd.Long)
	}
}
