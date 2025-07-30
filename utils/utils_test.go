package utils

import (
	"testing"
)

func TestCheckCommand(t *testing.T) {
	// Test with a command that should exist on most systems
	if !CheckCommand("go") {
		t.Skip("Go command not found, skipping test")
	}
	
	// Test with a command that definitely doesn't exist
	if CheckCommand("this-command-definitely-does-not-exist-12345") {
		t.Error("Expected false for non-existent command")
	}
}

func TestFileExists(t *testing.T) {
	// Test with current file (should exist)
	if !FileExists("utils.go") {
		t.Error("Expected utils.go to exist")
	}
	
	// Test with non-existent file
	if FileExists("non-existent-file-12345.txt") {
		t.Error("Expected false for non-existent file")
	}
}