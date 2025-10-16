package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSwapInPlaceCopyWithBackup(t *testing.T) {
	// Setup test directory
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "original.txt")
	newPath := filepath.Join(tmpDir, "new.txt")
	backupPath := srcPath + ".original"

	// Create test files
	if err := os.WriteFile(srcPath, []byte("original content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	if err := os.WriteFile(newPath, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Test with backup (noBackup=false)
	if err := SwapInPlaceCopy(srcPath, newPath, false); err != nil {
		t.Fatalf("SwapInPlaceCopy with backup failed: %v", err)
	}

	// Verify backup exists with original content
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Backup file does not exist: %v", err)
	}
	if string(backupContent) != "original content" {
		t.Errorf("Backup content mismatch: got %q, want %q", string(backupContent), "original content")
	}

	// Verify original path has new content
	srcContent, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Source file does not exist: %v", err)
	}
	if string(srcContent) != "new content" {
		t.Errorf("Source content mismatch: got %q, want %q", string(srcContent), "new content")
	}

	// Verify new file is removed (moved to srcPath)
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Error("New file should have been removed")
	}
}

func TestSwapInPlaceCopyNoBackup(t *testing.T) {
	// Setup test directory
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "original.txt")
	newPath := filepath.Join(tmpDir, "new.txt")
	backupPath := srcPath + ".original"
	tmpBackupPath := srcPath + ".tmp_delete"

	// Create test files
	if err := os.WriteFile(srcPath, []byte("original content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	if err := os.WriteFile(newPath, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Test without backup (noBackup=true)
	if err := SwapInPlaceCopy(srcPath, newPath, true); err != nil {
		t.Fatalf("SwapInPlaceCopy without backup failed: %v", err)
	}

	// Verify backup does NOT exist
	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		t.Error("Backup file should not exist when noBackup=true")
	}

	// Verify temp backup is also cleaned up
	if _, err := os.Stat(tmpBackupPath); !os.IsNotExist(err) {
		t.Error("Temporary backup file should have been cleaned up")
	}

	// Verify original path has new content
	srcContent, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Source file does not exist: %v", err)
	}
	if string(srcContent) != "new content" {
		t.Errorf("Source content mismatch: got %q, want %q", string(srcContent), "new content")
	}

	// Verify new file is removed
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Error("New file should have been removed")
	}
}
