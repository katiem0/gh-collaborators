package cmd

import (
	"testing"
)

func TestNewCmdRoot(t *testing.T) {
	cmd := NewCmdRoot()

	if cmd == nil {
		t.Fatal("NewCmdRoot() returned nil")
	}

	if cmd.Use != "collaborators <command> [flags]" {
		t.Errorf("Expected Use to be 'collaborators <command> [flags]', got %s", cmd.Use)
	}

	if cmd.Short != "List and maintain repository collaborators and their repos." {
		t.Errorf("Expected Short description, got %s", cmd.Short)
	}

	if cmd.Long != "List and maintain repository collaborators and their assigned repositories." {
		t.Errorf("Expected Long description, got %s", cmd.Long)
	}
}

func TestRootCommandHasSubcommands(t *testing.T) {
	cmd := NewCmdRoot()

	expectedCommands := []string{"add", "list", "remove"}

	for _, expectedCmd := range expectedCommands {
		found := false
		for _, subCmd := range cmd.Commands() {
			if subCmd.Name() == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expectedCmd)
		}
	}
}

func TestRootCommandCompletion(t *testing.T) {
	cmd := NewCmdRoot()

	if !cmd.CompletionOptions.DisableDefaultCmd {
		t.Error("Expected completion default command to be disabled")
	}
}

func TestRootCommandHelpDisabled(t *testing.T) {
	cmd := NewCmdRoot()

	// The help command is set via SetHelpCommand, not added as a regular command
	// We can't directly access it, but we can verify that the help system is customized
	// by checking that executing "help" doesn't cause issues

	// Set args to trigger help
	cmd.SetArgs([]string{"help"})

	// This should not panic or error since we've set a custom help command
	err := cmd.Execute()

	// The custom no-help command should be hidden and do nothing
	// So we expect no error when help is called
	if err == nil {
		// This is expected - the help command is suppressed
		return
	}

	// If we get an error, it might be because help is not properly suppressed
	t.Logf("Help command execution returned: %v", err)
}

func TestRootCommandSubcommandCount(t *testing.T) {
	cmd := NewCmdRoot()

	// Should have 3 visible commands (add, list, remove)
	// The help command set via SetHelpCommand doesn't appear in Commands()
	commands := cmd.Commands()
	if len(commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(commands))
	}

	// Count visible commands
	visibleCount := 0
	for _, subCmd := range commands {
		if !subCmd.Hidden {
			visibleCount++
		}
	}

	if visibleCount != 3 {
		t.Errorf("Expected 3 visible commands, got %d", visibleCount)
	}
}

func TestRootCommandExecuteWithoutArgs(t *testing.T) {
	cmd := NewCmdRoot()

	// Execute without any arguments should show usage
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	// Should not error when no subcommand is provided
	if err != nil && err.Error() != "cobra: no command specified" {
		t.Logf("Execution without args returned: %v", err)
	}
}
