package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/katiem0/gh-collaborators/cmd"
)

func TestMain(t *testing.T) {
	// Test that main function doesn't panic when called
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()
}

func TestRootCommand(t *testing.T) {
	// Test that the root command can be instantiated
	rootCmd := cmd.NewCmdRoot()
	if rootCmd == nil {
		t.Fatal("NewCmdRoot() returned nil")
	}

	// Test that help flag works
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	if err != nil {
		t.Errorf("Execute() with --help returned error: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected help output, got empty string")
	}
}

func TestRootCommandWithInvalidCommand(t *testing.T) {
	rootCmd := cmd.NewCmdRoot()
	rootCmd.SetArgs([]string{"invalid-command"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}
}

func TestMainExecution(t *testing.T) {
	main()
}
