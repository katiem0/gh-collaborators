package add

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCmdAdd(t *testing.T) {
	cmd := NewCmdAdd()

	if cmd == nil {
		t.Fatal("NewCmdAdd() returned nil")
	}

	if cmd.Use != "add [flags] <organization>" {
		t.Errorf("Expected Use to be 'add [flags] <organization>', got %s", cmd.Use)
	}

	if cmd.Short != "Add repo access for repository collaborators." {
		t.Errorf("Expected Short description, got %s", cmd.Short)
	}

	if cmd.Long != "Add repositories and permissions for repository collaborators." {
		t.Errorf("Expected Long description, got %s", cmd.Long)
	}

	// Check that command requires at least 1 argument
	if cmd.Args == nil {
		t.Error("Expected Args to be set")
	}
}

func TestAddCommandFlags(t *testing.T) {
	cmd := NewCmdAdd()

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

func TestAddCommandRequiredFlags(t *testing.T) {
	cmd := NewCmdAdd()

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

func TestAddCommandDefaultValues(t *testing.T) {
	cmd := NewCmdAdd()

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

func TestCmdFlagsStruct(t *testing.T) {
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

func TestAddCommandArgsValidation(t *testing.T) {
	cmd := NewCmdAdd()

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

func TestAddCommandLongDescription(t *testing.T) {
	cmd := NewCmdAdd()

	expectedLong := "Add repositories and permissions for repository collaborators."
	if cmd.Long != expectedLong {
		t.Errorf("Expected Long description to be '%s', got '%s'", expectedLong, cmd.Long)
	}
}

func TestAddCommandUsage(t *testing.T) {
	cmd := NewCmdAdd()

	expectedUse := "add [flags] <organization>"
	if cmd.Use != expectedUse {
		t.Errorf("Expected Use to be '%s', got '%s'", expectedUse, cmd.Use)
	}
}

func TestAddCommandFlagDescriptions(t *testing.T) {
	cmd := NewCmdAdd()

	expectedDescriptions := map[string]string{
		"token":     `GitHub Personal Access Token (default "gh auth token")`,
		"hostname":  "GitHub Enterprise Server hostname",
		"from-file": "Path and Name of CSV file to create access from (required)",
		"debug":     "To debug logging",
	}

	for flagName, expectedDesc := range expectedDescriptions {
		flag := cmd.Flag(flagName)
		if flag == nil {
			t.Errorf("Flag '%s' not found", flagName)
			continue
		}

		if flag.Usage != expectedDesc {
			t.Errorf("Expected flag '%s' usage to be '%s', got '%s'", flagName, expectedDesc, flag.Usage)
		}
	}
}

func TestAddCommandFlagTypes(t *testing.T) {
	cmd := NewCmdAdd()

	// Test that flags have correct types
	stringFlags := []string{"token", "hostname", "from-file"}
	for _, flagName := range stringFlags {
		flag := cmd.Flag(flagName)
		if flag == nil {
			t.Errorf("Flag '%s' not found", flagName)
			continue
		}

		if flag.Value.Type() != "string" {
			t.Errorf("Expected flag '%s' to be string type, got %s", flagName, flag.Value.Type())
		}
	}

	// Test debug flag is boolean
	debugFlag := cmd.Flag("debug")
	if debugFlag != nil && debugFlag.Value.Type() != "bool" {
		t.Errorf("Expected debug flag to be bool type, got %s", debugFlag.Value.Type())
	}
}

func TestAddCommandHelpText(t *testing.T) {
	cmd := NewCmdAdd()

	// Test that help can be generated without errors
	help := cmd.UsageString()
	if len(help) == 0 {
		t.Error("Expected non-empty help text")
	}

	// Check that key elements are in the help text
	expectedElements := []string{"add", "organization", "from-file"}
	for _, element := range expectedElements {
		if !strings.Contains(help, element) {
			t.Errorf("Expected help text to contain '%s'", element)
		}
	}
}

func TestAddCommandPersistentFlags(t *testing.T) {
	cmd := NewCmdAdd()

	// Test that certain flags are persistent
	persistentFlags := []string{"token", "hostname", "debug"}
	for _, flagName := range persistentFlags {
		flag := cmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected '%s' to be a persistent flag", flagName)
		}
	}

	// Test that from-file is not persistent (it's a local flag)
	flag := cmd.PersistentFlags().Lookup("from-file")
	if flag != nil {
		t.Error("Expected 'from-file' to not be a persistent flag")
	}

	// Test that from-file is a local flag
	flag = cmd.Flags().Lookup("from-file")
	if flag == nil {
		t.Error("Expected 'from-file' to be a local flag")
	}
}

func TestCmdFlagsEmptyValues(t *testing.T) {
	flags := cmdFlags{}

	if flags.token != "" {
		t.Errorf("Expected empty token, got %s", flags.token)
	}

	if flags.hostname != "" {
		t.Errorf("Expected empty hostname, got %s", flags.hostname)
	}

	if flags.fileName != "" {
		t.Errorf("Expected empty fileName, got %s", flags.fileName)
	}

	if flags.debug {
		t.Error("Expected debug to be false")
	}
}
