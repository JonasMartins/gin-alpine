package utils

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestGetExecutableDir(t *testing.T) {
	dir, err := GetExecutableDir()
	if err != nil {
		t.Fatalf("GetExecutableDir failed: %v", err)
	}

	if dir == "" {
		t.Error("GetExecutableDir returned empty string")
	}

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Executable directory does not exist: %s", dir)
	}
}

func TestGetProjectRoot(t *testing.T) {
	root, err := GetProjectRoot()
	if err != nil {
		t.Fatalf("GetProjectRoot failed: %v", err)
	}

	if root == "" {
		t.Error("GetProjectRoot returned empty string")
	}

	// Check if the directory exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		t.Errorf("Project root directory does not exist: %s", root)
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err = os.Remove(tmpFile.Name()); err != nil {
			log.Printf("close error %v", err)
		}
	}()
	if err = tmpFile.Close(); err != nil {
		log.Printf("close error %v", err)
	}

	// Test existing file
	if !FileExists(tmpFile.Name()) {
		t.Errorf("FileExists should return true for existing file: %s", tmpFile.Name())
	}

	// Test non-existing file
	nonExistentFile := filepath.Join(os.TempDir(), "non_existent_file_12345.txt")
	if FileExists(nonExistentFile) {
		t.Errorf("FileExists should return false for non-existing file: %s", nonExistentFile)
	}

	// Test directory (should return false)
	tmpDir, err := os.MkdirTemp("", "test_dir_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err = os.RemoveAll(tmpDir); err != nil {
			log.Printf("close error %v", err)
		}
	}()

	if FileExists(tmpDir) {
		t.Errorf("FileExists should return false for directory: %s", tmpDir)
	}
}

func TestReadFile(t *testing.T) {
	// Create a temporary file with content
	tmpFile, err := os.CreateTemp("", "test_read_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err = os.Remove(tmpFile.Name()); err != nil {
			log.Printf("close error %v", err)
		}
	}()

	testContent := "Hello, World!\nThis is a test file."
	if _, err = tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err = tmpFile.Close(); err != nil {
		log.Printf("close error %v", err)
	}

	// Test reading existing file
	content, err := ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, string(content))
	}
}
